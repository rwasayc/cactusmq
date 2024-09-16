package server

type options struct {
}

type option interface {
	apply(*options)
}
