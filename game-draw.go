/*
//  Implementation of the Draw method for the Game structure
//  This method is called once at every frame (60 frames per second)
//  by ebiten, juste after calling the Update method (game-update.go)
//  Provided with a few utilitary methods:
//    - DrawLaunch
//    - DrawResult
//    - DrawRun
//    - DrawSelectScreen
//    - DrawWelcomeScreen
*/

package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//Nombre de joueurs connecté
var playerscon string = "0"

//Bool de connexion
var isConnected bool = false

//Tableau de numéro des joueurs (remplacé p0,p1,p2,p3 par de vrais numéros)
var numPlayers = []int{
	0,
	0,
	0,
	0,
}

//Var de mode de jeu 0=JOIN 1=HOST
var gameMode int = 0

//nbr de bot pour un host
var nbrBot int = 0

//canal qui gère l'apparition du SPACE START toutes les 0.5s
var tmp = time.After(1000 * time.Millisecond)

//Valeur de l'appirtion du SPACE START
var varTmp = 0

//Images des carrés de sélection de background
var squares = []*ebiten.Image{
	ebiten.NewImage(15, 15),
	ebiten.NewImage(15, 15),
	ebiten.NewImage(15, 15),
	ebiten.NewImage(15, 15),
}

//Hauteurs des carrés de sélection de background
var heights = []float64{
	30,
	60,
	90,
	120,
}

//Couleurs des carrés de sélection de background
var bgColors = []color.RGBA{
	{174, 97, 255, 255},
	{255, 69, 69, 255},
	{255, 190, 128, 255},
	{125, 156, 71, 255},
	{141, 200, 235, 255},
}

// DrawWelcomeScreen displays the title screen in the game window
func (g *Game) DrawWelcomeScreen(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprint("COURS!!! : HELL X SUNAJ EDITION"),
		screenWidth/2-100,
		screenHeight/2+5,
	)
	//Affiche le nombre de joueurs connectés a la partie
	if isConnected {
		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprint(playerscon, "/4 PLAYERS CONNECTED"),
			screenWidth/2-60,
			screenHeight/2+40,
		)
	} else {
		//Affiche avant la connexion et selection des différents modes

		//Affichage PRESS SPACE TO START
		select {
		case <-tmp:
			varTmp = (varTmp + 1) % 2
			tmp = time.After(500 * time.Millisecond)
		default:

		}
		if varTmp == 1 {
			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprint("Press SPACE to play"),
				screenWidth/2-60,
				screenHeight/2+40,
			)
		}

		//Selection gamemode
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			gameMode = (gameMode + 1) % 2
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			gameMode = (gameMode - 1) % 2
			if gameMode < 0 {
				gameMode = 1
			}
		}
		if gameMode == 0 {
			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprint("< JOIN >"),
				screenWidth/2-26,
				10,
			)
		} else if gameMode == 1 {
			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprint("< HOST >"),
				screenWidth/2-26,
				10,
			)

			//Selection Bot pour Hosting
			if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
				if nbrBot != 3 {
					nbrBot = (nbrBot + 1) % 4
				}
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
				if nbrBot != 0 {
					nbrBot = (nbrBot - 1) % 4
				}
			}
			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprint("BOTS :"),
				screenWidth/2-26,
				35,
			)
			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprint(nbrBot),
				screenWidth/2+25,
				35,
			)
			if nbrBot != 3 {
				ebitenutil.DebugPrintAt(
					screen,
					fmt.Sprintf("+"),
					screenWidth/2+25,
					22,
				)
			}
			if nbrBot != 0 {
				ebitenutil.DebugPrintAt(
					screen,
					fmt.Sprintf("_"),
					screenWidth/2+25,
					43,
				)
			}
		}

	}

	if cheatValue == 3 {
		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("CHEAT CODE 'FAUCOU' ACTIVATED --- press 'R' to disable"),
			5,
			screenHeight-15,
		)
	}

	//Option d'image ()
	var options = []*ebiten.DrawImageOptions{
		{},
		{},
		{},
		{},
	}

	for i, j := range heights {
		squares[i].Fill(bgColors[i])
		options[i].GeoM.Translate(10, j)
		screen.DrawImage(squares[i], options[i])
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			if (30 >= mx) && (mx >= 10) && (my >= int(j)) && (int(j)+15 >= my) {
				bgColor = bgColors[i]
				println(bgColors)
				bgColors[i], bgColors[4] = bgColors[4], bgColors[i]
			}
		}
	}

}

// DrawSelectScreen displays the runner selection screen in the game window
func (g *Game) DrawSelectScreen(screen *ebiten.Image) {
	if g.runners[0].colorSelected {
		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprint("Player Selected"),
			screenWidth/2-60,
			20,
		)
	} else {
		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprint("Select your player"),
			screenWidth/2-60,
			20,
		)
	}

	xStep := (screenWidth - 100) / 8
	xPadding := (xStep - 32) / 2
	yPos := (screenHeight - 32) / 2
	for i := 0; i < 8; i++ {
		options := &ebiten.DrawImageOptions{}
		xPos := 50 + i*xStep + xPadding
		options.GeoM.Translate(float64(xPos), float64(yPos))
		screen.DrawImage(g.runnerImage.SubImage(image.Rect(0, i*32, 32, i*32+32)).(*ebiten.Image), options)
	}
	for i := range g.runners {
		g.runners[i].DrawSelection(screen, xStep, i)
	}
}

// DrawLaunch displays the countdown before a run in the game window
func (g *Game) DrawLaunch(screen *ebiten.Image) {
	if g.launchStep > 1 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprint(5-g.launchStep), screenWidth/2-10, 10)
	}
}

// DrawRun displays the current state of the run in the game window
func (g *Game) DrawRun(screen *ebiten.Image, drawChrono bool) {
	g.f.Draw(screen, drawChrono)
	for i := range g.runners {
		g.runners[i].Draw(screen)
	}
}

// DrawResult displays the results of the run in the game window
func (g *Game) DrawResult(screen *ebiten.Image) {
	ranking := [4]int{-1, -1, -1, -1}
	for i := range g.runners {
		rank := 0
		for j := range g.runners {
			if g.runners[i].runTime > g.runners[j].runTime {
				rank++
			}
		}
		for ranking[rank] != -1 {
			rank++
		}
		ranking[rank] = i
	}

	for i := 1; i < g.resultStep && i <= 4; i++ {
		//s, ms := GetSeconds(g.runners[ranking[i-1]].runTime.Milliseconds())
		ebitenutil.DebugPrintAt(screen, fmt.Sprint(i, ". P", numPlayers[ranking[i-1]], "     ", timerunners[numPlayers[ranking[i-1]]-1]), screenWidth/2-40, 55+ranking[i-1]*20)
	}

	if g.resultStep > 4 {
		ebitenutil.DebugPrintAt(screen, "Press SPACE to restart", screenWidth/2-60, 10)
	}
}

// Draw is the main drawing function of the game. It is called by ebiten at
// each frame (60 times per second) just after calling Update (game-update.go)
// Depending of the current state of the game it calls the above utilitary
// function to draw what is needed in the game window

//Couleur du background par défaut. Cette couleur est modifiable par l'utilisateur..
var bgColor = color.RGBA{141, 200, 235, 255}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)

	if g.getTPS {
		ebitenutil.DebugPrint(screen, fmt.Sprint(ebiten.CurrentTPS()))
	}

	switch g.state {
	case StateWelcomeScreen:
		g.DrawWelcomeScreen(screen)
	case StateChooseRunner:
		g.DrawSelectScreen(screen)
	case StateLaunchRun:
		g.DrawLaunch(screen)
		g.DrawRun(screen, false)
	case StateRun:
		g.DrawRun(screen, true)
	case StateResult:
		g.DrawResult(screen)
		g.DrawRun(screen, false)
	}
}
