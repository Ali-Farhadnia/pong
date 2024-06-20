package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Paddle1Y float32
	Paddle2Y float32
	BallX    float32
	BallY    float32
}

func (g *Game) Update() error {
	return nil
}

const (
	// Define layout
	width  = 640
	height = 480

	// Define paddle dimensions
	paddleWidth  = float32(10.0)
	paddleHeight = float32(100.0)

	// Define ball dimensions
	ballSize = float32(10.0)
)

var (
	// Set the color for the paddles and ball
	paddleColor = color.RGBA{255, 255, 255, 255}
	ballColor   = color.RGBA{255, 255, 255, 255}
)

func (g *Game) Draw(screen *ebiten.Image) {

	// Draw the left paddle
	vector.DrawFilledRect(screen, paddleWidth, g.Paddle1Y, paddleWidth, paddleHeight, paddleColor, false)

	// Draw the right paddle
	vector.DrawFilledRect(screen, width-paddleWidth, g.Paddle2Y, paddleWidth, paddleHeight, paddleColor, false)

	// Draw the ball
	vector.DrawFilledCircle(screen, g.BallX, g.BallY, ballSize, ballColor, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	game := &Game{
		Paddle1Y: (height / 2) - (paddleHeight / 2),
		Paddle2Y: (height / 2) - (paddleHeight / 2),
		BallX:    width / 2,
		BallY:    height / 2,
	}
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Pong")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
