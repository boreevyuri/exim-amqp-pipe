package reader

import (
	"errors"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
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

func ReadInput(files chan File, emlFiles []string, conf config.ParseConfig) {
	var msg *mail.Message
	messages := make([]*mail.Message, 0, 1)

	if len(emlFiles) == 0 {
		msg = ReadStdin()
		messages = append(messages, msg)
	}

	for _, filename := range emlFiles {
		msg = readFile(filename)
		messages = append(messages, msg)
	}

	for _, msg := range messages {
		fileSlice := ScanEmail(conf, msg)
		for _, file := range fileSlice {
			files <- file
		}
	}
	close(files)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
