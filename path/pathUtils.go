// Package path interfaces and utilities for filesystem operations
package path

import (
	"io/fs"
	"os"
)

var PathUtilsInterface PathInterface
var OSUtilsInterface OSInterface

// PathInterface defines the contract for path-related operations.
// Implementations must provide methods for retrieving various
// application paths and handling filesystem operations.
type PathInterface interface {
	GetDefaultPath() (string, error)
	GetLogFilePath(fileName string) (string, error)
	GetConfigFilePath() (string, error)
}

// OSInterface defines the contract for OS-level filesystem operations.
// Provides a mockable interface for testing filesystem interactions
// and ensures consistent behavior across different platforms.
type OSInterface interface {
	UserHomeDir() (string, error)
	Stat(name string) (fs.FileInfo, error)
	IsNotExist(err error) bool
	Mkdir(name string, perm fs.FileMode) error
	MkdirAll(name string, perm fs.FileMode) error
	OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error)
	Open(name string) (*os.File, error)
	ReadFile(path string) ([]byte, error)
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

// MkdirAll creates a directory and any necessary parent directories with specified permissions.
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

// ReadFile reads the entire contents of a file into a byte slice.
func (o OSUtils) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes data to a file, creating it if necessary, with specified permissions.
func (o OSUtils) WriteFile(name string, content []byte, perm fs.FileMode) error {
	return os.WriteFile(name, content, perm)
}
