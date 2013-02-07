package goatpress

import (
  "net"
  "os"
  "fmt"
  "time"
  "net/http"
  "html/template"
  "sort"
)

const serverAddress = "localhost:4123"

var newPlayers = make(chan Player)
var removePlayers = make(chan string)

type Server struct {
  Tournament *Tournament
  dataPath   string
}
var server *Server

func newServer(dataPath string) *Server {
  gameType := newGameType(5, DefaultWordSet)
  tourney := newTournament(*gameType, dataPath)
  randomPlayer := newInternalPlayer("Random", newRandomFinder(DefaultWordSet))
  tourney.RegisterPlayer(randomPlayer)
  server = &Server{tourney, dataPath}
  return server
}

func (c *Server) Run() {
  go c.RunWeb()
  c.RunTournament()
}

type HomePage struct {
  PlayerCount int
  Players     []PlayerStats
}

type PlayerStats struct {
  Name string
  Score int
  Games int
  Moves int
  Wins int
  Draws int
  Losses int
}

func homePage(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("src/goatpress/views/index.html")
  if err != nil {
    fmt.Println(err)
    w.Write([]byte("error"))
    return
  }
  pc := len(server.Tournament.Players)
  stats := make([]PlayerStats, 0)
  scores := server.Tournament.Scores
  for name, _ := range server.Tournament.AllPlayerNames {
    g := scores.Games[name]
    m := scores.Moves[name]
    w := scores.Wins[name]
    l := scores.Losses[name]
    d := g - w - l
    s := g + m/10 + d + 10*w
    stat := PlayerStats{name, s, g, m ,w, d, l}
    stats = append(stats, stat)
  }
  vs := NewValSorter(stats, func (s PlayerStats) int { return s.Score*-1 })
  vs.Sort()
  
  t.Execute(w, &HomePage{pc, vs.Keys})
}

func (c *Server) RunWeb() {
  http.HandleFunc("/", homePage)
  http.ListenAndServe(":5123", nil)
}

func (c *Server) RunTournament() {
  listener, err := net.Listen("tcp", serverAddress)
  if err != nil {
    fmt.Printf("error listening:", err.Error())
    os.Exit(1)
  }
  matchTicker := 0
  go AcceptPlayers(listener)

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
    default:
      for _, player := range c.Tournament.Players {
        player.Ping()
      }
      if c.Tournament.Size() > 1 {
        c.Tournament.PlayMatch()
        matchTicker++
        if matchTicker > 100 {
          c.Tournament.Save()
        }
      } else {
        time.Sleep(0.2*1e9)
      }
    }
    fmt.Printf("Players: %s\n", c.Tournament.PlayerList())
    time.Sleep(1)
  }
}

const serverSig = "goatpress<VERSION=1> ; \n"

func AcceptPlayers(listener net.Listener) {
  for {
    conn, err := listener.Accept()
    if err != nil {
      println("Error accept:", err.Error())
      return
    }
    conn.Write([]byte(serverSig))
    player := newClientPlayer(conn, removePlayers)
    if player != nil {
      newPlayers <- player
    }
  }
}

type ValSorter struct {
        Keys []PlayerStats
        Vals []int
}
 
func NewValSorter(values []PlayerStats, mapper func(PlayerStats) int) *ValSorter {
        vs := &ValSorter{
                Keys: make([]PlayerStats, 0, len(values)),
                Vals: make([]int, 0, len(values)),
        }
        for _, v := range values {
                vs.Keys = append(vs.Keys, v)
                vs.Vals = append(vs.Vals, mapper(v))
        }
        return vs
}
 
func (vs *ValSorter) Sort() {
        sort.Sort(vs)
}
 
func (vs *ValSorter) Len() int           { return len(vs.Vals) }
func (vs *ValSorter) Less(i, j int) bool { return vs.Vals[i] < vs.Vals[j] }
func (vs *ValSorter) Swap(i, j int) {
        vs.Vals[i], vs.Vals[j] = vs.Vals[j], vs.Vals[i]
        vs.Keys[i], vs.Keys[j] = vs.Keys[j], vs.Keys[i]
}


