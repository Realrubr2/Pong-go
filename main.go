package main

import (
	"log"

	"github.com/Realrubr2/Pong-go/colours"
	"github.com/hajimehoshi/ebiten/v2"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"

	"fmt"

	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// some constants
const (
	screenWidth  = 640
	screenHeight = 480
	ballSpeed = 4
	paddleSpeed = 6
)

type Object struct {
	X, Y, W, H int
}

type Paddle struct {
	Object
	color      color.Color
}

type Ball struct {
	Object
	dxdt int // x velocity of the ball
	dydt int // y velocity of the ball
}

type Game struct {
	paddle Paddle
	ball Ball
	score int
	highScore int
	currentColor int
	colors []color.Color
}

// main function
func main() {
	ebiten.SetWindowTitle("Pong Ramon")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	paddle := NewPaddle(600, 200, 15, 100, color.White)
	ball := NewBall(0, 0, 15, 15, ballSpeed, ballSpeed)
	g := &Game{
		paddle: paddle,
		ball: ball,
		colors: colours.InitColor(),
	}

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
// basicly returns the screen width and height
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// draws the game state
func (g *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 
		float32(g.paddle.X), float32(g.paddle.Y), 
		float32(g.paddle.W), float32(g.paddle.H), 
		g.paddle.color, false,
	)
	vector.DrawFilledRect(screen, 
		float32(g.ball.X), float32(g.ball.Y), 
		float32(g.ball.W), float32(g.ball.H), 
		color.White, false,
	)

	scoreStr := "Score: " + fmt.Sprint(g.score)
	text.Draw(screen, scoreStr, basicfont.Face7x13, 10, 10, color.White)

	highScoreStr := "High Score: " + fmt.Sprint(g.highScore)
	text.Draw(screen, highScoreStr, basicfont.Face7x13, 10, 30, color.White)

}


// updates the game state
func (g *Game) Update() error {	
	g.paddle.MoveOnKeyPress()
	g.ball.Move()
	g.CollideWithWall()
	g.CollideWithPaddle()
	return nil
}

// moves the paddle on key press
func (p *Paddle) MoveOnKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.Y += paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.Y -= paddleSpeed
	}
}

//moving the ball
func (b *Ball) Move() {
	b.X += b.dxdt
	b.Y += b.dydt
}


// reset the count if lost
func (g *Game) Reset() {
	g.ball.X = 0
	g.ball.Y = 0

	g.score = 0
}

// Right wall causes a game over
func (g *Game) CollideWithWall() {

	if g.ball.X >= screenWidth {
		g.Reset()
	}else if g.ball.X <= 0{
		g.ball.dxdt = ballSpeed
	}else if g.ball.Y <= 0 {
		g.ball.dydt = ballSpeed
	}else if g.ball.Y >= screenHeight {
		g.ball.dydt = -ballSpeed
	}
}

func (g *Game) CollideWithPaddle() {
	if g.ball.X + g.ball.W >= g.paddle.X &&
	g.paddle.X + g.paddle.W >= g.ball.X &&
	g.ball.Y + g.ball.H >= g.paddle.Y &&
	g.paddle.Y + g.paddle.H >= g.ball.Y {
		g.ball.dxdt = -g.ball.dxdt
		g.score++
		if g.score > g.highScore {
			g.highScore = g.score
		}
		g.UpdatePaddleColour()
	}
}

// updates the paddle colour
func(g *Game) UpdatePaddleColour(){
	g.currentColor = (g.currentColor + 1) % len(g.colors)
	g.paddle.color = g.colors[g.currentColor]
}

func NewPaddle(x, y, w, h int, color color.Color) Paddle {
	return Paddle{
		Object: Object{
			X: x,
			Y: y,
			W: w,
			H: h,
		},
		color: color,
	}
}
// here we return a new Ball object
func NewBall(x, y, w, h, dxdt, dydt int) Ball {
	return Ball{
		Object: Object{
			X: x,
			Y: y,
			W: w,
			H: h,
		},
		dxdt: dxdt,
		dydt: dydt,
	}
}