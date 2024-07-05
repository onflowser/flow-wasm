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
	// If we don't explicitly convert os.FileMode to uint32, the call will fail due to serialization errors.
	_, err := parseResult(resolvePromise(f.target.Call("writeFile", filename, string(data), uint32(perm))))

	if err != nil {
		return err
	}

	return nil
}

func (f *FileSystem) MkdirAll(path string, perm os.FileMode) error {
	// If we don't explicitly convert os.FileMode to uint32, the call will fail due to serialization errors.
	_, err := parseResult(resolvePromise(f.target.Call("mkdirAll", path, uint32(perm))))

	if err != nil {
		return err
	}

	return nil
}

func (f *FileSystem) Stat(path string) (os.FileInfo, error) {
	value, err := parseResult(resolvePromise(f.target.Call("stat", path)))

	if err != nil {
		return nil, err
	}

	return NewFileInfo(value), nil
}

var _ flowkit.ReaderWriter = &FileSystem{}
