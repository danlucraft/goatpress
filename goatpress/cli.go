package goatpress

func Demo() {
  gameType := newGameType(5, DefaultWordSet)
  tourney := newTournament(*gameType)
  player1 := newInternalPlayer("Alice", newRandomFinder(DefaultWordSet))
  player2 := newInternalPlayer("Bob", newRandomFinder(DefaultWordSet))

  tourney.RegisterPlayer(player1)
  tourney.RegisterPlayer(player2)

  tourney.PlayMatch()
}

func ServerStart() {
  server := newServer()
  server.Run()
}

func ClientStart() {
  client := newClient()
  client.Run()
}
