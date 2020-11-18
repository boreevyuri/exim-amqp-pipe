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

	e.ContentType = msg.Header.Get(contentTypeHeader)
	contentType, params, err := parseContentType(e.ContentType)
	failOnError(err, "Unable to parse email Content-Type")

	switch contentType {
	case mixed:
		e.Attachments, e.EmbeddedFiles, err = parseMixed(msg.Body, params["boundary"])
	case alternative:
		e.EmbeddedFiles, err = parseAlternative(msg.Body, params["boundary"])
	case related:
		e.EmbeddedFiles, err = parseRelated(msg.Body, params["boundary"])
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

func parseMixed(msg io.Reader, boundary string) (
	attachments []Attachment, embeddedFiles []EmbeddedFile, err error) {

	r := multipart.NewReader(msg, boundary)
	for {
		part, err := r.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return attachments, embeddedFiles, err
		}

		contentType, params, err := mime.ParseMediaType(part.Header.Get(contentTypeHeader))
		if err != nil {
			return attachments, embeddedFiles, err
		}

		switch contentType {
		case plain, html:
			break
		case alternative:
			embeddedFiles, err = parseAlternative(part, params["boundary"])
			if err != nil {
				return attachments, embeddedFiles, err
			}
		case related:
			embeddedFiles, err = parseRelated(part, params["boundary"])
			if err != nil {
				return attachments, embeddedFiles, err
			}
		default:
			if !isAttachment(part) {
				return attachments, embeddedFiles, fmt.Errorf(
					"unknown multipart/mixed nested mime type: %s", contentType)
			}
			at, err := createAttachment(part)
			if err != nil {
				return attachments, embeddedFiles, err
			}
			attachments = append(attachments, at)

		}

	}
	return attachments, embeddedFiles, err
}

func parseAlternative(msg io.Reader, boundary string) (embeddedFiles []EmbeddedFile, err error) {
	r := multipart.NewReader(msg, boundary)
	for {
		part, err := r.NextPart()

		if err != nil {
			if err == io.EOF {
				break
			}
			return embeddedFiles, err
		}

		contentType, params, err := mime.ParseMediaType(part.Header.Get(contentTypeHeader))
		if err != nil {
			return embeddedFiles, err
		}

		switch contentType {
		case plain, html:
			break
		case related:
			ef, err := parseRelated(part, params["boundary"])
			if err != nil {
				return embeddedFiles, err
			}
			embeddedFiles = append(embeddedFiles, ef...)
		default:
			if isEmbeddedFile(part) {
				ef, err := createEmbedded(part)
				if err != nil {
					return embeddedFiles, err
				}
				embeddedFiles = append(embeddedFiles, ef)
			} else {
				return embeddedFiles, fmt.Errorf(
					"can't process multipart/alternative inner type: %s", contentType)
			}
		}
	}

	return embeddedFiles, err
}

func parseRelated(msg io.Reader, boundary string) (embeddedFiles []EmbeddedFile, err error) {
	r := multipart.NewReader(msg, boundary)
	for {
		part, err := r.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return embeddedFiles, err
		}

		contentType, params, err := mime.ParseMediaType(part.Header.Get(contentTypeHeader))
		if err != nil {
			return embeddedFiles, err
		}

		switch contentType {
		case plain, html:
			break
		case alternative:
			ef, err := parseAlternative(part, params["boundary"])
			if err != nil {
				return embeddedFiles, err
			}
			embeddedFiles = append(embeddedFiles, ef...)
		default:
			if isEmbeddedFile(part) {
				ef, err := createEmbedded(part)
				if err != nil {
					return embeddedFiles, err
				}
				embeddedFiles = append(embeddedFiles, ef)
			} else {
				return embeddedFiles, fmt.Errorf("can't process multipart/alternative inner type: %s", contentType)
			}
		}
	}

	return embeddedFiles, err
}

func createEmbedded(part *multipart.Part) (ef EmbeddedFile, err error) {
	ef.CID = part.Header.Get("Content-Id")
	ef.ContentType = part.Header.Get(contentTypeHeader)
	//ef.Data, err = ioutil.ReadAll(part)
	data, err := ioutil.ReadAll(part)
	if err != nil {
		return
	}
	ef.Data = string(data)
	return ef, err
}

func createAttachment(part *multipart.Part) (at Attachment, err error) {
	at.Filename = part.FileName()
	at.ContentType = strings.Split(part.Header.Get(contentTypeHeader), ";")[0]
	//at.Data, err = ioutil.ReadAll(part)
	data, err := ioutil.ReadAll(part)
	if err != nil {
		return
	}
	at.Data = string(data)
	return at, err
}

func isAttachment(part *multipart.Part) bool {
	return part.FileName() != ""
}

func isEmbeddedFile(part *multipart.Part) bool {
	return part.Header.Get("Content-Transfer-Encoding") != ""
}

type Attachment struct {
	Filename    string
	ContentType string
	Data        string
}

type EmbeddedFile struct {
	CID         string
	ContentType string
	Data        string
}

type ParsedMail struct {
	//Header mail.Header
	ContentType string

	Attachments   []Attachment
	EmbeddedFiles []EmbeddedFile
}
