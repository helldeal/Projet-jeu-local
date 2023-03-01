package main

import (
	"log"
	"net"
	"strconv"
	"strings"
)

//ON CREE UNE STRUCTURE AFIN DE SIMPLIFIER L'APPEL A CERTAINES VARIABLES ESSENTIELLES AU FONCTIONNEMENT DU SERVER
type server struct {
	connProtocol        string
	connAddr            string
	numberOfConnections int
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//CONNECT VA GERER LA RECEPTION DE CONNEXION DES JOUEURS, L'ENVOI DE RÉPONSES AUX CLIENTS ET LE DÉROULEMENT DU JEU
func connect(srv server, numberOfClients int) {
	//Initialisation du listener
	listener, err := net.Listen(srv.connProtocol, srv.connAddr)
	if err != nil {
		log.Println("SERVER >>>>>> ", "listen error:", err)
		return
	}
	defer listener.Close()

	//Init des chan
	//CHAN ENTRANT --> reponse du client
	var inChans = []chan string{
		make(chan string),
		make(chan string),
		make(chan string),
		make(chan string),
	}
	//CHAN SORTANT --> envoie au client
	var outChans = []chan string{
		make(chan string),
		make(chan string),
		make(chan string),
		make(chan string),
	}

	log.Println("SERVER >>>>>> ", numberOfClients, "clients sont attendus")

	//Acceptation des clients
	for srv.numberOfConnections < numberOfClients {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("SERVER >>>>>> ", "accept error:", err)
			return
		}
		srv.numberOfConnections++

		defer conn.Close()
		log.Println("SERVER >>>>>> ", srv.numberOfConnections, "/", numberOfClients, "clients sont connectés")

		//Créer un goroutine de gestion de client pour chaque client qui se connecte
		go handleClient(conn, srv, inChans[srv.numberOfConnections-1], outChans[srv.numberOfConnections-1], srv.numberOfConnections)

		//Envoie le nombre de personne connectées a chaque joueurs connectés
		for i := 0; i < srv.numberOfConnections; i++ {
			outChans[i] <- "PLY :" + strconv.Itoa(srv.numberOfConnections)
		}
	}

	log.Println("SERVER >>>>>> ", "-------NEXTSTAGE")
	//ON ATTEND QUE LES CLIENTS SE CONNECTENT, PUIS
	//ON PASSE A LA SCENE SUIVANTE, ON ENVOIE LE CODE 2000 -> pour passer à la suite dans le jeu de chaque client
	for i := 0; i < numberOfClients; i++ {
		outChans[i] <- "2000"
	}

	//On attends que tout le monde soit bien à la scène suivante
	for i := 0; i < numberOfClients; i += 0 {
		if <-inChans[i] == "20" {
			i += 1
		}
	}

	log.Println("SERVER >>>>>> ", "-------EVERYBODY IS SELECTING HIS PLAYER SKIN..")
	//liste des skins choisis
	var skins = []string{
		"",
		"",
		"",
		"",
	}
	//liste de boolean pour la selection du skin de chaque joueurs
	var selected = []bool{
		false,
		false,
		false,
		false,
	}
	//Boucle de sélection des skins
	for z := 0; z < 1; z += 0 {
		z = 1
		//ON DEMANDE LES SKINS DONC ON ENVOI LE CODE "3000"
		for i := 0; i < numberOfClients; i++ {
			outChans[i] <- "3000"
		}

		//RECEPTION DES SKINS CHOISIS à l'instant t
		log.Println("SERVER >>>>>> ", "SELECTION :")
		for i := 0; i < numberOfClients; i++ {
			choix := <-inChans[i]
			log.Println("SERVER >>>>>> ", choix)
			skins[i] = strings.Split(choix, " ")[1]
			selected[i], _ = strconv.ParseBool(strings.Split(choix, " ")[2])
			log.Println("SERVER >>>>>> ", "player", i+1, "select", skins[i], selected[i])
		}
		//On envoie à chaque joueurs le skin des autres
		for i := 0; i < numberOfClients; i++ {
			for y := 0; y < numberOfClients; y++ {
				if y != i {
					outChans[i] <- "SKIN :" + skins[y] + "/ PLAYER :" + strconv.Itoa(y+1) + "/ SELECTED :" + strconv.FormatBool(selected[y])
				}
			}
		}
		//Tant que tout les joueurs n'ont pas selectioné un skin on reste dans la boucle
		for i := range selected {
			if !selected[i] {
				z = 0
			}
		}
	}

	//liste des positions de chaque joueurs durant la course
	var pos = []float64{
		0,
		0,
		0,
		0,
	}
	//liste des chronos de chaque joueurs
	var chronos = []float64{
		0,
		0,
		0,
		0,
	}
	//classement
	var ranking = []int{}
	//numéro du joueur vainqueur
	var winner = 0

	//Boucle de la course
	log.Println("SERVER >>>>>> ", "-------LET'S RACE !!!")
	for y := 0; y < 4; y += 0 {
		//ON DEMANDE LES POSITION DONC ON ENVOI LE CODE "4000"
		log.Println("SERVER >>>>>> ", "---------------------------------|----------")
		for i := 0; i < numberOfClients; i++ {
			outChans[i] <- "4000"
		}

		//ON RECOIT LES POSITIONS DE TOUS LES JOUEURS ET GESTION DE L'ARRIVÉE
		for i := 0; i < numberOfClients; i++ {
			pos[i], _ = strconv.ParseFloat(<-inChans[i], 64)
			if contains(ranking, i+1) {
				pos[i] = 750
			}
			if pos[i] > 750 {
				y += 1
				pos[i] = 750
				if winner == 0 {
					winner = i + 1

				}
				ranking = append(ranking, i+1)
			}
			log.Println("SERVER >>>>>> ", i+1, ":", pos[i], "/ 750")
		}
		//ON ENVOIE LA POS DES AUTRES JOUEURS A CHAQUE JOUEURS
		for i := 0; i < numberOfClients; i++ {
			for y := 0; y < numberOfClients; y++ {
				if y != i {
					outChans[i] <- "POS :" + strconv.FormatFloat(pos[y], 'f', 2, 64) + "/ PLAYER :" + strconv.Itoa(y+1)
				}
			}
		}
		//fmt.Println("SERVER >>>>>> ",ranking)
		//fmt.Println("SERVER >>>>>> ","winner : " + strconv.Itoa(winner))

	}

	//LA COURSE FINI DEMANDE LES CHRONOS AVEC CODE "6000"
	for i := 0; i < numberOfClients; i++ {
		outChans[i] <- "6000"
	}

	//ON RECUPERE LES CHRONOS
	for i := 0; i < numberOfClients; i += 0 {
		time := string(<-inChans[i])
		if strings.Contains(time, "TIME") {
			chronos[i], _ = strconv.ParseFloat(strings.Split(time, ":")[1], 64)
			i++
		}
	}

	//ON ENVOIE LES CHRONOS DES AUTRES A CHAQUE JOUEURS
	for i := 0; i < numberOfClients; i++ {
		for y := 0; y < numberOfClients; y++ {
			outChans[i] <- "TIMEINFO :" + strconv.FormatFloat(chronos[y], 'f', 3, 64) + "/ PLAYER :" + strconv.Itoa(y+1)
		}
	}
	log.Println("SERVER >>>>>> ", chronos)

	//LES CLIENTS N'ONT PLUS DE RAISON D'ETRE ON ENVOIE CODE "5000" pour les kill
	for i := 0; i < numberOfClients; i++ {
		outChans[i] <- "5000"
	}
	log.Println("SERVER >>>>>> ", "----------PLAYER", winner, "WIN-----------|----FIN----")

}

//HANDELCLIENT (QUI EST APPELÉ DANS CONNECT SOUS LA FORME D'UNE GOROUTINE) S'OCCUPE DE L'ECHANGE D'INFORMATIONS ENTRE LE CLIENT ET LE SERVEUR.
//POUR COMMUNIQUER AVEC LE CLIENT, HANDLE CLIENT PASSE PAR DES SOCKETS TCP (CONN)
//POUR COMMUNIQUER AVECX LA FONCTION CONNECT, HANDLECLIENT PASSE PAR DES CANAUX (IN POUR LA RECEPTION, ET OUT POUR L'ENVOI).
func handleClient(c net.Conn, srv server, in chan string, out chan string, num int) {
	for {
		out := <-out
		if strings.Contains(out, "PLY") {
			// Envoie au client le nombre de joueurs connectés
			//log.Println("SERVER >>>>>> ",out)
			c.Write([]byte(out))
			readConn(c, in, num)
		} else if out == "2000" {
			//"2000" --> Envoie au client l'info du changement de scène
			c.Write([]byte("NEXTSTEP"))
			readConn(c, in, num)
		} else if out == "3000" {
			//"3000" --> Demande au client son skin choisi pendant la sélection des skins
			c.Write([]byte("PLAYERSELECTION"))
			readConn(c, in, num)
		} else if strings.Contains(out, "SKIN :") {
			// Envoie au client le skin de chaque joueurs durant la selection des skins
			c.Write([]byte(out))
			readConn(c, in, num)
		} else if out == "4000" {
			//"4000" --> Demande au client sa position durant la course
			c.Write([]byte("RACE"))
			readConn(c, in, num)
		} else if strings.Contains(out, "POS :") {
			// Envoie au client la position de chaque joueurs durant la course
			c.Write([]byte(out))
			readConn(c, in, num)
		} else if out == "6000" {
			//"6000" --> Demande au client son chrono final
			c.Write([]byte("TIMEASK"))
			readConn(c, in, num)
		} else if strings.Contains(out, "TIMEINFO :") {
			// Envoie au client les chronos finaux des autres joueurs
			c.Write([]byte(out))
			readConn(c, in, num)
		} else if out == "5000" {
			//"5000" --> Envoie au client l'info de kill le client
			c.Write([]byte("END"))
			break
		}
	}
}

// A CHAQUE ECRITURE AU CLIENT ON ATTEND UNE REPONSE AVEC readConn()
func readConn(c net.Conn, in chan string, num int) {
	buf := make([]byte, 128)
	_, err := c.Read(buf)
	var msg = strings.Split(string(buf), "\x00")[0]
	if err != nil {
	} else {
		//log.Println("SERVER >>>>>> ","RECEPTION DE", num, ":", msg)
		if strings.Contains(msg, "SUITE") {
			//Le client a bien changé de scène   --> Réponse au code "2000"
			in <- "20"
		} else if strings.Contains(msg, "SELECT") {
			//Le client envoie son skin sélectionné durant la sélection   --> Réponse au code "3000"
			in <- msg
		} else if strings.Contains(msg, "OK") {
			//Le client confirme la réception de l'écriture (pas d'information à envoyé du client au serv)
		} else if strings.Contains(msg, "POS") {
			//Le client envoie sa postion durant la course    --> Réponse au code "4000"
			pos := strings.Split(msg, ":")[1]
			in <- pos
		} else if strings.Contains(msg, "TIME") {
			//Le client envoie son chrono final   --> Réponse au code "6000"
			in <- msg
		}
	}
}

func Server() {

	//Création du server
	srv := server{"tcp", "localhost:8080", 0}
	connect(srv, 4)
}
