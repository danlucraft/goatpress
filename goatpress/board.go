package goatpress

import (
	"fmt"
	"math/rand"
)

// *** Tile

type Tile []int

func newTile(i int, j int) Tile {
	return []int{i, j}
}

func (t *Tile) X() int {
	return []int(*t)[0]
}

func (t *Tile) Y() int {
	return []int(*t)[1]
}

func (t *Tile) Key() string {
	return fmt.Sprintf("%d-%d", t.X(), t.Y())
}

// *** BoardGenerator

const boardRandomizationCount int = 100

type BoardGenerator struct {
	Words WordSet
}

func defaultBoardGenerator() *BoardGenerator {
	return newBoardGenerator(DefaultWordSet)
}

func newBoardGenerator(words WordSet) *BoardGenerator {
	return &BoardGenerator{words}
}

func (bg *BoardGenerator) newBoard(size int) *Board {
	letters := make([][]string, size)
	for i := 0; i < size; i++ {
		letters[i] = make([]string, size)
	}
	current := 0
	for current < size*size {
		word := bg.Words.ChooseRandom()
		for _, char := range word {
			if current < size*size {
				letters[current/size][current%size] = string(char)
			}
			current += 1
		}
	}
	for i := 0; i < boardRandomizationCount; i++ {
		// swap row
		row1 := rand.Intn(size)
		row2 := rand.Intn(size)
		for j := 0; j < size; j++ {
			tmp := letters[row1][j]
			letters[row1][j] = letters[row2][j]
			letters[row2][j] = tmp
		}
		// swap column
		col1 := rand.Intn(size)
		col2 := rand.Intn(size)
		for j := 0; j < size; j++ {
			tmp := letters[j][col1]
			letters[j][col1] = letters[j][col2]
			letters[j][col2] = tmp
		}
	}
	return &Board{size, letters}
}

// *** Board: a set of letters arranged in a grid

type Board struct {
	Size    int
	Letters [][]string
}

func newEmptyBoard(size int) Board {
	board := Board{size, make([][]string, size)}
	for i, _ := range board.Letters {
		board.Letters[i] = make([]string, size)
	}
	return board
}

func (board *Board) MoveFromTiles(tiles []Tile) Move {
	word := ""
	isPass := (len(tiles) == 0)
	for _, tile := range tiles {
        if tile.X() > 4 || tile.Y() > 4 || tile.X() < 0 || tile.Y() < 0 {
          return MakePassMove()
        }
		word += board.Letters[tile.X()][tile.Y()]
	}
	return Move{isPass, tiles, word}
}

func (board *Board) WordFromTiles(tiles []Tile) string {
	word := ""
	for _, c := range tiles {
		word += board.Letters[c.X()][c.Y()]
	}
	return word
}

func (board *Board) TilesForLetterExcluding(letter string, tiles []Tile) []Tile {
	hasTiles := make(map[string]bool)
	for _, tile := range tiles {
		hasTiles[tile.Key()] = true
	}
	result := make([]Tile, 0)
	for i := 0; i < board.Size; i++ {
		for j := 0; j < board.Size; j++ {
			if board.Letters[i][j] == letter {
				tile := newTile(i, j)
				if !hasTiles[tile.Key()] {
					result = append(result, tile)
					hasTiles[tile.Key()] = true
				}
			}
		}
	}
	return result
}

func (board *Board) RandomMoveFromWord(word string) Move {
	hasLetters := board.HasLetters()

	// prefilter on has the right letters
	found1 := true
	j := 0
	for found1 && j < len(word) {
		if !hasLetters[string(word[j])] {
			found1 = false
		}
		j++
	}

	// construct the move, if possible
	if found1 {
		moveTiles := make([]Tile, 0)
		for _, char := range word {
			candidateTiles := board.TilesForLetterExcluding(string(char), moveTiles)
			if len(candidateTiles) > 0 {
				tile := candidateTiles[rand.Intn(len(candidateTiles))]
				moveTiles = append(moveTiles, tile)
			} else {
				return MakePassMove()
			}
		}
		return board.MoveFromTiles(moveTiles)
	}
	return MakePassMove()
}

func (board *Board) HasLetters() map[string]bool {
	result := make(map[string]bool)
	for _, row := range board.Letters {
		for _, letter := range row {
			result[letter] = true
		}
	}
	return result
}

func (board *Board) ToString() string {
	r := ""
	l := board.Size
	for i := 0; i < l; i++ {
		for j := 0; j < l; j++ {
			r += board.Letters[i][j]
		}
		if i < l-1 {
			r += " "
		}
	}
	return r
}

type ColorMask [][]int

func isNotPlayers(color int, player int) bool {
	otherPlayer := (((player - 1) + 1) % 2) + 1
	return color == 0 || color == otherPlayer
}

func isDark(b [][]int, player int, x int, y int) bool {
	offsets := [][]int{[]int{-1, 0},
		[]int{0, -1}, []int{0, 0}, []int{0, 1},
		[]int{1, 0}}
	size := len(b)

	for _, offset := range offsets {
		x1 := x + offset[0]
		y1 := y + offset[1]
		if x1 >= 0 && x1 < size && y1 >= 0 && y1 < size {
			if isNotPlayers(b[x1][y1], player) {
				return false
			}
		}
	}
	return true
}

func newColorMask(b *Board, moves []Move) ColorMask {
	currColors := make([][]int, b.Size)
	prevColors := make([][]int, b.Size)
	for i := 0; i < b.Size; i++ {
		currColors[i] = make([]int, b.Size)
		prevColors[i] = make([]int, b.Size)
	}

	curr := 0
	prev := 1
	boards := [][][]int{currColors, prevColors}
	moveCount := 0
	for _, move := range moves {
		player := (moveCount % 2) + 1
		otherPlayer := ((moveCount + 1) % 2) + 1
		for i := 0; i < b.Size; i++ {
			for j := 0; j < b.Size; j++ {
				boards[curr][i][j] = boards[prev][i][j]
			}
		}
		for _, tile := range move.Tiles {
			x := tile.X()
			y := tile.Y()
			if !isDark(boards[prev], otherPlayer, x, y) {
				boards[curr][x][y] = player
			}
		}
		moveCount += 1
		curr = (curr + 1) % 2
		prev = (prev + 1) % 2
	}
	return ColorMask(boards[prev])
}

func (cm *ColorMask) Score(player int) int {
	score := 0
	for _, row := range [][]int(*cm) {
		for _, c := range row {
			if c == player || c == player+2 {
				score++
			}
		}
	}
	return score
}

func (cm *ColorMask) ToString() string {
	r := ""
	l := len([][]int(*cm))
	for i := 0; i < l; i++ {
		for j := 0; j < l; j++ {
			r += fmt.Sprintf("%d", [][]int(*cm)[i][j])
		}
		if i < l-1 {
			r += " "
		}
	}
	return r
}
