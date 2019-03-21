package persistence

import (
	"errors"
	"os"
	"path"
)

type (
	// Saver is an interface for I/O operations.
	Saver interface {
		Read() ([]byte, error)
		Save(data []byte) error
		Close() error
	}

	// SaveHandler is responsible for saving operations.
	SaveHandler struct {
		file *os.File
	}
)

// NewSaver allocates and returns a new Saver.
func NewSaver(filePath string) (*SaveHandler, error) {
	file, err := createFile(filePath)
	if err != nil {
		return nil, err
	}

	return &SaveHandler{
		file: file,
	}, nil
}

// createFile creates file, and returns error in case of any.
func createFile(filePath string) (*os.File, error) {
	if filePath == "" {
		err := errors.New("saver: File path should not be empty")

		return nil, &os.PathError{Path: filePath, Op: "parse", Err: err}
	}

	dirPath := path.Dir(filePath)
	err := os.Mkdir(dirPath, 0755)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)
}

// Read reads and returns saved data.
func (s *SaveHandler) Read() ([]byte, error) {
	// 4K buffer
	data := make([]byte, 4096)
	n, err := s.file.Read(data)

	if err != nil {
		return nil, err
	}

	// Trimming the data.
	return data[:n], nil
}

// Save saves given data by truncating and overriding the current saved data.
func (s *SaveHandler) Save(data []byte) error {
	_, err := s.file.Seek(0, 0)
	if err != nil {
		return err
	}

	err = s.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = s.file.WriteAt(data, 0)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the underlying layer used for saving.
func (s *SaveHandler) Close() error {
	return s.file.Close()
}
