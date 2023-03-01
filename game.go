/*
//  Data structure for representing a game. Implements the ebiten.Game
//  interface (Update in game-update.go, Draw in game-draw.go, Layout
//  in game-layout.go). Provided with a few utilitary functions:
//    - initGame
*/

package main

import (
	"bytes"
	"course/assets"
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	state       int           // Current state of the game
	runnerImage *ebiten.Image // Image with all the sprites of the runners
	runners     [4]Runner     // The four runners used in the game
	f           Field         // The running field
	launchStep  int           // Current step in StateLaunchRun state
	resultStep  int           // Current step in StateResult state
	getTPS      bool          // Help for debug
}

// These constants define the five possible states of the game
const (
	StateWelcomeScreen int = iota // Title screen
	StateChooseRunner             // Player selection screen
	StateLaunchRun                // Countdown before a run
	StateRun                      // Run
	StateResult                   // Results announcement
)

// InitGame builds a new game ready for being run by ebiten
func InitGame() (g Game) {

	// Open the png image for the runners sprites
	img, _, err := image.Decode(bytes.NewReader(assets.RunnerImage))
	if err != nil {
		log.Fatal(err)
	}
	g.runnerImage = ebiten.NewImageFromImage(img)

	// Define game parameters
	start := 50.0
	finish := float64(screenWidth - 50)
	frameInterval := 20

	// Create the runners
	for i := range g.runners {
		interval := 0
		if i == 0 {
			interval = frameInterval
		}
		g.runners[i] = Runner{
			xpos: start, ypos: 50 + float64(i*20),
			maxFrameInterval: interval,
			colorScheme:      0,
		}
	}

	// Create the field
	g.f = Field{
		xstart:   start,
		xarrival: finish,
		chrono:   time.Now(),
	}

	return g
}
