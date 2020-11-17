package reader

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

func ReadStdin() (data []byte) {
	var errInput = errors.New("no input specified")

	inputData, err := os.Stdin.Stat()
	failOnError(err, "no input specified(1):")

	if (inputData.Mode() & os.ModeNamedPipe) == 0 {
		failOnError(errInput, "no input specified(2):")
	}

	data, err = ioutil.ReadAll(os.Stdin)
	failOnError(err, "Unable to read os.Stdin:")

	return data
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
