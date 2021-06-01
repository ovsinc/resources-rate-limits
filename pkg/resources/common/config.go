package common

import (
	"os"

	"go.uber.org/atomic"

	"github.com/ovsinc/errors"
)

type ResourceConfiger interface {
	Type() ResourceType
	File(string) ReadSeekCloser
	Init() error
	Stop()
}

type resourceConfig struct {
	rtype  ResourceType
	fnames []string
	ff     map[string]ReadSeekCloser
	isinit *atomic.Bool
}

func NewResourceConfig(rtype ResourceType, fnames ...string) ResourceConfiger {
	return &resourceConfig{
		fnames: append(make([]string, 0, len(fnames)), fnames...),
		rtype:  rtype,
		ff:     make(map[string]ReadSeekCloser),
		isinit: atomic.NewBool(false),
	}
}

func (rc *resourceConfig) Type() ResourceType {
	return rc.rtype
}

func (rc *resourceConfig) Init() error {
	if rc.isinit.Load() {
		return nil
	}

	var err error

	for _, name := range rc.fnames {
		var e error
		rc.ff[name], e = os.Open(name)
		if e != nil {
			err = errors.Wrap(err, e)
		}
	}

	rc.isinit.Store(true)

	return err
}

func (rc *resourceConfig) File(name string) ReadSeekCloser {
	return rc.ff[name]
}

func (rc *resourceConfig) Stop() {
	if !rc.isinit.Load() {
		return
	}

	for _, v := range rc.ff {
		if v != nil {
			_ = v.Close()
		}
	}

	rc.ff = make(map[string]ReadSeekCloser)

	rc.isinit.Store(false)
}
