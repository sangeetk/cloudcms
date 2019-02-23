package item

import (
	"log"
)

// File fieldz
type File struct {
	Name  string `json:"name"`
	URI   string `json:"uri"`
	Size  int64  `json:"size"`
	Bytes []byte `json:"bytes"`
}

// Write byte array from p to f.Bytes
func (f *File) Write(p []byte) (int, error) {
	f.Bytes = append(f.Bytes, p...)
	return len(p), nil
}

//  Read byte array from f.Bytes and copy to p
func (f *File) Read(p []byte) (int, error) {
	log.Println("item.File.Read()", len(p))

	return len(p), nil
}
