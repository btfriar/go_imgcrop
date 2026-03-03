# go_imgcrop

A Go package for intelligently cropping and resizing images to specific dimensions or aspect ratios.

[![Go Reference](https://pkg.go.dev/badge/github.com/friar/go_imgcrop.svg)](https://pkg.go.dev/github.com/friar/go_imgcrop)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- Supports multiple image formats (JPEG, PNG, GIF, BMP)
- Center-weighted cropping to preserve image subjects
- High-quality resizing with configurable quality levels
- Configurable anchor points for crop alignment
- Production-ready with comprehensive error handling

## Installation

```bash
go get github.com/friar/go_imgcrop
```

## Quick Start

```go
package main

import (
    "os"
    "image/jpeg"
    "github.com/friar/go_imgcrop"
)

func main() {
    // Open image file
    file, _ := os.Open("photo.jpg")
    defer file.Close()

    // Crop and resize to 800x600
    result, err := imgcrop.CropAndResize(file, imgcrop.Options{
        Width:   800,
        Height:  600,
        Quality: imgcrop.QualityHigh,
    })
    if err != nil {
        panic(err)
    }

    // Save result
    out, _ := os.Create("thumbnail.jpg")
    defer out.Close()
    jpeg.Encode(out, result.Image, &jpeg.Options{Quality: 90})
}
```

## License

MIT License - see [LICENSE](LICENSE) file.

---

# Why
this is a package create to use with summercampscout.com camp image uploading.
Feel free to use it and let me know if you have issues.

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestCalculateCropRegion ./...

# Run benchmarks
go test -bench=. ./...

# Run with race detector (good for production code)
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```
