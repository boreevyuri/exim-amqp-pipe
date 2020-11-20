package reader

import (
	"errors"
	"fmt"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"log"
	"net/mail"
	"os"
)

//func ReadStdin() string {
//	var errInput = errors.New("no input specified")
//
//	inputData, err := os.Stdin.Stat()
//	failOnError(err, "no input specified(1):")
//
//	if (inputData.Mode() & os.ModeNamedPipe) == 0 {
//		failOnError(errInput, "no input specified(2):")
//	}
//
//	data, err := ioutil.ReadAll(os.Stdin)
//	failOnError(err, "Unable to readAll os.Stdin:")
//
//	return string(data)
//}

func ReadMail() (msg *mail.Message) {
	var errInput = errors.New("no input specified")

	inputData, err := os.Stdin.Stat()
	failOnError(err, "no input specified(1):")

	if (inputData.Mode() & os.ModeNamedPipe) == 0 {
		failOnError(errInput, "no input specified(2):")
	}

	msg, err = mail.ReadMessage(os.Stdin)
	failOnError(err, "Unable to read mail from os.Stdin:")

	return msg
}

//func Parse(outgoing chan File, parseConf config.ParseConfig) {
func Parse(outgoing chan File, parseConf config.ParseConfig) {

	//TODO: доделать конфигурацию с публикацией всего письма.
	//TODO: м.б. обернуть в структ

	//if !parseConf.AttachmentsOnly {
	//	data := ReadStdin()
	//	outgoing <- data
	//	close(outgoing)
	//}

	message := ReadMail()
	email, err := ParseMail(message)
	if err != nil {
		failOnError(err, "Unable to parse email:")
	}

	log.Printf("Attachments found: %d", len(email.Files))

	for _, file := range email.Files {
		fmt.Printf("Got file with name %s\n", file.Filename)
		outgoing <- file
	}
	fmt.Printf("All files gone. Closing Parse\n")
	close(outgoing)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
