package reader

import (
	"errors"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"log"
	"net/mail"
	"os"
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

func Parse(files chan File, conf config.ParseConfig) {
	msg := ReadStdin()
	fileSlice := ScanEmail(conf, msg)

	for _, file := range fileSlice {
		files <- file
	}
	close(files)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
