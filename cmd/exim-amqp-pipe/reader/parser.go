package reader

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

const (
	contentTypeHeader = "Content-Type"
	mixed             = "multipart/mixed"
	alternative       = "multipart/alternative"
	related           = "multipart/related"
	html              = "text/html"
	plain             = "text/plain"
)

func ParseMail(msg *mail.Message) (e ParsedMail, err error) {

	contentType, params, err := parseContentType(msg.Header.Get(contentTypeHeader))
	failOnError(err, "Unable to parse email Content-Type")

	switch contentType {
	case mixed:
		e.Files, err = parseMixed(msg.Body, params["boundary"])
	case alternative, related:
		e.Files, err = parseMultipart(msg.Body, params["boundary"])
	default:
		return
	}

	return
}

func parseContentType(contentTypeHeader string) (
	contentType string, params map[string]string, err error) {

	if contentTypeHeader == "" {
		contentType = plain
		return
	}

	return mime.ParseMediaType(contentTypeHeader)
}

func parseMixed(msg io.Reader, boundary string) (files []File, err error) {

	r := multipart.NewReader(msg, boundary)
	for {
		part, err := r.NextPart()
		if err != nil {
			if err == io.EOF {
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
		case alternative, related:
			files, err = parseMultipart(part, params["boundary"])
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
			if err == io.EOF {
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
			ef, err := parseMultipart(part, params["boundary"])
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

func createEmbedded(part *multipart.Part) (ef File, err error) {
	cid := decodeMimeSentence(part.Header.Get("Content-Id"))
	ef.Filename = strings.Trim(cid, "<>")
	ef.ContentType = part.Header.Get(contentTypeHeader)
	ef.ContentEncoding = part.Header.Get("Content-Transfer-Encoding")
	//ef.Data, err = decodeContent(part, part.Header.Get("Content-Transfer-Encoding"))
	ef.Data, err = ioutil.ReadAll(part)
	if err != nil {
		return
	}

	return ef, err
}

func createAttachment(part *multipart.Part) (at File, err error) {
	at.Filename = decodeMimeSentence(part.FileName())
	//at.ContentType = strings.Split(part.Header.Get(contentTypeHeader), ";")[0]
	at.ContentType = part.Header.Get(contentTypeHeader)
	at.ContentEncoding = part.Header.Get("Content-Transfer-Encoding")
	//at.Data, err = decodeContent(part, part.Header.Get("Content-Transfer-Encoding"))
	at.Data, err = ioutil.ReadAll(part)
	if err != nil {
		return
	}
	return at, err
}

func isAttachment(part *multipart.Part) bool {
	return part.FileName() != ""
}

func isEmbeddedFile(part *multipart.Part) bool {
	return part.Header.Get("Content-Transfer-Encoding") != ""
}

func decodeMimeSentence(s string) string {
	ss := strings.Fields(s)
	result := make([]string, 0, len(ss))
	for _, word := range ss {
		dec := new(mime.WordDecoder)
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

//func decodeContent(content io.Reader, encoding string) ([]byte, error) {
//	switch encoding {
//	case "base64":
//		fmt.Printf("Found encoding base64")
//		dec := base64.NewDecoder(base64.StdEncoding, content)
//		data, err := ioutil.ReadAll(dec)
//		if err != nil {
//			return nil, err
//		}
//
//		return data, nil
//	default:
//		data, err := ioutil.ReadAll(content)
//		if err != nil {
//			return nil, err
//		}
//		return data, nil
//	}
//
//}

type File struct {
	Filename        string
	ContentType     string
	ContentEncoding string
	Data            []byte
}

type ParsedMail struct {
	Files []File
}
