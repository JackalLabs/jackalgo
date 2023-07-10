package types

import (
	"bytes"
	"time"
)

type Details struct {
	Name         string
	LastModified time.Time
	FileType     string
	Size         int64
}

type File struct {
	Details Details
	Buffer  *bytes.Buffer
}

func (f File) Read(p []byte) (n int, err error) {
	return f.Buffer.Read(p)
}

func (f File) Name() string {
	return f.Details.Name
}
