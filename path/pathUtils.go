// Package path provides all path related functions
package path

import (
	"io/fs"
	"os"
)

var PathUtilsInterface PathInterface
var OSUtilsInterface OSInterface

// PathInterface defines the interface for path-related operations
type PathInterface interface {
	GetDefaultPath() (string, error)
	GetLogFilePath(fileName string) (string, error)
	GetConfigFilePath() (string, error)
}

// OSInterface defines the interface for OS-related operations
type OSInterface interface {
	UserHomeDir() (string, error)
	Stat(name string) (fs.FileInfo, error)
	IsNotExist(err error) bool
	Mkdir(name string, perm fs.FileMode) error
	MkdirAll(name string, perm fs.FileMode) error
	OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error)
	Open(name string) (*os.File, error)
	WriteFile(name string, content []byte, perm fs.FileMode) error
}

// PathUtils implements the PathInterface
type PathUtils struct{}

// OSUtils implements the OSInterface
type OSUtils struct{}

// UserHomeDir returns the home directory of the current user
func (o OSUtils) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

// Stat returns the FileInfo structure describing file
func (o OSUtils) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

// IsNotExist returns a boolean indicating whether the error is known to report that a file or directory does not exist
func (o OSUtils) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// Mkdir creates a new directory with the specified name and permission bits
func (o OSUtils) Mkdir(name string, perm fs.FileMode) error {
	return os.Mkdir(name, perm)
}

func (o OSUtils) MkdirAll(name string, perm fs.FileMode) error {
	return os.MkdirAll(name, perm)
}

// OpenFile is the generalized open call; most users will use Open or Create instead
func (o OSUtils) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

// Open opens the named file for reading
func (o OSUtils) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (o OSUtils) WriteFile(name string, content []byte, perm fs.FileMode) error {
	return os.WriteFile(name, content, perm)
}
