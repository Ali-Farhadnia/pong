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
	// Original game dimensions
	originalWidth  = 640
	originalHeight = 480

	// Define paddle dimensions relative to the game height
	paddleWidthRatio  = 0.015
	paddleHeightRatio = 0.2

	// Define ball size relative to the game height
	ballSizeRatio = 0.02

	paddleSpeedRatio = 0.01

	// Define border thickness relative to the game height
	borderThicknessRatio = 0.02
)

var (
	// Set the color for the paddles, ball, and border
	paddleColor = color.RGBA{255, 255, 255, 255}
	ballColor   = color.RGBA{255, 255, 255, 255}
	borderColor = color.RGBA{255, 255, 255, 255}
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
	screenWidth, screenHeight := ebiten.WindowSize()
	paddleSpeed := float32(screenHeight) * paddleSpeedRatio

	g.BallX += g.BallDX
	g.BallY += g.BallDY

	g.handleCollisions(screenWidth, screenHeight)

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

func (g *Game) resetBallPosition(screenWidth, screenHeight int) {
	g.BallX = float32(screenWidth) / 2
	g.BallY = float32(screenHeight) / 2
}

func (g *Game) handleCollisions(screenWidth, screenHeight int) {
	paddleWidth := float32(screenWidth) * paddleWidthRatio
	paddleHeight := float32(screenHeight) * paddleHeightRatio
	ballSize := float32(screenHeight) * ballSizeRatio
	borderThickness := float32(screenHeight) * borderThicknessRatio

	// Constrain left paddle within the screen
	if g.Paddle1Y < borderThickness {
		g.Paddle1Y = borderThickness
	}
	if g.Paddle1Y > float32(screenHeight)-paddleHeight-borderThickness {
		g.Paddle1Y = float32(screenHeight) - paddleHeight - borderThickness
	}

	// Constrain right paddle within the screen
	if g.Paddle2Y < borderThickness {
		g.Paddle2Y = borderThickness
	}
	if g.Paddle2Y > float32(screenHeight)-paddleHeight-borderThickness {
		g.Paddle2Y = float32(screenHeight) - paddleHeight - borderThickness
	}

	// Ball exit from left
	if g.BallX < borderThickness {
		g.Player2Score++
		g.resetBallPosition(screenWidth, screenHeight)
	}

	// Ball exit from right
	if g.BallX > float32(screenWidth)-borderThickness {
		g.Player1Score++
		g.resetBallPosition(screenWidth, screenHeight)
	}

	// Check collision with top and bottom borders
	if g.BallY <= borderThickness {
		g.BallY = borderThickness // Prevent ball from moving into the border
		g.BallDY *= -1
	}
	if g.BallY >= float32(screenHeight)-borderThickness {
		g.BallY = float32(screenHeight) - borderThickness // Prevent ball from moving into the border
		g.BallDY *= -1
	}

	// Check collision with the left paddle
	if g.BallX-ballSize <= paddleWidth+borderThickness && g.BallY-ballSize >= g.Paddle1Y && g.BallY+ballSize <= g.Paddle1Y+paddleHeight {
		g.BallDX *= -1
		g.BallX = ballSize + paddleWidth + borderThickness // Adjust position to prevent sticking
	}

	// Check collision with the right paddle
	if g.BallX+ballSize >= float32(screenWidth)-paddleWidth-borderThickness && g.BallY-ballSize >= g.Paddle2Y && g.BallY+ballSize <= g.Paddle2Y+paddleHeight {
		g.BallDX *= -1
		g.BallX = float32(screenWidth) - ballSize - paddleWidth - borderThickness // Adjust position to prevent sticking
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := ebiten.WindowSize()
	paddleWidth := float32(screenWidth) * paddleWidthRatio
	paddleHeight := float32(screenHeight) * paddleHeightRatio
	ballSize := float32(screenHeight) * ballSizeRatio
	borderThickness := float32(screenHeight) * borderThicknessRatio

	// Draw the left paddle
	vector.DrawFilledRect(screen, borderThickness, g.Paddle1Y, paddleWidth, paddleHeight, paddleColor, false)

	// Draw the right paddle
	vector.DrawFilledRect(screen, float32(screenWidth)-paddleWidth-borderThickness, g.Paddle2Y, paddleWidth, paddleHeight, paddleColor, false)

	// Draw the ball
	vector.DrawFilledCircle(screen, g.BallX, g.BallY, ballSize, ballColor, true)

	// Draw the borders
	vector.DrawFilledRect(screen, 0, 0, float32(screenWidth), borderThickness, borderColor, false)                                     // Top border
	vector.DrawFilledRect(screen, 0, float32(screenHeight)-borderThickness, float32(screenWidth), borderThickness, borderColor, false) // Bottom border
	vector.DrawFilledRect(screen, 0, 0, borderThickness, float32(screenHeight), borderColor, false)                                    // Left border
	vector.DrawFilledRect(screen, float32(screenWidth)-borderThickness, 0, borderThickness, float32(screenHeight), borderColor, false) // Right border

	// Draw the scores
	scoreText1 := fmt.Sprintf("Player 1: %d", g.Player1Score)
	scoreText2 := fmt.Sprintf("Player 2: %d", g.Player2Score)
	text.Draw(screen, scoreText1, basicfont.Face7x13, 20, 20, color.White)
	text.Draw(screen, scoreText2, basicfont.Face7x13, screenWidth-160, 20, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	game := &Game{
		Paddle1Y:     originalHeight/2 - originalHeight*paddleHeightRatio/2,
		Paddle2Y:     originalHeight/2 - originalHeight*paddleHeightRatio/2,
		BallX:        originalWidth / 2,
		BallY:        originalHeight / 2,
		BallDX:       3,
		BallDY:       3,
		Player1Score: 0,
		Player2Score: 0,
	}
	ebiten.SetWindowSize(originalWidth, originalHeight)
	ebiten.SetWindowTitle("Pong")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
