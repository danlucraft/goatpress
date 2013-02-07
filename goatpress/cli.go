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

func ServerStart(dataPath string) {
  server := newServer(dataPath)
  server.Run()
}

func ClientStart(name string) {
  client := newClient(name)
  client.Run()
}
