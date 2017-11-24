package wo

import "io"

// ByteReader is a function that gets bytes by name.
type ByteReader func(name string) ([]byte, error)

// AssetReader is a function that get an Reader by name.
type AssetReader func(name string) (io.Reader, error)

// readCloserWrapper wraps a Reader to support Close. A
// call to Close will be propagated to the wrapped value
// if it is supported.
type readCloserWrapper struct {
	// Reader is the wrapped Reader
	io.Reader
}

// Close closes the wrapped value, if it is supported.
func (w *readCloserWrapper) Close() error {
	switch t := w.Reader.(type) {
	case io.Closer:
		return t.Close()
	default:
		return nil
	}
}
