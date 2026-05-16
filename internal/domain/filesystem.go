package domain

import (
	"io"
	"os"
)

// FileSystem defines the interface for file system operations.
type FileSystem interface {
	Exists(path string) bool
	IsDir(path string) bool
	CreateDirectory(path string) error
	CopyFile(src, dst string) error
	ReadDir(path string) ([]os.DirEntry, error)
	Stat(path string) (os.FileInfo, error)
	Open(path string) (io.ReadCloser, error)
	Create(path string) (io.WriteCloser, error)
}
