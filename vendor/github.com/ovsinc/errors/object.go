package errors

import "bytes"

type Objecter interface {
	String() string
	Bytes() []byte
	Buffer() *bytes.Buffer
}

type object struct {
	data []byte
}

func NewObjectEmpty() Objecter {
	return &object{}
}

func NewObjectFromBytes(v []byte) Objecter {
	return &object{
		data: v,
	}
}

func NewObjectFromString(s string) Objecter {
	return &object{
		data: []byte(s),
	}
}

func (o *object) String() string {
	if o == nil || len(o.data) == 0 {
		return ""
	}
	return string(o.data)
}

func (o *object) Bytes() []byte {
	if o == nil || len(o.data) == 0 {
		return []byte{}
	}
	return o.data
}

func (o *object) Buffer() *bytes.Buffer {
	if o == nil {
		return &bytes.Buffer{}
	}
	return bytes.NewBuffer(o.data)
}

//

type Objects []Objecter

var _ interface {
	Append(oo ...Objecter) Objects
	AppendString(ss ...string) Objects
	AppendBytes(vv ...[]byte) Objects
} = (*Objects)(nil)

func NewObjects(oo ...Objecter) Objects {
	return append(make(Objects, 0, len(oo)), oo...)
}

func NewObjectsFromBytes(vv ...[]byte) Objects {
	objs := make(Objects, 0, len(vv))
	for _, v := range vv {
		objs = append(objs, NewObjectFromBytes(v))
	}
	return objs
}

func NewObjectsFromString(ss ...string) Objects {
	objs := make(Objects, 0, len(ss))
	for _, s := range ss {
		objs = append(objs, NewObjectFromString(s))
	}
	return objs
}

func (os Objects) copy(cap int) Objects {
	objs := make(Objects, 0, cap)
	if len(os) > 0 {
		objs = append(objs, os...)
	}
	return objs
}

func (os Objects) Append(oo ...Objecter) Objects {
	return append(os.copy(len(os)+len(oo)), oo...)
}

func (os Objects) AppendString(ss ...string) Objects {
	objs := os.copy(len(os) + len(ss))
	for _, v := range ss {
		objs = append(objs, NewObjectFromString(v))
	}
	return objs
}

func (os Objects) AppendBytes(vv ...[]byte) Objects {
	objs := os.copy(len(os) + len(vv))
	for _, v := range vv {
		objs = append(objs, NewObjectFromBytes(v))
	}
	return objs
}
