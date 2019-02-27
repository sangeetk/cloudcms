package item

import (
	"log"
)

const (
	// FileImageType - .jpg, .png, .gif etc
	FileImageType = "image"
	// FileAudioType - .mp3, .wav etc
	FileAudioType = "audio"
	// FileVideoType - .mp4, .avi, .mpg, etc
	FileVideoType = "video"
	// FileDocumentType - .doc, .pdf, .docx, .txt etc
	FileDocumentType = "document"
	// FileOtherType - any other type of file
	FileOtherType = "other"
)

// File fieldz
type File struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	URI   string `json:"uri"`
	Size  int64  `json:"size"`
	Bytes []byte `json:"bytes,omitempty"`
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
