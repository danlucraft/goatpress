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

const serverAddress = "10.32.4.142"

var newPlayers = make(chan Player)
var removePlayers = make(chan string)
var matchResults = make(chan MatchResult)

type Server struct {
	Tournament    *Tournament
	dataPath      string
	clientTimeout string
	serverPort    int
	webPort       int
}

var server *Server

func newServer(dataPath string, clientTimeout string, serverPort int, webPort int) *Server {
	gameType := newGameType(5, DefaultWordSet)
	tourney := newTournament(*gameType, dataPath)
	randomPlayer := newInternalPlayer("Random", newRandomFinder(DefaultWordSet))
	randomPlayer2 := newInternalPlayer("Random2", newRandomFinder(DefaultWordSet))
	tourney.RegisterPlayer(randomPlayer)
	tourney.RegisterPlayer(randomPlayer2)
	server = &Server{tourney, dataPath, clientTimeout, serverPort, webPort}
	return server
}

func (c *Server) Run() {
	fmt.Printf("Starting server on %d (web %d)\n", c.serverPort, c.webPort)
	go c.RunWeb()
	c.RunTournament()
}

type HomePage struct {
	PlayerCount int
	Players     []*PlayerStats
	MatchOffs   []Matchoff
}

type PlayerStats struct {
	Name   string
	Score  int
	Games  int
	Moves  int
	Wins   int
	Draws  int
	Losses int
}

type SortableStatsList []*PlayerStats

func (l *SortableStatsList) Len() int {
	return len([]*PlayerStats(*l))
}

func (l *SortableStatsList) Less(i int, j int) bool {
	a := []*PlayerStats(*l)
	return a[i].Score > a[j].Score // reverse order
}

func (l *SortableStatsList) Swap(i int, j int) {
	a := []*PlayerStats(*l)
	tmp := a[i]
	a[i] = a[j]
	a[j] = tmp
}

type Matchoff struct {
	People string
	Count  int
}

func homePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/goatpress/views/index.html")
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("error"))
		return
	}

	pc := len(server.Tournament.Players)
	stats := SortableStatsList{}
	scores := server.Tournament.Scores
	for name, _ := range server.Tournament.AllPlayerNames {
		g := scores.Games[name]
		m := scores.Moves[name]
		w := scores.Wins[name]
		l := scores.Losses[name]
		d := g - w - l
		s := 10*g + m + 10*d + 100*w
		stat := PlayerStats{name, s / 10, g, m, w, d, l}
		stats = append(stats, &stat)
	}

	matchOffs := make([]Matchoff, 0)
	for matchOff, count := range scores.WinProduct {
		matchOffs = append(matchOffs, Matchoff{matchOff, count})
	}
	sort.Sort(&stats)
	t.Execute(w, &HomePage{pc, stats, matchOffs})
}

func (c *Server) RunWeb() {
	http.HandleFunc("/", homePage)
	port := fmt.Sprintf(":%d", c.webPort)
	println(port)
	http.ListenAndServe(port, nil)
}

func (c *Server) RunTournament() {
	address := serverAddress + fmt.Sprintf(":%d", c.serverPort)
	println(address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("error listening:", err.Error())
		os.Exit(1)
	}
	matchTicker := 0
	playersTicker := 0
	go AcceptPlayers(listener, c.clientTimeout)

	for {
		select {
		case newPlayer := <-newPlayers:
			if newPlayer.Name() != "" {
				fmt.Printf("Player Online: %s\n", newPlayer.Name())
				c.Tournament.RegisterPlayer(newPlayer)
			}
		case removePlayerName := <-removePlayers:
			if removePlayerName != "" {
				c.Tournament.DeregisterPlayer(removePlayerName)
			}
		case result := <-matchResults:
			c.Tournament.RecordResult(result)
		default:
			for name, player := range c.Tournament.Players {
				if !c.Tournament.InProgress[name] {
					player.Ping()
				}
			}
			if c.Tournament.NonPlayingCount() > 1 {
				c.Tournament.PlayMatch()
				matchTicker++
				if matchTicker > 10 {
					c.Tournament.Save()
					matchTicker = 0
				}
			} else {
				time.Sleep(0.2 * 1e9)
			}
		}
		playersTicker++
		if playersTicker > 100000 {
			fmt.Printf("Players: %s\n", c.Tournament.PlayerList())
			playersTicker = 0
		}
		time.Sleep(100)
	}
}

const serverSig = "goatpress<VERSION=1> ; \n"

func AcceptPlayers(listener net.Listener, clientTimeout string) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accept:", err.Error())
			return
		}
		conn.Write([]byte(serverSig))
		player := newClientPlayer(conn, removePlayers, clientTimeout)
		if player != nil {
			newPlayers <- player
		}
	}
}
