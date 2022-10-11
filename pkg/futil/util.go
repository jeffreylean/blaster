package futil

import (
	"errors"
	"io/ioutil"
	"net/url"
)

func ReadFile(path string) ([]byte, error) {
	// Get file based on scheme
	u, err := url.Parse(path)
	if err != nil {
		return []byte{}, err
	}

	switch u.Scheme {
	case "file":
		b, err := ioutil.ReadFile(u.Path)
		if err != nil {
			return []byte{}, err
		}
		return b, nil

	}

	return []byte{}, errors.New("file: Incorrect path format.")
}
