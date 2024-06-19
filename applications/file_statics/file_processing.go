package file_statics

type ConvertTo int8

const (
	WEBP = iota
	JPG
)

type FileProcessing interface {
	CompressImage(buffer []byte, to ConvertTo) ([]byte, string)
	AddWatermark(buffer []byte) []byte
}
