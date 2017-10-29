package persistence

import (
	"errors"
	"os"
	"path"
)

type (
	// Saver is an interface for I/O operations.
	Saver struct {
		file *os.File
	}
)

// NewSaver allocates and returns a new Saver.
func NewSaver(filePath string) *Saver {
	file, err := createFile(filePath)

	if err != nil {
		return nil
	}

	return &Saver{
		file: file,
	}
}

// createFile creates file, and returns error in case of any.
func createFile(filePath string) (*os.File, error) {
	if filePath == "" {
		return nil, errors.New("saver: File path should not be empty")
	}
	dirPath := path.Dir(filePath)
	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {
		os.Mkdir(dirPath, 0755)
	}

	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)
}

// Read reads and returns saved data.
func (s *Saver) Read() ([]byte, error) {
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
func (s *Saver) Save(data []byte) error {
	s.file.Seek(0, 0)
	s.file.Truncate(0)

	_, err := s.file.WriteAt(data, 0)

	if err != nil {
		return err
	}

	return nil
}

// Close closes the underlying layer used for saving.
func (s *Saver) Close() error {
	return s.file.Close()
}
