package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
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
	Paddle1Y     float32
	Paddle2Y     float32
	BallX        float32
	BallY        float32
	BallDX       float32
	BallDY       float32
	Player1Score int
	Player2Score int
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

	// Ball exit from left
	if g.BallX < 0 {
		g.Player2Score++
		g.resetBallPosition()
	}

	// Ball exit from right
	if g.BallX > width {
		g.Player1Score++
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

	// Draw the scores
	scoreText1 := fmt.Sprintf("Player 1: %d", g.Player1Score)
	scoreText2 := fmt.Sprintf("Player 2: %d", g.Player2Score)
	text.Draw(screen, scoreText1, basicfont.Face7x13, 20, 20, color.White)
	text.Draw(screen, scoreText2, basicfont.Face7x13, 460, 20, color.White)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	game := &Game{
		Paddle1Y:     (height / 2) - (paddleHeight / 2),
		Paddle2Y:     (height / 2) - (paddleHeight / 2),
		BallX:        width / 2,
		BallY:        height / 2,
		BallDX:       3,
		BallDY:       3,
		Player1Score: 0,
		Player2Score: 0,
	}
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Pong")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
