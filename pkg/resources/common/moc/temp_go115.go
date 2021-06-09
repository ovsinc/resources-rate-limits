// +build !go1.16

package moc

import (
	"io/ioutil"
)

func CreateTemp(dir string, data []byte) (Temp, error) {
	f, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(data)
	if err != nil {
		return nil, err
	}

	return &commonTemp{
		f:    f,
		file: f.Name(),
	}, nil
}
