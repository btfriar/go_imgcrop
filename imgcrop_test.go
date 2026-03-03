package imgcrop

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// Test files in Go use the _test.go suffix and the testing package.
// Run tests with: go test ./...
// Run with verbose output: go test -v ./...
// Run specific test: go test -run TestCalculateCropRegion ./...

// createTestImage is a helper function that creates a solid color image
// for testing purposes. Using helpers keeps tests clean and readable.
//
// WHY create test images programmatically:
// - Tests are self-contained (no external files needed)
// - Easy to create images of exact sizes needed
// - Tests run faster than loading files from disk
// - No need to maintain test fixture files
func createTestImage(width, height int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

// encodeTestImage encodes an image to PNG format in a buffer.
// Returns a bytes.Reader that implements io.Reader.
func encodeTestImage(img image.Image) *bytes.Reader {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return bytes.NewReader(buf.Bytes())
}

// TestCalculateCropRegion tests the crop region calculation logic.
// This is one of the most important functions to test thoroughly.
func TestCalculateCropRegion(t *testing.T) {
	tests := []struct {
		name        string
		bounds      image.Rectangle
		targetRatio float64
		want        image.Rectangle
	}{
		{
			name:        "square to wide (16:9) - should crop top and bottom",
			bounds:      image.Rect(0, 0, 1000, 1000),
			targetRatio: 16.0 / 9.0,
			want:        image.Rect(0, 219, 1000, 781),
		},
		{
			name:        "square to tall (9:16) - should crop left and right",
			bounds:      image.Rect(0, 0, 1000, 1000),
			targetRatio: 9.0 / 16.0,
			want:        image.Rect(219, 0, 781, 1000),
		},
		{
			name:        "wide to square (1:1) - should crop left and right",
			bounds:      image.Rect(0, 0, 1920, 1080),
			targetRatio: 1.0,
			want:        image.Rect(420, 0, 1500, 1080),
		},
		{
			name:        "tall to square (1:1) - should crop top and bottom",
			bounds:      image.Rect(0, 0, 1080, 1920),
			targetRatio: 1.0,
			want:        image.Rect(0, 420, 1080, 1500),
		},
		{
			name:        "same ratio - no crop needed",
			bounds:      image.Rect(0, 0, 1920, 1080),
			targetRatio: 16.0 / 9.0,
			want:        image.Rect(0, 0, 1920, 1080),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateCropRegion(tt.bounds, tt.targetRatio)

			if got != tt.want {
				t.Errorf("calculateCropRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestOptionsValidate tests the options validation logic.
func TestOptionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr error
	}{
		{
			name:    "valid options",
			opts:    Options{Width: 800, Height: 600},
			wantErr: nil,
		},
		{
			name:    "zero width",
			opts:    Options{Width: 0, Height: 600},
			wantErr: ErrInvalidWidth,
		},
		{
			name:    "negative width",
			opts:    Options{Width: -100, Height: 600},
			wantErr: ErrInvalidWidth,
		},
		{
			name:    "zero height",
			opts:    Options{Width: 800, Height: 0},
			wantErr: ErrInvalidHeight,
		},
		{
			name:    "dimensions too large",
			opts:    Options{Width: MaxDimension + 1, Height: 600},
			wantErr: ErrDimensionsTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCropAndResize tests the full crop and resize workflow.
// This is an integration test that exercises the entire pipeline.
func TestCropAndResize(t *testing.T) {

	tests := []struct {
		name       string
		srcWidth   int
		srcHeight  int
		opts       Options
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{
			name:       "basic resize",
			srcWidth:   1000,
			srcHeight:  1000,
			opts:       Options{Width: 500, Height: 500},
			wantWidth:  500,
			wantHeight: 500,
			wantErr:    false,
		},
		{
			name:       "crop and resize landscape to portrait",
			srcWidth:   1920,
			srcHeight:  1080,
			opts:       Options{Width: 400, Height: 600},
			wantWidth:  400,
			wantHeight: 600,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := createTestImage(tt.srcWidth, tt.srcHeight, color.RGBA{255, 0, 0, 255})
			r := encodeTestImage(img)

			result, err := CropAndResize(r, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("CropAndResize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				bounds := result.Image.Bounds()
				gotWidth := bounds.Dx()
				gotHeight := bounds.Dy()

				if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
					t.Errorf("CropAndResize() dimensions = %dx%d, want %dx%d",
						gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
				}
			}
		})
	}
}

// TestDecodeImage tests image decoding with various formats.
func TestDecodeImage(t *testing.T) {

	img := createTestImage(100, 100, color.RGBA{0, 255, 0, 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	decoded, format, err := decodeImage(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("decodeImage() error = %v", err)
	}

	if format != "png" {
		t.Errorf("decodeImage() format = %v, want png", format)
	}

	bounds := decoded.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("decodeImage() size = %dx%d, want 100x100", bounds.Dx(), bounds.Dy())
	}
}

// BenchmarkCropAndResize measures performance of the full pipeline.
// Run with: go test -bench=. ./...
func BenchmarkCropAndResize(b *testing.B) {

	// Setup: create a test image once
	img := createTestImage(1920, 1080, color.RGBA{100, 100, 100, 255})
	var buf bytes.Buffer
	png.Encode(&buf, img)
	data := buf.Bytes()

	opts := Options{Width: 800, Height: 600, Quality: QualityMedium}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(data)
		_, err := CropAndResize(r, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}
