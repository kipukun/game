package engine

import (
	"bytes"
	"io"
	"io/fs"
)

func head[T any](array []T) T {
	return array[len(array)-1]
}

func asset(f fs.FS, path string) io.ReadSeekCloser {
	bs, err := fs.ReadFile(f, path)
	if err != nil {
		panic(err)
	}
	return &bscloser{bytes.NewReader(bs)}
}
