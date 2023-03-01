/*
//  Data structure for representing the running field
//  used in the game. Provided with a few utilitary
//  methods:
//    - Draw
//    - Reset
*/

package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Field struct {
	xstart   float64   // Position of the start line for the run
	xarrival float64   // Position of the arrival line for the run
	chrono   time.Time // Time recording of the run duration
}

// Reset allows to reset the field state in order to be able to start a new run
func (f *Field) Reset() {
	f.chrono = time.Now()
}

// Draw displays the field on the screen
func (f *Field) Draw(screen *ebiten.Image, drawChrono bool) {
	ebitenutil.DrawLine(screen, f.xstart, 55, f.xstart, 135, color.White)
	ebitenutil.DrawLine(screen, f.xarrival, 55, f.xarrival, 135, color.White)
	for i := 0; i < 5; i++ {
		ebitenutil.DrawLine(screen, f.xstart, float64(55+i*20), f.xarrival, float64(55+i*20), color.White)
	}
	if drawChrono {
		s, ms := GetSeconds(time.Since(f.chrono).Milliseconds())
		ebitenutil.DebugPrintAt(screen, fmt.Sprint(s, ":", ms), screenWidth/2-25, 10)
	}
}
