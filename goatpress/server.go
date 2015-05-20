package goatpress

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"sort"
	"time"
)

const (
	serverAddress      = "172.30.1.246"
	scoreForConnection = 1
	scoreForMove       = 10
	scoreForDraw       = 100
	scoreForWin        = 1000

	clientPingInterval      = 50 * time.Millisecond
	newMatchAttemptInterval = 10 * time.Millisecond

	iterationsBeforeSaving = 100
)

var (
	server        *Server
	newPlayers    = make(chan Player)
	removePlayers = make(chan string)
	matchResults  = make(chan MatchResult)
)

type Server struct {
	Tournament    *Tournament
	dataPath      string
	clientTimeout string
	serverPort    int
	webPort       int
}

func newServer(dataPath string, clientTimeout string, serverPort int, webPort int) *Server {
	gameType := newGameType(5, DefaultWordSet)
	tourney := newTournament(*gameType, dataPath)
	randomPlayer := newInternalPlayer("Random", newRandomFinder(DefaultWordSet))
	randomPlayer2 := newInternalPlayer("Random2", newRandomFinder(HalfWordSet))
	println(DefaultWordSet.Length())
	println(HalfWordSet.Length())
	tourney.RegisterPlayer(randomPlayer)
	tourney.RegisterPlayer(randomPlayer2)
	server = &Server{tourney, dataPath, clientTimeout, serverPort, webPort}
	return server
}

func (s *Server) Run() {
	fmt.Printf("Starting server on %d (web %d)\n", s.serverPort, s.webPort)
	go s.RunWeb()
	s.RunTournament()
}

type HomePage struct {
	PlayerCount int
	Players     []*PlayerInfo
	MatchOffs   []Matchoff
}

type PlayerInfo struct {
	Name             string
	Score            int
	Games            int
	Moves            int
	Wins             int
	Draws            int
	Losses           int
	MeanTime         int64
	IsConnectedClass string
}

type PlayerInfoList []*PlayerInfo

func (l PlayerInfoList) Len() int {
	return len([]*PlayerInfo(l))
}

func (l PlayerInfoList) Less(i int, j int) bool {
	a := []*PlayerInfo(l)
	return a[i].Score > a[j].Score // reverse order
}

func (l PlayerInfoList) Swap(i int, j int) {
	a := []*PlayerInfo(l)
	tmp := a[i]
	a[i] = a[j]
	a[j] = tmp
}

type Matchoff struct {
	People string
	Count  int
}

func homePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("goatpress/views/index.html")
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("error"))
		return
	}

	pc := len(server.Tournament.Players)
	stats := PlayerInfoList{}
	scores := server.Tournament.Scores
	for name, _ := range server.Tournament.AllPlayerNames {
		isConnectedClass := "❌"
		if server.Tournament.Players[name] != nil {
			isConnectedClass = "✅"
		}
		games := scores.Games[name]
		moves := scores.Moves[name]
		wins := scores.Wins[name]
		losses := scores.Losses[name]

		draws := games - wins - losses
		score := scoreForConnection*games + scoreForMove*moves + scoreForDraw*draws + scoreForWin*wins
		mt := int64(0)
		if scores.MoveCounts[name] > 0 {
			mt = (scores.Times[name] / int64(scores.MoveCounts[name])) / 1000
		}
		stat := PlayerInfo{name, score, games, moves, wins, draws, losses, mt, isConnectedClass}
		stats = append(stats, &stat)
	}

	matchOffs := make([]Matchoff, 0)
	for matchOff, count := range scores.WinProduct {
		matchOffs = append(matchOffs, Matchoff{matchOff, count})
	}
	sort.Sort(stats)
	t.Execute(w, &HomePage{pc, stats, matchOffs})
}

func (c *Server) RunWeb() {
	http.HandleFunc("/", homePage)
	port := fmt.Sprintf(":%d", c.webPort)
	http.ListenAndServe(port, nil)
}

func (c *Server) RunTournament() {
	address := serverAddress + fmt.Sprintf(":%d", c.serverPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("error listening: %s\n", err.Error())
		os.Exit(1)
	}
	iterationCount := 0
	pingTicker := time.NewTicker(clientPingInterval)
	matchTicker := time.NewTicker(newMatchAttemptInterval)

	go AcceptPlayers(listener, c.clientTimeout)

	for {
		iterationCount += 1
		// Every so often, display the player list and save the tournament
		if iterationCount > iterationsBeforeSaving {
			fmt.Printf("Players: %s\n", c.Tournament.PlayerList())
			c.Tournament.Save()
			iterationCount = 0
		}
		select {

		case <-pingTicker.C:
			// Pings are
			//for name, player := range c.Tournament.Players {
			//go func() {
			//if !c.Tournament.InProgress[name] {
			//player.Ping()
			//}
			//}()
			//}

		case <-matchTicker.C:
			c.Tournament.PlayMatch()

		case newPlayer := <-newPlayers:
			if newPlayer.Name() != "" {
				fmt.Printf("Player online: %s\n", newPlayer.Name())
				c.Tournament.RegisterPlayer(newPlayer)

				// attempt to get the player to a match right away
				c.Tournament.PlayMatch()
			}

		case removePlayerName := <-removePlayers:
			fmt.Printf("Player offline %s\n", removePlayerName)
			if removePlayerName != "" {
				c.Tournament.DeregisterPlayer(removePlayerName)
			}

		case result := <-matchResults:
			c.Tournament.RecordResult(result)

		}
	}
}

const serverSig = "goatpress<VERSION=1> ; \n"

func AcceptPlayers(listener net.Listener, clientTimeout string) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accept:", err.Error())
			return
		} else {
			conn.Write([]byte(serverSig))
			go IdentifyPlayer(conn, clientTimeout)
		}
	}
}

func IdentifyPlayer(conn net.Conn, clientTimeout string) {
	player := newClientPlayer(conn, removePlayers, clientTimeout)
	if player != nil {
		newPlayers <- player
	}
}
