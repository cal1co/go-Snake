package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	width  = 40
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

type Position struct {
	x int
	y int
}

func newPosition(x, y int) Position {
	return Position{x, y}
}

type Snake struct {
	direction Movement
	head      Position
	length    int
	body      []Position
	bodyMap   map[Position]bool
}

type Fruit struct {
	pos   Position
	value int
}

func newFruit() *Fruit {
	pos := randomPosition()
	fruit := &Fruit{pos, 10}
	return fruit
}

func randomPosition() Position {
	posX, posY := rand.Intn(width-2)+1, rand.Intn(height-2)+1
	return Position{posX, posY}
}

type Game struct {
	score    int
	gameOver bool
	snake    *Snake
	fruit    *Fruit
}

func newSnake() *Snake {
	headPos := Position{startX, startY}
	return &Snake{Right, headPos, 1, []Position{headPos}, map[Position]bool{headPos: true}}
}

func (s *Snake) move() {
	butt := s.body[0]
	s.body = s.body[1:]
	delete(s.bodyMap, butt)

	switch s.direction {
	case Up:
		s.head.y--
	case Down:
		s.head.y++
	case Left:
		s.head.x--
	case Right:
		s.head.x++
	}
	s.body = append(s.body, s.head)
	s.bodyMap[s.head] = true
}

func (g *Game) checkEat() {
	if g.snake.head == g.fruit.pos {
		g.fruit.pos = randomPosition()
	}
}

func newSnakeGame() *Game {
	game := &Game{0, false, newSnake(), newFruit()}
	return game
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (g *Game) draw() {
	clearScreen()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				fmt.Print("#")
			} else if x == g.fruit.pos.x && y == g.fruit.pos.y {
				fmt.Print("■")
			} else {
				var isSnakePart bool
				pos := newPosition(x, y)
				if g.snake.bodyMap[pos] {
					isSnakePart = true
					fmt.Print("s")
				}
				if !isSnakePart {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}
	fmt.Printf("Score: %d-%d\n", g.score, g.fruit.pos)
}

func (g *Game) update() {
	g.checkEat()
	g.draw()
	g.snake.move()
}

func main() {
	// handle fruit eat
	// handle growing snake
	// handle collisions
	// handle gameover
	// handle rampup speed

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()
	game := newSnakeGame()

	go userInputListener(game)

	for !game.gameOver {
		game.update()
		time.Sleep(200 * time.Millisecond)
	}
}

func userInputListener(game *Game) {
	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			log.Fatal(err)
		}
		if key == keyboard.KeyEsc {
			game.gameOver = true
			return
		}
		if key == keyboard.KeyArrowUp || char == 'w' && game.snake.direction != Down {
			game.snake.direction = Up
			continue
		}
		if key == keyboard.KeyArrowDown || char == 's' && game.snake.direction != Up {
			game.snake.direction = Down
			continue
		}
		if key == keyboard.KeyArrowLeft || char == 'a' && game.snake.direction != Right {
			game.snake.direction = Left
			continue
		}
		if key == keyboard.KeyArrowRight || char == 'd' && game.snake.direction != Left {
			game.snake.direction = Right
			continue
		}
	}
}
