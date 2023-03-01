package main

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

/*Cette fonction permet a un BOT de se connecter a un serveur*/

//Fonctionne exactement comme CLient.go mais simule des envois de données (game_update virtuel)
func ClientBot() {
	rand.Seed(time.Now().UnixNano())

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}

	defer conn.Close()

	log.Println("BOT connecté")
	var msg string
	var race float64 = 50
	var skin int = rand.Intn(8)
	var startchrono bool = false
	var arriver bool = false
	var chrono time.Time
	var runtime time.Duration
	for {
		buf := make([]byte, 128)
		_, err := conn.Read(buf)
		msg = strings.Split(string(buf), "\x00")[0]
		if err != nil {
		} else {
			if strings.Contains(msg, "PLY :") {
				conn.Write([]byte("OK"))
			}
			if strings.Contains(msg, "NEXTSTEP") {
				conn.Write([]byte("SUITE"))
			}
			if strings.Contains(msg, "PLAYERSELECTION") {
				conn.Write([]byte("SELECT " + strconv.Itoa(skin) + " " + strconv.FormatBool(true)))
			}
			if strings.Contains(msg, "SKIN") {
				conn.Write([]byte("OK"))
			}
			if strings.Contains(msg, "RACE") {
				if !startchrono {
					chrono = time.Now()
					startchrono = true
				}
				race += rand.Float64() * 10
				if !arriver && race > 750 {
					arriver = true
					runtime = time.Since(chrono)
				}
				conn.Write([]byte("POS:" + strconv.FormatFloat(race, 'f', 2, 64)))
			}
			if strings.Contains(msg, "POS") {
				//log.Println(msg)
				conn.Write([]byte("OK"))
			}
			if strings.Contains(msg, "END") {
				break
			}
			if strings.Contains(msg, "TIMEASK") {
				s, ms := GetSeconds(runtime.Milliseconds())
				msg := "TIME :" + strconv.FormatInt(s-4, 10) + "." + strconv.FormatInt(ms-50, 10)
				conn.Write([]byte(msg))
			}
			if strings.Contains(msg, "TIMEINFO") {
				conn.Write([]byte("OK"))
			}
		}
	}

}
