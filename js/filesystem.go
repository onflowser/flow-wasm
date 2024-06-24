package js

import (
	"fmt"
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
	fmt.Println("ReadFile", source)
	result := resolvePromise(f.target.Call("readFile", source))
	fmt.Println("ReadFile", result)
	// TODO: Handle file not found
	return []byte(result.String()), nil
}

func (f *FileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	fmt.Println("WriteFile", filename, data, perm)
	result := resolvePromise(f.target.Call("writeFile", filename, string(data), perm))
	fmt.Println("WriteFile", result)
	// TODO: Handle errors
	return nil
}

func (f *FileSystem) MkdirAll(path string, perm os.FileMode) error {
	fmt.Println("MkdirAll", path, perm)
	result := resolvePromise(f.target.Call("mkdirAll", path, perm))
	fmt.Println("MkdirAll", result)
	// TODO: Handle errors
	return nil
}

func (f *FileSystem) Stat(path string) (os.FileInfo, error) {
	fmt.Println("Stat", path)
	result := resolvePromise(f.target.Call("stat", path))
	fmt.Println("Stat", result)
	// TODO: Handle errors
	return NewFileInfo(result), nil
}

var _ flowkit.ReaderWriter = &FileSystem{}

func resolvePromise(promise js.Value) js.Value {
	wait := make(chan interface{})
	var result js.Value
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result = args[0]
		wait <- nil
		return nil
	}))
	<-wait
	return result
}
