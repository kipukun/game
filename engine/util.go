package engine

import (
	"bytes"
	"io"
	"io/fs"

	"github.com/kipukun/game/engine/errors"
)

func head[T any](array []T) T {
	return array[len(array)-1]
}

func asset(f fs.FS, path string) (io.ReadSeekCloser, error) {
	var op errors.Op = "asset"

	bs, err := fs.ReadFile(f, path)
	if err != nil {
		return nil, errors.Error(op, "error reading file", err)
	}
	return &bscloser{bytes.NewReader(bs)}, nil
}
