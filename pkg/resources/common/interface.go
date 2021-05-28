package common

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
