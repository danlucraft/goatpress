package goatpress

func Demo() {
  gameType := newGameType(5, DefaultWordSet)
  tourney := newTournament(*gameType,"asdfasdfasdf")
  player1 := newInternalPlayer("Alice", newRandomFinder(DefaultWordSet))
  player2 := newInternalPlayer("Bob", newRandomFinder(DefaultWordSet))

  tourney.RegisterPlayer(player1)
  tourney.RegisterPlayer(player2)

  tourney.PlayMatch()
}

func ServerStart(dataPath string, clientTimeout string, serverPort int, webPort int) {
  server := newServer(dataPath, clientTimeout, serverPort, webPort)
  server.Run()
}

func ClientStart(name string, serverPort int) {
  client := newClient(name, serverPort)
  client.Run()
}
