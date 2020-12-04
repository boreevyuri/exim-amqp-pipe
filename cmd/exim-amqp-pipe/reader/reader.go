package reader

import (
	"errors"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
)

func ReadStdin() (msg *mail.Message) {
	var errInput = errors.New("no input specified")

	inputData, err := os.Stdin.Stat()
	failOnError(err, "no input specified(1):")

	if (inputData.Mode() & os.ModeNamedPipe) == 0 {
		failOnError(errInput, "no pipe")
	}

	msg, err = mail.ReadMessage(os.Stdin)
	failOnError(err, "Unable to read mail from os.Stdin:")

	return msg
}

func readDir(path string) []*mail.Message {
	msgs := make([]*mail.Message, 0, 1)
	file, err := os.Open(path)
	failOnError(err, "unable to open")
	defer func() {
		err := file.Close()
		failOnError(err, "unable to close")
	}()

	fi, err := file.Stat()
	failOnError(err, "unable to stat")

	var walkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		msg := readFile(path)
		msgs = append(msgs, msg)
		return nil
	}

	switch {
	case fi.IsDir():
		err := filepath.Walk(path, walkFunc)
		failOnError(err, "unable to filepath.Walk")
	default:
		msg := readFile(path)
		msgs = append(msgs, msg)
	}

	return msgs
}

func readFile(filename string) (msg *mail.Message) {
	file, err := os.Open(filename)
	failOnError(err, "unable to open file")
	defer func() {
		err := file.Close()
		failOnError(err, "unable to close file")
	}()

	data, err := ioutil.ReadAll(file)
	msg, err = mail.ReadMessage(strings.NewReader(string(data)))
	failOnError(err, "unable to parse mail message")

	return msg
}

func ReadInput(out chan Email, emlFiles []string, conf config.ParseConfig) {
	var msg *mail.Message
	messages := make([]*mail.Message, 0, 1)

	if len(emlFiles) == 0 {
		msg = ReadStdin()
		messages = append(messages, msg)
	}

	for _, filename := range emlFiles {
		msg := readDir(filename)
		messages = append(messages, msg...)
	}

	for _, msg := range messages {
		email, err := ScanEmail(conf, msg)
		if err != nil {
			failOnError(err, "oops")
		}
		out <- email
	}
	close(out)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
