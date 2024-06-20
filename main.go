package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	// Define layout
	width  = 640
	height = 480

	// Define paddle dimensions
	paddleWidth  = float32(10.0)
	paddleHeight = float32(100.0)

	// Define ball dimensions
	ballSize = float32(10.0)

	paddleSpeed = 5.0
)

var (
	// Set the color for the paddles and ball
	paddleColor = color.RGBA{255, 255, 255, 255}
	ballColor   = color.RGBA{255, 255, 255, 255}
)

type Game struct {
	Paddle1Y float32
	Paddle2Y float32
	BallX    float32
	BallY    float32
	BallDX   float32
	BallDY   float32
}

func (g *Game) Update() error {
	g.BallX += g.BallDX
	g.BallY += g.BallDY

	g.handleCollisions()

	// Handle left paddle movement (W key)
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Paddle1Y -= paddleSpeed
	}

	// Handle left paddle movement (S key)
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Paddle1Y += paddleSpeed
	}

	// Handle right paddle movement (up arrow key)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.Paddle2Y -= paddleSpeed
	}

	// Handle right paddle movement (down arrow key)
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.Paddle2Y += paddleSpeed
	}

	return nil
}

func (g *Game) resetBallPosition() {
	g.BallX = width / 2
	g.BallY = height / 2

}

func (g *Game) handleCollisions() {
	// Constrain left paddle within the screen
	if g.Paddle1Y < 0 {
		g.Paddle1Y = 0
	}
	if g.Paddle1Y > height-paddleHeight {
		g.Paddle1Y = height - paddleHeight
	}

	// Constrain right paddle within the screen
	if g.Paddle2Y < 0 {
		g.Paddle2Y = 0
	}
	if g.Paddle2Y > height-paddleHeight {
		g.Paddle2Y = height - paddleHeight
	}

	// Ball Out of Window
	if g.BallX < 0 || g.BallX > width {
		g.resetBallPosition()
	}

	// Check collision with top and bottom
	if g.BallY <= 0 || g.BallY >= height {
		g.BallDY *= -1
	}

	// Check collision with the left paddle
	if g.BallX-ballSize <= paddleWidth && g.BallY-ballSize >= g.Paddle1Y && g.BallY+ballSize <= g.Paddle1Y+paddleHeight {
		g.BallDX *= -1                   // Reverse horizontal direction
		g.BallX = ballSize + paddleWidth // Adjust position to prevent sticking
	}

	// Check collision with the right paddle
	if g.BallX+ballSize >= width-paddleWidth && g.BallY-ballSize >= g.Paddle2Y && g.BallY+ballSize <= g.Paddle2Y+paddleHeight {
		g.BallDX *= -1                           // Reverse horizontal direction
		g.BallX = width - ballSize - paddleWidth // Adjust position to prevent sticking
	}

}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the left paddle
	vector.DrawFilledRect(screen, 0, g.Paddle1Y, paddleWidth, paddleHeight, paddleColor, false)

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
		BallDX:   3,
		BallDY:   3,
	}
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Pong")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
