/*
//  Implementation of the Update method for the Game structure
//  This method is called once at every frame (60 frames per second)
//  by ebiten, juste before calling the Draw method (game-draw.go).
//  Provided with a few utilitary methods:
//    - CheckArrival
//    - ChooseRunners
//    - HandleLaunchRun
//    - HandleResults
//    - HandleWelcomeScreen
//    - Reset
//    - UpdateAnimation
//    - UpdateRunners
*/

package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// HandleWelcomeScreen waits for the player to push SPACE in order to
// start the game

var cheatCheck string

func (g *Game) HandleWelcomeScreen() bool {

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		cheatCheck += "a"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		cheatCheck = ""
		cheatCheck += "f"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		cheatCheck += "u"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		cheatCheck += "c"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		cheatCheck += "o"
	}

	if len(cheatCheck) == 6 {
		if cheatCheck == "faucou" && cheatValue == 1 {
			cheatValue = 3
			println("CHEAT CODE \"FAUCOU\": ON ")
			cheatCheck = ""
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		cheatValue = 1
		println("CHEAT CODE : OFF ")
	}

	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// ChooseRunners loops over all the runners to check which sprite each
// of them selected
//var skins string

func (g *Game) ChooseRunners() (done bool) {
	done = false
	//Selected donne le nombre le joueurs ayant validés la selection de leur skin
	selected := 0
	for i := range g.runners {
		//Pour le joueur, on active le choix manuel
		if i == 0 {
			g.runners[i].ManualChoose()
			//Une fois le runner choisi, on envoi au server notre selection
			out <- strconv.Itoa(g.runners[i].colorScheme) + " " + strconv.FormatBool(g.runners[i].colorSelected)

		} else {
			//Ici, on s'occupe des autres runners. On récupère les infos envoyées par le server.
			skinselected := string(<-in)
			skin, _ := strconv.Atoi(strings.Split(skinselected, " ")[0])
			sel, _ := strconv.ParseBool(strings.Split(skinselected, " ")[1])
			g.runners[i].colorScheme = skin
			g.runners[i].colorSelected = sel
		}
	}
	skinalreaadyselected = nil
	for i := range g.runners {
		if g.runners[i].colorSelected {
			selected += 1
			skinalreaadyselected = append(skinalreaadyselected, g.runners[i].colorScheme)
		}
	}
	//Si les 4 joueurs ont validés la selection, alors on fini cette étape
	if selected == 4 {
		done = true
	}
	return done
}

// HandleLaunchRun countdowns to the start of a run
func (g *Game) HandleLaunchRun() bool {
	if time.Since(g.f.chrono).Milliseconds() > 1000 {
		g.launchStep++
		g.f.chrono = time.Now()
	}
	if g.launchStep >= 5 {
		g.launchStep = 0
		return true
	}
	return false
}

// UpdateRunners loops over all the runners to update each of them
func (g *Game) UpdateRunners() {
	for i := range g.runners {
		if i == 0 {
			//A chaque update envoie de la position au serv
			g.runners[i].ManualUpdate()
			pos := strconv.FormatFloat(g.runners[i].xpos, 'f', 2, 64)
			out <- pos
			//g.runners[i].xpos = race

		} else {
			//A chaque update reception de la position des autres (et mise a jour vitesse pour l'animation)
			pos, _ := strconv.ParseFloat(strings.Split(strings.Split(<-in, "/")[0], ":")[1], 64)
			if g.runners[i].xpos < pos {
				g.runners[i].xpos = pos
				g.runners[i].speed = 1
			} else {
				g.runners[i].speed = 0
			}
		}
	}
}

// CheckArrival loops over all the runners to check which ones are arrived
func (g *Game) CheckArrival() (finished bool) {
	finished = true
	for i := range g.runners {
		g.runners[i].CheckArrival(&g.f)
		finished = finished && g.runners[i].arrived
	}
	return finished
}

// Reset resets all the runners and the field in order to start a new run
func (g *Game) Reset() {
	for i := range g.runners {
		g.runners[i].Reset(&g.f)
	}
	//Reset de toutes les variables de jeu
	playerscon = "0"
	timeChecker = false
	isConnected = false
	numPlayers[0] = 0
	g.f.Reset()
}

// UpdateAnimation loops over all the runners to update their sprite
func (g *Game) UpdateAnimation() {
	for i := range g.runners {
		g.runners[i].UpdateAnimation(g.runnerImage)
	}
}

// HandleResults computes the resuls of a run and prepare them for
// being displayed
var timeChecker = false

func (g *Game) HandleResults() bool {
	if !timeChecker {
		s, ms := GetSeconds(g.runners[0].runTime.Milliseconds())
		out <- strconv.FormatInt(s, 10) + "." + strconv.FormatInt(ms, 10)
		timeChecker = true
	}
	if time.Since(g.f.chrono).Milliseconds() > 1000 || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resultStep++
		g.f.chrono = time.Now()
	}
	if g.resultStep >= 4 && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resultStep = 0
		return true
	}
	return false
}

// Update is the main update function of the game. It is called by ebiten
// at each frame (60 times per second) just before calling Draw (game-draw.go)
// Depending of the current state of the game it calls the above utilitary
// function and then it may update the state of the game

var in = make(chan string)
var out = make(chan string)

func (g *Game) Update() error {
	switch g.state {
	case StateWelcomeScreen:
		//Gestion du menu d'acceuil (Gamemode / Gestion Client)
		done := g.HandleWelcomeScreen()
		if done {
			if !isConnected {
				if gameMode == 1 {
					go Server()
					time.Sleep(100 * time.Millisecond)
					go Client()
					for i := 0; i < nbrBot; i++ {
						go ClientBot()
					}
				} else {
					go Client()
				}
				isConnected = true
			}
		}
		if isConnected {
			select {
			case next := <-in:
				if next == "next" {
					g.state++
				}
			default:
			}
		}

	case StateChooseRunner:
		done := g.ChooseRunners()
		if done {
			g.UpdateAnimation()
			g.state++
		}
	case StateLaunchRun:
		done := g.HandleLaunchRun()
		if done {
			g.state++
		}
	case StateRun:
		g.UpdateRunners()
		finished := g.CheckArrival()
		g.UpdateAnimation()
		if finished {
			g.state++
		}
	case StateResult:
		done := g.HandleResults()
		if done {
			g.Reset()
			g.state = StateWelcomeScreen
		}
	}
	return nil
}
