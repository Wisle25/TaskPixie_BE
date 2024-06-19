package file_statics

// FileUpload handling file uploading, manipulating (like adding watermark), removing, etc
type FileUpload interface {
	// UploadFile before uploading file, it will generate a new name.
	// Receiving the buffer and extension file that want to be uploaded.
	// Returning uploaded link.
	UploadFile(buffer []byte, extension string) string

	// GetFile Getting the object as buffer
	GetFile(fileName string) []byte

	// RemoveFile deleting specified file by its link
	// Do nothing if it's really not existed
	RemoveFile(oldFileLink string)
}
