package engine

import (
	"bytes"
	"io"
	"io/fs"
)

func head[T any](array []T) T {
	return array[len(array)-1]
}

func asset(f fs.FS, path string) (io.ReadSeekCloser, error) {
	bs, err := fs.ReadFile(f, path)
	if err != nil {
		return nil, err
	}
	return &bscloser{bytes.NewReader(bs)}, nil
}
