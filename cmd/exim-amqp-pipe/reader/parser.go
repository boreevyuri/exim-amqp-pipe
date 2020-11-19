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
	case alternative:
		e.Files, err = parseAlternative(msg.Body, params["boundary"])
	case related:
		e.Files, err = parseRelated(msg.Body, params["boundary"])
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
		case alternative:
			files, err = parseAlternative(part, params["boundary"])
			if err != nil {
				return files, err
			}
		case related:
			files, err = parseRelated(part, params["boundary"])
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

func parseAlternative(msg io.Reader, boundary string) (files []File, err error) {
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
		case related:
			ef, err := parseRelated(part, params["boundary"])
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
					"can't process multipart/alternative inner type: %s", contentType)
			}
		}
	}

	return files, err
}

func parseRelated(msg io.Reader, boundary string) (files []File, err error) {
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
		case alternative:
			ef, err := parseAlternative(part, params["boundary"])
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
				return files, fmt.Errorf("can't process multipart/alternative inner type: %s", contentType)
			}
		}
	}

	return files, err
}

func createEmbedded(part *multipart.Part) (ef File, err error) {
	ef.Filename = part.Header.Get("Content-Id")
	ef.ContentType = part.Header.Get(contentTypeHeader)
	ef.Data, err = ioutil.ReadAll(part)
	if err != nil {
		return
	}
	return ef, err
}

func createAttachment(part *multipart.Part) (at File, err error) {
	at.Filename = part.FileName()
	at.ContentType = strings.Split(part.Header.Get(contentTypeHeader), ";")[0]
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

type File struct {
	Filename    string
	ContentType string
	Data        []byte
}

type ParsedMail struct {
	Files []File
}
