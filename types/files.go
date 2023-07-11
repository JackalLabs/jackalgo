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
	Details Details       `json:"details"`
	Buffer  *bytes.Buffer `json:"bytes"`
}

func (f File) Read(p []byte) (n int, err error) {
	return f.Buffer.Read(p)
}

func (f File) Name() string {
	return f.Details.Name
}
