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

type Movement string

const (
	width           = 40
	height          = 20
	startX          = 5
	startY          = 5
	Up     Movement = "up"
	Down   Movement = "down"
	Left   Movement = "left"
	Right  Movement = "right"
)

type Position struct {
	x int
	y int
}

func newPosition(x, y int) Position {
	return Position{x, y}
}

func randomPosition() Position {
	posX, posY := rand.Intn(width-2)+1, rand.Intn(height-2)+1
	return Position{posX, posY}
}

type Snake struct {
	direction Movement
	head      Position
	length    int
	body      []Position
	bodyMap   map[Position]bool
	justAte   bool
}

func newSnake() *Snake {
	headPos := Position{startX, startY}
	return &Snake{Right, headPos, 1, []Position{headPos}, map[Position]bool{headPos: true}, false}
}

func (s *Snake) move() {
	butt := s.body[0]
	if !s.justAte {
		s.body = s.body[1:]
		delete(s.bodyMap, butt)
	}
	s.justAte = false

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

type Fruit struct {
	pos   Position
	value int
}

func newFruit() *Fruit {
	pos := randomPosition()
	fruit := &Fruit{pos, 10}
	return fruit
}

type Game struct {
	score    int
	gameOver bool
	snake    *Snake
	fruit    *Fruit
	pace     time.Duration
}

func newSnakeGame() *Game {
	game := &Game{0, false, newSnake(), newFruit(), 200}
	return game
}

func (g *Game) checkEat() {
	if g.snake.head == g.fruit.pos {
		g.eat()
	}
}

func (g *Game) eat() {
	g.fruit.pos = randomPosition()
	g.score += 10
	g.snake.justAte = true
	if g.pace > 75 {
		g.pace -= 10
	}
}

func (g *Game) checkCollision() {
	if g.snake.head.x <= 0 || g.snake.head.x >= width || g.snake.head.y <= 0 || g.snake.head.y >= height {
		g.gameOver = true
	}
	g.checkSelfCollision()
}

func (g *Game) checkSelfCollision() {
	if !g.snake.justAte {
		for i := len(g.snake.body) - 2; i >= 0; i-- {
			if g.snake.head == g.snake.body[i] {
				g.gameOver = true
				return
			}
		}
	}
}

func (g *Game) updateState() {
	g.checkCollision()
	g.snake.move()
	g.checkEat()
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
	fmt.Printf("Score: %d\n", g.score)
}

func (g *Game) update() {
	g.updateState()
	g.draw()

}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
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

func main() {
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()
	game := newSnakeGame()

	go userInputListener(game)

	for !game.gameOver {
		game.update()
		time.Sleep(game.pace * time.Millisecond)
	}
}
