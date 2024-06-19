package file_statics

import (
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/wisle25/be-template/applications/file_statics"
	"os"
	"path/filepath"
)

type VipsFileProcessing struct {
}

func NewVipsFileProcessing() file_statics.FileProcessing {
	vips.Startup(nil)
	return &VipsFileProcessing{}
}

func (v *VipsFileProcessing) CompressImage(buffer []byte, to file_statics.ConvertTo) ([]byte, string) {
	if buffer == nil {
		return nil, ""
	}

	var result []byte
	var extension string
	var err error

	image, err := vips.NewImageFromBuffer(buffer)
	if err != nil {
		panic(fmt.Errorf("vips: new compress image: %v", err))
	}

	switch to {
	case file_statics.JPG:
		options := vips.NewJpegExportParams()
		options.Quality = 40
		options.StripMetadata = true
		result, _, err = image.ExportJpeg(options)
		if err != nil {
			panic(fmt.Errorf("vips: compress image jpeg: %v", err))
		}

		extension = ".jpg"
	case file_statics.WEBP:
		options := vips.NewWebpExportParams()
		options.Quality = 40
		options.StripMetadata = true
		options.Lossless = false

		result, _, err = image.ExportWebp(options)
		if err != nil {
			panic(fmt.Errorf("vips: compress image webp: %v", err))
		}

		extension = ".webp"
	}

	return result, extension
}

func (v *VipsFileProcessing) AddWatermark(buffer []byte) []byte {
	// Open original image
	originalImage, err := vips.NewImageFromBuffer(buffer)
	if err != nil {
		panic(fmt.Errorf("add_watermark_err: opening original image: %v", err))
	}
	defer originalImage.Close()

	// Open watermark image
	rootDir, _ := os.Getwd()
	watermarkImage, err := vips.NewImageFromFile(filepath.Join(rootDir, "resources", "watermark.png"))
	if err != nil {
		panic(fmt.Errorf("add_watermark_err: opening watermark image: %v", err))
	}
	defer watermarkImage.Close()

	// Resize watermark image to fit the original image
	err = watermarkImage.ResizeWithVScale(
		float64(originalImage.Width())/float64(watermarkImage.Width()),
		float64(originalImage.Height())/float64(watermarkImage.Height()),
		vips.KernelLanczos3,
	)
	if err != nil {
		panic(fmt.Errorf("add_watermark_err: resizing watermark: %v", err))
	}

	// Composite
	err = originalImage.Composite(watermarkImage, vips.BlendModeAdd, 0, 0)
	if err != nil {
		panic(fmt.Errorf("add_watermark_err: compositing watermark: %v", err))
	}

	// Get the buffer of the result
	resultBuffer, _, _ := originalImage.ExportNative()

	return resultBuffer
}
