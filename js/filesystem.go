package js

import (
	"github.com/onflow/flowkit/v2"
	"os"
	"syscall/js"
)

type FileSystem struct {
	target js.Value
}

func NewFileSystem(target js.Value) *FileSystem {
	return &FileSystem{
		target,
	}
}

func (f *FileSystem) ReadFile(source string) ([]byte, error) {
	value, err := parseResult(resolvePromise(f.target.Call("readFile", source)))

	if err != nil {
		return nil, err
	}

	return []byte(value.String()), nil
}

func (f *FileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	resolvePromise(f.target.Call("writeFile", filename, string(data), perm))
	// TODO: Handle errors
	return nil
}

func (f *FileSystem) MkdirAll(path string, perm os.FileMode) error {
	resolvePromise(f.target.Call("mkdirAll", path, perm))
	// TODO: Handle errors
	return nil
}

func (f *FileSystem) Stat(path string) (os.FileInfo, error) {
	result := resolvePromise(f.target.Call("stat", path))
	// TODO: Handle errors
	return NewFileInfo(result), nil
}

var _ flowkit.ReaderWriter = &FileSystem{}
