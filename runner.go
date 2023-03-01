/*
//  Data structure for representing the four runners
//  used in the game. Provided with a few utilitary
//  methods:
//    - CheckArrival
//    - Draw
//    - DrawSelection
//    - ManualChoose
//    - ManualUpdate
//    - RandomChoose
//    - RandomUpdate
//    - Reset
//    - UpdateAnimation
//    - UpdatePos
//    - UpdateSpeed
*/

package main

import (
	"fmt"
	"image"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Runner struct {
	xpos, ypos        float64       // Position of the runner on the screen
	speed             float64       // Current speed of the runner
	framesSinceUpdate int           // Number of frames since last speed update
	maxFrameInterval  int           // Maximum number of frames between two speed updates
	arrived           bool          // Tells if the runner has finished running or not
	runTime           time.Duration // Records the duration of the run for ranking
	image             *ebiten.Image // Current image used to display the runner
	colorScheme       int           // Number of the color scheme of the runner
	colorSelected     bool          // Tells if the color scheme is fixed or not
	animationStep     int           // Current step of the runner animation
	animationFrame    int           // Number of frames since the last animation step
}

// ManualUpdate allows to use the keyboard in order to control a runner
// when the game is in the StateRun state (i.e. during a run)
func (r *Runner) ManualUpdate() {
	r.UpdateSpeed(inpututil.IsKeyJustPressed(ebiten.KeySpace))
	r.UpdatePos()
}

// RandomUpdate allows to randomly control a runner when the game is in
// the StateRun state (i.e. during a run)
func (r *Runner) RandomUpdate() {
	r.UpdateSpeed(rand.Intn(3) == 0)
	r.UpdatePos()
}

// UpdateSpeed sets the speed of a runner. It is used when the game is in
// StateRun state (i.e. during a run)
func (r *Runner) UpdateSpeed(keyPressed bool) {
	if !r.arrived {
		r.framesSinceUpdate++
		if keyPressed {
			r.speed = 1500 / float64(r.framesSinceUpdate*r.framesSinceUpdate*r.framesSinceUpdate)
			if r.speed > 10 {
				r.speed = 10
			}
			r.framesSinceUpdate = 0
		} else if r.framesSinceUpdate > r.maxFrameInterval {
			r.speed = 0
		}
	}
}

// UpdatePos sets the current (x) position of a runner according to the current
// speed and the previous (x) position. It is used when the game is in StateRun
// state (i.e. during a run)

var cheatValue int = 1

func (r *Runner) UpdatePos() {
	if !r.arrived {
		r.xpos += r.speed * float64(cheatValue)
	}
}

// UpdateAnimation determines the next image that should be displayed for a
// runner, depending of whether or not the runner is running, the current
// animationStep and the current animationFrame
func (r *Runner) UpdateAnimation(runnerImage *ebiten.Image) {
	r.animationFrame++
	if r.speed == 0 || r.arrived {
		r.image = runnerImage.SubImage(image.Rect(0, r.colorScheme*32, 32, r.colorScheme*32+32)).(*ebiten.Image)
		r.animationFrame = 0
	} else {
		if r.animationFrame > 1 {
			r.animationStep = r.animationStep%6 + 1
			r.image = runnerImage.SubImage(image.Rect(32*r.animationStep, r.colorScheme*32, 32*r.animationStep+32, r.colorScheme*32+32)).(*ebiten.Image)
			r.animationFrame = 0
		}
	}
}

// ManualChoose allows to use the keyboard for selecting the appearance of a
// runner when the game is in StateChooseRunner state (i.e. at player selection
// screen)

var skinalreaadyselected []int
var skinalreaadyselectedbool bool

func (r *Runner) ManualChoose() (done bool) {
	skinalreaadyselectedbool = false
	for _, i := range skinalreaadyselected {
		if i == r.colorScheme {
			skinalreaadyselectedbool = true
		}
	}
	if !skinalreaadyselectedbool {
		r.colorSelected =
			(!r.colorSelected && inpututil.IsKeyJustPressed(ebiten.KeySpace)) ||
				(r.colorSelected && !inpututil.IsKeyJustPressed(ebiten.KeySpace))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		r.colorSelected = false
		r.colorScheme = (r.colorScheme + 1) % 8
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		r.colorSelected = false
		r.colorScheme = (r.colorScheme + 7) % 8
	}
	return r.colorSelected
}

// RandomChoose allows to randomly select the appearance of a
// runner when the game is in StateChooseRunner state (i.e. at player selection
// screen)
func (r *Runner) RandomChoose() (done bool) {
	if !r.colorSelected {
		r.colorScheme = rand.Intn(8)
	}
	r.colorSelected = true
	return r.colorSelected
}

// CheckArrival allows to test if a runner has passed the arrival line
func (r *Runner) CheckArrival(f *Field) {
	if !r.arrived {
		r.arrived = r.xpos >= f.xarrival
		r.runTime = time.Since(f.chrono)
	} else {
		r.xpos = 750
	}
}

// Reset allows to reset a player state in order to be able to start a new run
func (r *Runner) Reset(f *Field) {
	r.xpos = f.xstart
	r.speed = 0
	r.framesSinceUpdate = 0
	r.arrived = false
	r.animationStep = 0
	r.animationFrame = 0
	r.colorSelected = false
}

// Draw draws a runner on screen at the good position (defined by xpos and ypos)
func (r *Runner) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(r.xpos-16, r.ypos-16)
	screen.DrawImage(r.image, options)
}

// DrawSelection draws the current selection of a runner appearance for the
// player select screen
func (r *Runner) DrawSelection(screen *ebiten.Image, xStep, playerNum int) {
	xMod := 32
	if (playerNum/2)%2 == 0 {
		xMod = -32
	}
	xPadding := (xStep + xMod) / 2
	xPos := 43 + xStep*r.colorScheme + xPadding
	yMod := 32
	if playerNum%2 == 0 {
		yMod = -62
	}
	yPos := (screenHeight + yMod) / 2
	ebitenutil.DebugPrintAt(screen, fmt.Sprint("P", numPlayers[playerNum]), xPos, yPos)
	if r.colorSelected {
		ebitenutil.DebugPrintAt(screen, fmt.Sprint("SELECTED"), xPos, yPos+15)
	}
}
