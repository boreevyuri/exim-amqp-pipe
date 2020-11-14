package reader

import (
	"bufio"
	"errors"
	"os"
)

func ReadStdin() (data []byte, err error) {
	var errInput = errors.New("no input specified")

	inputData, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (inputData.Mode() & os.ModeNamedPipe) == 0 {
		return nil, errInput
	}

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}
	data = scanner.Bytes()

	return data, nil
}
