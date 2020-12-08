package reader

import (
	"errors"
	"fmt"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

const (
	contentTypeHeader        = "Content-Type"
	contentIDHeader          = "Content-Id"
	contentTransferEncHeader = "Content-Transfer-Encoding"
	contentDispositionHeader = "Content-Disposition"
	mixed                    = "multipart/mixed"
	alternative              = "multipart/alternative"
	related                  = "multipart/related"
	html                     = "text/html"
	plain                    = "text/plain"
	binary                   = "application/octet-stream"
	boundaryParam            = "boundary"
)

type Email struct {
	Sender      string
	Rcpt        string
	Attachments []File
}

type File struct {
	Filename           string
	ContentType        string
	ContentDisposition string
	ContentEncoding    string
	Data               []byte
}

func ScanEmail(conf config.ParseConfig, msg *mail.Message) (email Email, err error) {
	var files []File
	switch conf.AttachmentsOnly {
	case true:
		files, err = GetFilesFrom(msg)
	default:
		files, err = ScanFullLetter(msg)
	}
	if err != nil {
		failOnError(err, "Unable to parse email:")
	}

	email = Email{
		Sender:      decodeMimeSentence(msg.Header.Get("From")),
		Rcpt:        GetRecipients(&msg.Header),
		Attachments: files,
	}

	return email, nil
}

func ScanFullLetter(msg *mail.Message) (files []File, err error) {
	data := new(File)

	contentType := msg.Header.Get(contentTypeHeader)
	if len(contentType) == 0 {
		contentType = plain
	}

	fileName := msg.Header.Get("Message-ID")

	if len(fileName) == 0 {
		fileName = decodeMimeSentence(msg.Header.Get("From"))
	}

	data.Data, err = ioutil.ReadAll(msg.Body)

	data.ContentType = contentType
	data.Filename = fileName
	data.ContentEncoding = plain

	files = append(files, *data)

	return files, err
}

func GetFilesFrom(msg *mail.Message) (files []File, err error) {
	contentType, params, err := parseContentType(msg.Header.Get(contentTypeHeader))
	failOnError(err, "Unable to parse email Content-Type")

	switch contentType {
	case mixed:
		files, err = parseMixed(msg.Body, params[boundaryParam])
	case alternative, related:
		files, err = parseMultipart(msg.Body, params[boundaryParam])
	default:
		return
	}

	return files, err
}

func parseContentType(header string) (contentType string, params map[string]string, err error) {
	if header == "" {
		contentType = plain
		return
	}

	return mime.ParseMediaType(header)
}

func parseMixed(msg io.Reader, boundary string) (files []File, err error) {
	r := multipart.NewReader(msg, boundary)

	for {
		part, err := r.NextPart()
		if err != nil {
			//Если нет вложенного part - прерываем обработку
			if errors.Is(err, io.EOF) {
				break
			}

			return files, err
		}

		contentType, params, err := mime.ParseMediaType(part.Header.Get(contentTypeHeader))
		if err != nil {
			contentType = binary
		}

		switch contentType {
		case plain, html:
			break
		case alternative, related:
			files, err = parseMultipart(part, params[boundaryParam])
			if err != nil {
				return files, err
			}
		default:
			if !isAttachment(part) {
				return files, fmt.Errorf(
					"unknown multipart/mixed nested mime type: %s", contentType)
			}

			at, err := createAttachment(part)
			if err != nil {
				return files, err
			}

			files = append(files, at)
		}
	}

	return files, err
}

func parseMultipart(msg io.Reader, boundary string) (files []File, err error) {
	r := multipart.NewReader(msg, boundary)

	for {
		part, err := r.NextPart()
		if err != nil {
			//Если нет вложенного part - прерываем обработку
			if errors.Is(err, io.EOF) {
				break
			}

			return files, err
		}

		contentType, params, err := mime.ParseMediaType(part.Header.Get(contentTypeHeader))
		if err != nil {
			return files, err
		}

		switch contentType {
		case plain, html:
			break
		case related, alternative:
			ef, err := parseMultipart(part, params[boundaryParam])
			if err != nil {
				return files, err
			}

			files = append(files, ef...)

		default:
			if isEmbeddedFile(part) {
				ef, err := createEmbedded(part)
				if err != nil {
					return files, err
				}

				files = append(files, ef)
			} else {
				return files, fmt.Errorf(
					"can't process multipart/(alternative|related) inner type: %s", contentType)
			}
		}
	}

	return files, err
}

func createEmbedded(part *multipart.Part) (file File, err error) {
	cid := decodeMimeSentence(part.Header.Get(contentIDHeader))
	file.Filename = strings.Trim(cid, "<>")
	file.ContentType = part.Header.Get(contentTypeHeader)
	file.ContentEncoding = part.Header.Get(contentTransferEncHeader)
	file.ContentDisposition = part.Header.Get(contentDispositionHeader)

	file.Data, err = ioutil.ReadAll(part)
	if err != nil {
		return
	}

	return file, err
}

func createAttachment(part *multipart.Part) (file File, err error) {
	file.Filename = decodeMimeSentence(part.FileName())
	file.ContentType = part.Header.Get(contentTypeHeader)
	if file.ContentType == "" {
		file.ContentType = binary + "; name=\"unknown\""
	}
	file.ContentEncoding = part.Header.Get(contentTransferEncHeader)
	file.ContentDisposition = part.Header.Get(contentDispositionHeader)
	file.Data, err = ioutil.ReadAll(part)
	if err != nil {
		return
	}

	return file, err
}

func isAttachment(part *multipart.Part) bool {
	return part.FileName() != ""
}

func isEmbeddedFile(part *multipart.Part) bool {
	return part.Header.Get(contentTransferEncHeader) != ""
}

var CharsetReader = func(label string, input io.Reader) (io.Reader, error) {
	label = strings.Replace(label, "windows-", "cp", -1)
	encoding, _ := charset.Lookup(label)
	return encoding.NewDecoder().Reader(input), nil
}

func decodeMimeSentence(s string) string {
	ss := strings.Fields(s)
	result := make([]string, 0, len(ss))

	for _, word := range ss {
		dec := mime.WordDecoder{CharsetReader: CharsetReader}

		w, err := dec.Decode(word)
		if err != nil {
			if len(result) == 0 {
				w = word
			} else {
				w = "_" + word
			}
		}

		result = append(result, w)
	}

	return strings.Join(result, "")
}

//func decodeHeader(s string) string {
//	dec := mime.WordDecoder{CharsetReader: CharsetReader}
//	header, err := dec.DecodeHeader(s)
//	if err != nil {
//		return s
//	}
//	return header
//}

func GetRecipients(head *mail.Header) (rcpts string) {
	var rcptHeaders = [...]string{
		"To",
		"Envelope-To",
		"X-Envelope-To",
		"Cc",
		"Bcc",
	}

	unique := make(map[string]bool)
	addrs := make([]string, 0, 1)

	for _, h := range rcptHeaders {
		t, _ := head.AddressList(h)
		for _, addr := range t {
			if !unique[addr.Address] {
				addrs = append(addrs, addr.Address)
				unique[addr.Address] = true
			}
		}
	}
	for _, key := range addrs {
		if rcpts != "" {
			rcpts = rcpts + ", "
		}
		rcpts = rcpts + key
	}

	return rcpts
}
