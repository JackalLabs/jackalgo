package types

import (
	"bytes"
	"time"
)

type Details struct {
	Name         string    `json:"name"`
	LastModified time.Time `json:"lastModified"`
	FileType     string    `json:"type"`
	Size         int64     `json:"size"`
}

type File struct {
	Details Details
	buffer  []byte
}

func NewFile(data []byte, details Details) *File {
	f := File{
		buffer:  data,
		Details: details,
	}

	return &f
}

func (f File) Buffer() *bytes.Buffer {
	b := make([]byte, len(f.buffer))
	copy(b, f.buffer)
	return bytes.NewBuffer(b)
}

func (f File) Read(p []byte) (n int, err error) {
	return f.Buffer().Read(p)
}

func (f File) Name() string {
	return f.Details.Name
}
