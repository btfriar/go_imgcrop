// Package imgcrop provides functionality for cropping and resizing images
// to specific aspect ratios and dimensions while preserving image quality.
//
// The package supports multiple image formats (JPEG, PNG, GIF, BMP) and uses
// center-weighted cropping to maintain the most important parts of an image.
//
// Basic usage:
//
//	result, err := imgcrop.CropAndResize(reader, imgcrop.Options{
//		Width:  800,
//		Height: 600,
//	})
//
// For aspect ratio based cropping:
//
//	result, err := imgcrop.CropToAspectRatio(reader, 16, 9)
package imgcrop

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"golang.org/x/image/draw"
)

// Result holds the processed image along with metadata about the operation.
// This provides consumers with both the image data and information about
// what transformations were applied.
type Result struct {
	// Image is the processed image ready for encoding
	Image image.Image

	// OriginalWidth is the width of the source image before processing
	OriginalWidth int

	// OriginalHeight is the height of the source image before processing
	OriginalHeight int

	// CroppedWidth is the width after cropping but before final resize
	CroppedWidth int

	// CroppedHeight is the height after cropping but before final resize
	CroppedHeight int

	// FinalWidth is the width of the output image
	FinalWidth int

	// FinalHeight is the height of the output image
	FinalHeight int

	// Format is the detected format of the source image (e.g., "jpeg", "png")
	Format string
}

// CropAndResize takes an image from a reader and crops/resizes it according to
// the provided options. It returns a Result containing the processed image
// and metadata about the transformations applied.
func CropAndResize(r io.Reader, opts Options) (*Result, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	img, format, err := decodeImage(r)
	if err != nil {
		return nil, err
	}

	originalBounds := img.Bounds()
	targetRatio := float64(opts.Width) / float64(opts.Height)

	cropRegion := calculateCropRegion(originalBounds, targetRatio, opts.Anchor)
	cropped := cropImage(img, cropRegion)
	final := resizeImage(cropped, opts.Width, opts.Height, opts.Quality)

	return &Result{
		Image:          final,
		OriginalWidth:  originalBounds.Dx(),
		OriginalHeight: originalBounds.Dy(),
		CroppedWidth:   cropRegion.Dx(),
		CroppedHeight:  cropRegion.Dy(),
		FinalWidth:     opts.Width,
		FinalHeight:    opts.Height,
		Format:         format,
	}, nil
}

// CropToAspectRatio crops an image to match the specified aspect ratio
// without resizing. Useful when you want to maintain original resolution
// but change the shape.
func CropToAspectRatio(r io.Reader, widthRatio, heightRatio int) (*Result, error) {
	if widthRatio <= 0 || heightRatio <= 0 {
		return nil, ErrInvalidAspectRatio
	}

	img, format, err := decodeImage(r)
	if err != nil {
		return nil, err
	}

	originalBounds := img.Bounds()
	targetRatio := float64(widthRatio) / float64(heightRatio)

	cropRegion := calculateCropRegion(originalBounds, targetRatio, AnchorCenter)
	cropped := cropImage(img, cropRegion)

	return &Result{
		Image:          cropped,
		OriginalWidth:  originalBounds.Dx(),
		OriginalHeight: originalBounds.Dy(),
		CroppedWidth:   cropRegion.Dx(),
		CroppedHeight:  cropRegion.Dy(),
		FinalWidth:     cropRegion.Dx(),
		FinalHeight:    cropRegion.Dy(),
		Format:         format,
	}, nil
}

// decodeImage reads and decodes an image from a reader.
// It returns the decoded image and the format string (e.g., "jpeg", "png").
func decodeImage(r io.Reader) (image.Image, string, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, "", &DecodeError{Format: format, Err: err}
	}
	return img, format, nil
}

// calculateCropRegion determines the optimal rectangle to crop from the source
// image to achieve the target aspect ratio. The anchor parameter controls which
// part of the image is preserved when cropping.
func calculateCropRegion(bounds image.Rectangle, targetRatio float64, anchor Anchor) image.Rectangle {
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()
	srcRatio := float64(srcWidth) / float64(srcHeight)

	var cropX, cropY, newWidth, newHeight int

	if srcRatio > targetRatio {
		// Image is wider than target: crop left/right
		newHeight = srcHeight
		newWidth = int(float64(srcHeight) * targetRatio)
		excess := srcWidth - newWidth
		switch anchor {
		case AnchorLeft:
			cropX = 0
		case AnchorRight:
			cropX = excess
		default:
			cropX = excess / 2
		}
		cropY = 0
	} else if srcRatio < targetRatio {
		// Image is taller than target: crop top/bottom
		newWidth = srcWidth
		newHeight = int(float64(srcWidth) / targetRatio)
		excess := srcHeight - newHeight
		switch anchor {
		case AnchorTop:
			cropY = 0
		case AnchorBottom:
			cropY = excess
		default:
			cropY = excess / 2
		}
		cropX = 0
	} else {
		return bounds
	}

	return image.Rect(
		bounds.Min.X+cropX,
		bounds.Min.Y+cropY,
		bounds.Min.X+cropX+newWidth,
		bounds.Min.Y+cropY+newHeight,
	)
}

// cropImage extracts a rectangular region from the source image.
func cropImage(img image.Image, region image.Rectangle) image.Image {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	if si, ok := img.(subImager); ok {
		return si.SubImage(region)
	}
	dst := image.NewRGBA(image.Rect(0, 0, region.Dx(), region.Dy()))
	draw.Draw(dst, dst.Bounds(), img, region.Min, draw.Src)
	return dst
}

// resizeImage scales an image to the exact target dimensions.
func resizeImage(img image.Image, width, height int, quality Quality) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	var scaler draw.Scaler
	switch quality {
	case QualityLow:
		scaler = draw.NearestNeighbor
	case QualityMedium:
		scaler = draw.BiLinear
	case QualityHigh:
		scaler = draw.CatmullRom
	default:
		scaler = draw.BiLinear
	}

	scaler.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}
