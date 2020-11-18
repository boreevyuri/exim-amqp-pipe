package reader

import (
	"errors"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
)

func ReadStdin() (msg []byte) {
	var errInput = errors.New("no input specified")

	inputData, err := os.Stdin.Stat()
	failOnError(err, "no input specified(1):")

	if (inputData.Mode() & os.ModeNamedPipe) == 0 {
		failOnError(errInput, "no input specified(2):")
	}

	msg, err = ioutil.ReadAll(os.Stdin)
	failOnError(err, "Unable to readAll os.Stdin:")

	return msg
}

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

//type commandAction int
//
//const (
//	readAll commandAction = iota
//	parseAttach
//	parseEmbed
//)
//
//type commandData struct {
//	action commandAction
//	result chan<- []byte
//}
//
//type processMessage chan commandData
//
//func (pm processMessage) ParseAttach() []byte {
//	reply := make(chan []byte)
//	pm <- commandData{
//		action: parseAttach,
//		result: reply,
//	}
//	return <-reply
//}
//
//func (pm processMessage) ParseEmbed() []byte {
//	reply := make(chan []byte)
//	pm <- commandData{
//		action: parseEmbed,
//		result: reply,
//	}
//	return <-reply
//}
//
//func (pm processMessage) run(parseConf config.ParseConfig) {
//	data := ReadStdin()
//	if !parseConf.WithEmbeddedFiles {
//
//	}
//	for command := range pm {
//		switch command.action {
//		case readAll:
//			log.Printf("Got readAll signal")
//		case parseAttach:
//			log.Printf("Got parseAttach signal")
//		case parseEmbed:
//			log.Printf("Got parseEmbed signal")
//		}
//	}
//}
//
//
//type ProcessMessage interface {
//	ParseAttach() bool
//	ParseEmbed() bool
//}
//
//func New(parseConfig config.ParseConfig) ProcessMessage {
//	readerPipe := make(processMessage)
//	go readerPipe.run(parseConfig)
//	return readerPipe
//}

func Parse(parseConf config.ParseConfig) (reply chan string) {
	reply = make(chan string)
	defer close(reply)

	//if !parseConf.AttachmentsOnly {
	//	data := ReadStdin()
	//	reply <- data
	//	return reply
	//}

	message := ReadMail()
	email, err := ParseMail(message)
	if err != nil {
		failOnError(err, "Unable to parse email:")
	}

	log.Printf("Attachments found: %d", len(email.Attachments))

	for _, attachment := range email.Attachments {

		reply <- attachment.Data
		//fmt.Println(reflect.TypeOf(attachment.Data).String())
	}

	for _, embeddedFile := range email.EmbeddedFiles {
		reply <- embeddedFile.Data
	}

	return reply
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
