package moc

import (
	"io"

	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

var _ rescommon.ResourceConfiger = (*ResourceConfigMoc)(nil)

type ResourceConfigMoc struct {
	Rtype rescommon.ResourceType
	FF    map[string]io.ReadSeekCloser
}

func (rc *ResourceConfigMoc) Init() error {
	return nil
}

func (rc *ResourceConfigMoc) Type() rescommon.ResourceType {
	return rc.Rtype
}

func (rc *ResourceConfigMoc) File(name string) io.ReadSeekCloser {
	return rc.FF[name]
}

func (rc *ResourceConfigMoc) Stop() {
	for _, v := range rc.FF {
		if v != nil {
			_ = v.Close()
		}
	}
}
