package reader

import (
	"errors"
	"exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"log"
	"net/mail"
	"os"
	"path/filepath"
)

func ReadInput(out chan<- Email, emlFiles []string, conf config.ParseConfig) {
	mailChan := make(chan *mail.Message)

	if len(emlFiles) > 0 {
		for _, filename := range emlFiles {
			go readDir(mailChan, filename)
		}
	} else {
		go readStdin(mailChan)
	}

	for msg := range mailChan {
		email, err := ScanEmail(conf, msg)
		failOnError(err, "ooops")
		out <- email
	}

	close(out)
}

func readStdin(job chan<- *mail.Message) {
	var errInput = errors.New("no input specified")

	inputData, err := os.Stdin.Stat()
	failOnError(err, "no input specified(1)")

	if (inputData.Mode() & os.ModeNamedPipe) == 0 {
		failOnError(errInput, "no pipe")
	}

	msg, err := mail.ReadMessage(os.Stdin)
	failOnError(err, "Unable to read mail from os.Stdin")

	job <- msg
}

func readDir(job chan<- *mail.Message, path string) {
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
		if filepath.Ext(path) == ".eml" {
			file, err := os.Open(path)
			failOnError(err, "unable to open file")
			// defer file.Close()
			m, err := mail.ReadMessage(file)
			failOnError(err, "unable to parse mail message")
			job <- m
		}
		return nil
	}

	switch {
	case fi.IsDir():
		err := filepath.Walk(path, walkFunc)
		failOnError(err, "unable to filepath.Walk")
	default:
		file, err := os.Open(path)
		failOnError(err, "unable to open file")
		// defer file.Close()
		m, err := mail.ReadMessage(file)
		failOnError(err, "unable to parse mail message")
		job <- m
	}

	close(job)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
