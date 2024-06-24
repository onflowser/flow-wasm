package js

import (
	"io/fs"
	"os"
	"syscall/js"
	"time"
)

type FileInfo struct {
	target js.Value
}

func NewFileInfo(target js.Value) *FileInfo {
	return &FileInfo{
		target: target,
	}
}

func (f *FileInfo) Name() string {
	return f.target.Get("name").String()
}

func (f *FileInfo) Size() int64 {
	return int64(f.target.Get("size").Int())
}

func (f *FileInfo) Mode() fs.FileMode {
	return fs.FileMode(f.target.Get("mode").Int())
}

func (f *FileInfo) ModTime() time.Time {
	return time.UnixMilli(int64(f.target.Get("modTime").Int()))
}

func (f *FileInfo) IsDir() bool {
	return f.target.Get("isDir").Bool()
}

func (f *FileInfo) Sys() any {
	return nil
}

var _ os.FileInfo = &FileInfo{}
