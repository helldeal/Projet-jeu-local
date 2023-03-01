package main

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

//var race float64 = 50
//liste chrono de chaque joueurs
var timerunners = []float64{
	0,
	0,
	0,
	0,
}

/*Cette fonction permet a un client de se connecter a un serveur*/
func Client() {
	rand.Seed(time.Now().UnixNano())

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Println("CLIENT >>>>>>", "Dial error:", err)
		isConnected = false
		return
	}

	defer conn.Close()

	log.Println("CLIENT >>>>>>", "Je suis connecté")

	var msg string
	for {
		buf := make([]byte, 128)
		_, err := conn.Read(buf)
		msg = strings.Split(string(buf), "\x00")[0]
		if err != nil {
		} else {
			if strings.Contains(msg, "PLY :") {
				//Recoit le nombre joueurs connectés
				isConnected = true
				playerscon = strings.Split(strings.Split(msg, "/")[0], ":")[1]
				psc, _ := strconv.Atoi(playerscon)
				if numPlayers[0] == 0 {
					numPlayers[0] = psc
					for i := range numPlayers {
						if i < psc && i != 0 {
							numPlayers[i] = i
						} else if i != 0 {
							numPlayers[i] = i + 1
						}
					}
				}
				conn.Write([]byte("OK"))
				log.Println("CLIENT >>>>>>", "--> envoyé au serv")

			} else if strings.Contains(msg, "NEXTSTEP") {
				//Recoit l'ordre de changer de scène
				log.Println("CLIENT >>>>>>", "RECEPTION :", msg)
				in <- "next"
				conn.Write([]byte("SUITE"))
				log.Println("CLIENT >>>>>>", "--> envoyé au serv")

			} else if strings.Contains(msg, "PLAYERSELECTION") {
				//Recoit la demande du skin durant la sélection
				log.Println("CLIENT >>>>>>", "RECEPTION :", msg)
				var s = "SELECT " + <-out
				conn.Write([]byte(s))
				log.Println("CLIENT >>>>>>", "--> envoyé au serv")

			} else if strings.Contains(msg, "SKIN") {
				//Recoit les skins des autres joueurs durant la sélection
				skin := strings.Split(strings.Split(msg, "/")[0], ":")[1]
				pl := strings.Split(strings.Split(msg, "/")[1], ":")[1]
				selected := strings.Split(strings.Split(msg, "/")[2], ":")[1]
				log.Println("CLIENT >>>>>>", "Player :", pl, "/ Skin :", skin, "/ Selected :", selected)
				in <- skin + " " + selected
				conn.Write([]byte("OK"))
				log.Println("CLIENT >>>>>>", skinalreaadyselected)

			} else if strings.Contains(msg, "RACE") {
				//Recoit la demande de position durant la course
				log.Println("CLIENT >>>>>>", "RECEPTION :", msg)
				race := <-out
				conn.Write([]byte("POS:" + race))
				//race += rand.Float64() * 7
				//conn.Write([]byte("POS:" + strconv.FormatFloat(race, 'f', 2, 64)))

			} else if strings.Contains(msg, "POS") {
				//Recoit la position des autres durant la course
				log.Println("CLIENT >>>>>>", msg)
				in <- msg
				conn.Write([]byte("OK"))

			} else if strings.Contains(msg, "END") {
				//Recoit l'ordre de se kill
				log.Println("CLIENT >>>>>>", "RECEPTION :", msg)
				log.Println("CLIENT >>>>>>", timerunners)
				log.Println("CLIENT >>>>>>", "GAME FINISHED")
				//race = 0
				//conn.Write([]byte("OK"))
				break

			} else if strings.Contains(msg, "TIMEASK") {
				//Recoit la demande du chrono final
				log.Println("CLIENT >>>>>>", "RECEPTION :", msg)
				s := "TIME :" + <-out
				log.Println("CLIENT >>>>>>", s)
				conn.Write([]byte(s))

			} else if strings.Contains(msg, "TIMEINFO") {
				//Recoit les chronos finaux des autres joueurs
				log.Println("CLIENT >>>>>>", "RECEPTION :", msg)
				times, _ := strconv.ParseFloat(strings.Split(strings.Split(msg, "/")[0], ":")[1], 64)
				timeindex, _ := strconv.Atoi(strings.Split(strings.Split(msg, "/")[1], ":")[1])
				timerunners[timeindex-1] = times
				conn.Write([]byte("OK"))
			}
		}
	}

}
