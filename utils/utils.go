package utils

import (
	"errors"
	"io"
	"os"
)

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)

	if err != nil {
		if errors.Is(err, io.EOF) {
			return true, nil
		} 

		return false, err
	}

	return false, nil
}
