package common

import "io"

type ResourceStopper interface {
	Stop()
}

type ResourceViewer interface {
	Used() float64
}

type Resourcer interface {
	ResourceStopper
	ResourceViewer
}

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
