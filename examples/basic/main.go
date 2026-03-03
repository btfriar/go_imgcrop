// Package main demonstrates basic usage of the imgcrop package.
//
// This example shows how to:
// 1. Open an image file
// 2. Crop and resize it to specific dimensions
// 3. Save the result to a new file
//
// Run this example:
//
//	go run examples/basic/main.go input.jpg output.jpg 800 600
package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// TODO: Uncomment this import once the package is ready
	// "github.com/friar/go_imgcrop"
)

func main() {
	// Parse command line arguments
	if len(os.Args) != 5 {
		fmt.Println("Usage: go run main.go <input> <output> <width> <height>")
		fmt.Println("Example: go run main.go photo.jpg thumbnail.jpg 800 600")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	width, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Invalid width: %v\n", err)
		os.Exit(1)
	}

	height, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Printf("Invalid height: %v\n", err)
		os.Exit(1)
	}

	// Process the image
	if err := processImage(inputPath, outputPath, width, height); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created %s (%dx%d)\n", outputPath, width, height)
}

func processImage(inputPath, outputPath string, width, height int) error {
	// Open the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// TODO: Uncomment and implement once the package is ready
	/*
		// Create options for the crop operation
		opts := imgcrop.Options{
			Width:   width,
			Height:  height,
			Quality: imgcrop.QualityHigh,
			Anchor:  imgcrop.AnchorCenter,
		}

		// Alternative: use the fluent API
		// opts := imgcrop.DefaultOptions().
		//     WithDimensions(width, height).
		//     WithQuality(imgcrop.QualityHigh)

		// Process the image
		result, err := imgcrop.CropAndResize(inputFile, opts)
		if err != nil {
			return fmt.Errorf("failed to process image: %w", err)
		}

		// Log some information about the processing
		fmt.Printf("Original: %dx%d (%s)\n",
			result.OriginalWidth,
			result.OriginalHeight,
			result.Format)
		fmt.Printf("Cropped to: %dx%d\n",
			result.CroppedWidth,
			result.CroppedHeight)
		fmt.Printf("Final size: %dx%d\n",
			result.FinalWidth,
			result.FinalHeight)
	*/

	// Create the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Encode the result based on output file extension
	ext := strings.ToLower(filepath.Ext(outputPath))

	// TODO: Uncomment once the package is ready
	/*
		switch ext {
		case ".jpg", ".jpeg":
			err = jpeg.Encode(outputFile, result.Image, &jpeg.Options{Quality: 90})
		case ".png":
			err = png.Encode(outputFile, result.Image)
		default:
			return fmt.Errorf("unsupported output format: %s", ext)
		}
	*/

	// Placeholder to avoid unused import errors
	_ = jpeg.DefaultQuality
	_ = png.DefaultCompression
	_ = ext

	if err != nil {
		return fmt.Errorf("failed to encode output: %w", err)
	}

	return nil
}

// Example of using CropToAspectRatio for when you want to maintain
// original resolution but change the shape.
//
// TODO: Uncomment once the package is ready
/*
func cropToAspectRatioExample(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Crop to 16:9 aspect ratio without resizing
	result, err := imgcrop.CropToAspectRatio(inputFile, 16, 9)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return jpeg.Encode(outputFile, result.Image, &jpeg.Options{Quality: 95})
}
*/
