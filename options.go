package imgcrop

// Options configures how an image should be cropped and resized.
// Use the functional options pattern or direct struct initialization.
//
// Example:
//
//	opts := imgcrop.Options{
//		Width:   800,
//		Height:  600,
//		Quality: imgcrop.QualityHigh,
//	}
type Options struct {
	Width   int
	Height  int
	Quality Quality
	Anchor  Anchor
}

// Quality represents the resampling quality level for resizing operations.
type Quality int

const (
	QualityLow Quality = iota
	QualityMedium
	QualityHigh
)

func (q Quality) String() string {
	switch q {
	case QualityLow:
		return "low"
	case QualityMedium:
		return "medium"
	case QualityHigh:
		return "high"
	default:
		return "unknown"
	}
}

type Anchor int

const (
	AnchorCenter Anchor = iota
	AnchorTop
	AnchorBottom
	AnchorLeft
	AnchorRight
)

// String returns a human-readable name for the anchor position.
func (a Anchor) String() string {
	switch a {
	case AnchorCenter:
		return "center"
	case AnchorTop:
		return "top"
	case AnchorBottom:
		return "bottom"
	case AnchorLeft:
		return "left"
	case AnchorRight:
		return "right"
	default:
		return "unknown"
	}
}

// Validate checks that the options are valid and returns an error if not.
// This should be called at the start of CropAndResize.
func (o Options) Validate() error {
	if o.Width <= 0 {
		return ErrInvalidWidth
	}
	if o.Height <= 0 {
		return ErrInvalidHeight
	}
	if o.Width > MaxDimension || o.Height > MaxDimension {
		return ErrDimensionsTooLarge
	}
	return nil
}

// DefaultOptions returns sensible default options.
// Users can modify specific fields while keeping good defaults for others.
func DefaultOptions() Options {
	return Options{
		Quality: QualityMedium,
		Anchor:  AnchorCenter,
	}
}

// WithDimensions is a convenience method that returns a copy of options
// with the specified dimensions.
func (o Options) WithDimensions(width, height int) Options {
	o.Width = width
	o.Height = height
	return o
}

// WithQuality is a convenience method that returns a copy of options
// with the specified quality level.
func (o Options) WithQuality(q Quality) Options {
	o.Quality = q
	return o
}

// WithAnchor is a convenience method that returns a copy of options
// with the specified anchor position.
func (o Options) WithAnchor(a Anchor) Options {
	o.Anchor = a
	return o
}
