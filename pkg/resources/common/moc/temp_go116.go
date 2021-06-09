// +build go1.16

package moc

import "os"

func CreateTemp(dir string, data []byte) (Temp, error) {
	f, err := os.CreateTemp(dir, pattern)
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
