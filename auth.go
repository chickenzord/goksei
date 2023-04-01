package goksei

import (
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/file"
)

type AuthStore gokv.Store

func NewFileAuthStore(dir string) (AuthStore, error) {
	fileStore, err := file.NewStore(file.Options{
		Directory: dir,
		Codec:     encoding.JSON,
	})
	if err != nil {
		return nil, err
	}

	return fileStore, err
}
