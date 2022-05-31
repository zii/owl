package android

import (
	"io/fs"

	"github.com/psanford/memfs"
)

var (
	mfs *memfs.FS
)

func raise(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	mfs = memfs.New()
	err := mfs.MkdirAll("screen", 0777)
	raise(err)
}

func GetFS() fs.FS {
	return mfs
}
