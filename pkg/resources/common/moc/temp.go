package moc

import "os"

const (
	pattern = "fs_*_moc"
)

type Temp interface {
	Remove() error
	File() *os.File
}

type commonTemp struct {
	file string
	f    *os.File
}

func (ct *commonTemp) Remove() error {
	defer ct.f.Close()
	return os.Remove(ct.file)
}

func (ct *commonTemp) File() *os.File {
	return ct.f
}
