package frame

import (
	"github.com/rjkroege/edwood/draw"
	"github.com/rjkroege/edwood/theme"
	"image"
	"log"
)

// drawTick is a centralised function to render a tick with consistent logic.
func (f *frameimpl) drawTick(target draw.Image, tickColor draw.Image, rect image.Rectangle) {
	target.Draw(rect, tickColor, nil, image.Point{})
}

// InitTick initialises the tick with consistent dark mode logic.
// TODO(rjk): doesn't appear to need to be exposed publically.
func (f *frameimpl) InitTick() {
	if f.cols[ColBack] == nil || f.display == nil {
		return
	}

	// Scaling factor
	f.tickscale = f.display.ScaleSize(1)
	b := f.display.ScreenImage()
	ft := f.font

	// Free existing tickimage if any
	if f.tickimage != nil {
		f.tickimage.Free()
	}

	height := ft.Height()

	var err error
	f.tickimage, err = f.display.AllocImage(image.Rect(0, 0, f.tickscale*frtickw, height), b.Pix(), false, draw.Transparent)
	if err != nil {
		log.Printf("InitTick: Failed to allocate tickimage: %v\n", err)
		return
	}

	var backgroundColor draw.Color
	if theme.IsDarkMode() {
		backgroundColor = theme.BackgroundColor
	} else {
		backgroundColor = theme.White
	}

	f.tickback, err = f.display.AllocImage(f.tickimage.R(), b.Pix(), false, backgroundColor)
	if err != nil {
		log.Printf("InitTick: Failed to allocate tickback image: %v\n", err)
		f.tickimage.Free()
		f.tickimage = nil
		return
	}

	// Draw the background of the tick
	f.tickback.Draw(f.tickback.R(), f.cols[ColBack], nil, image.Point{})

	// Clear the tick image with transparency
	f.tickimage.Draw(f.tickimage.R(), f.display.Transparent(), nil, image.Pt(0, 0))
	// vertical line
	f.tickimage.Draw(image.Rect(f.tickscale*(frtickw/2), 0, f.tickscale*(frtickw/2+1), height), f.display.Opaque(), nil, image.Pt(0, 0))
	// box on each end
	f.tickimage.Draw(image.Rect(0, 0, f.tickscale*frtickw, f.tickscale*frtickw), f.display.Opaque(), nil, image.Pt(0, 0))
	f.tickimage.Draw(image.Rect(0, height-f.tickscale*frtickw, f.tickscale*frtickw, height), f.display.Opaque(), nil, image.Pt(0, 0))

	// Flush the display buffer to ensure it's drawn
	err = f.display.Flush()
	if err != nil {
		log.Printf("InitTick: Failed to flush display: %v\n", err)
		f.tickimage = nil // Disable the tick functionality to avoid issues
	}
}
