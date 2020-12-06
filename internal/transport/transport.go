package transport

import "io"

type Transport interface {
	Send(file io.Reader, path string) error
}
