package wo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type reader struct {
}

func (r *reader) Read([]byte) (int, error) {
	return 0, nil
}

type readCloser struct {
	closed bool
}

func (r *readCloser) Read([]byte) (int, error) {
	return 0, nil
}

func (r *readCloser) Close() error {
	r.closed = true
	return nil
}

func TestReadCloserWrapper_Close_Closer(t *testing.T) {
	r := &readCloser{
		closed: false,
	}

	wrapped := &readCloserWrapper{r}
	wrapped.Close()

	assert.True(t, r.closed)
}

func TestReadCloserWrapper_Close_Reader(t *testing.T) {
	r := &reader{}

	wrapped := &readCloserWrapper{r}
	err := wrapped.Close()

	assert.Nil(t, err)
}
