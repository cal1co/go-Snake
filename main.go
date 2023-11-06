package main

import (
	"fmt"
	"time"
)

const (
	width  = 20
	height = 20
	startX = 5
	startY = 5
)

type Movement string

const (
	Up    Movement = "up"
	Down  Movement = "down"
	Left  Movement = "left"
	Right Movement = "right"
)

type Snake struct {
	direction Movement
	x         int
	y         int
	length    int
	body      []string
}
type Fruit struct {
	x     int
	y     int
	value int
}

func newFruit() *Fruit {
	fruit := &Fruit{1, 1, 10}
	return fruit
}

type Game struct {
	score    int
	gameOver bool
	snake    *Snake
	fruit    *Fruit
}

func newSnake() *Snake {
	return &Snake{Left, startX, startY, 1, []string{fmt.Sprintf("%d-%d", startX, startY)}}
}

func newSnakeGame() *Game {
	game := &Game{0, false, newSnake(), newFruit()}
	return game
}
func (g *Game) update() {

}

func main() {
	game := newSnakeGame()

	go func() {
		for {
			game.update()
			time.Sleep(200 * time.Millisecond)
		}
	}()
	newFruit()
}
