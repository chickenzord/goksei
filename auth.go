package goksei

import (
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/file"
)

// AuthStore provides an interface for persisting authentication tokens.
// It extends the gokv.Store interface to allow different storage backends
// for caching JWT tokens between client sessions.
type AuthStore gokv.Store

// NewFileAuthStore creates a new file-based authentication token store.
// The dir parameter specifies the directory where token files will be stored.
// Each username will have its own file within this directory.
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
