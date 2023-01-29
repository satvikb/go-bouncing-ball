// main.go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Vector2 struct {
	X float32
	Y float32
}

type Ball struct {
	Position Vector2
	Radius   float32
	Velocity Vector2
	Color    sdl.Color
}

const (
	screenWidth  = 800
	screenHeight = 600
)

var balls []*Ball = make([]*Ball, 0)

const ballThickness = 2

func main() {

	for i := 0; i < 10; i++ {
		balls = append(balls, MakeBall())
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return
	}

	err := ttf.Init()
	if err != nil {
		return
	}

	window, err := sdl.CreateWindow(
		"GO BOUNCING BALL",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return
	}

	font, err := ttf.OpenFont("OpenSans-Regular.ttf", 100)
	if err != nil {
		fmt.Println("ERROR ", err)
		return
	}

	//tick := time.Tick(16 * time.Millisecond)
	gameTickDelta := time.Second / 60

	gameTicker := time.NewTicker(gameTickDelta)

	then := time.Now()

	for {

		select {
		case <-gameTicker.C:
			elapsed := time.Since(then)

			// for {
			then = time.Now()

			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					return

				case *sdl.KeyboardEvent:
					keyEvent := event.(*sdl.KeyboardEvent)
					if keyEvent.Type == sdl.KEYDOWN {
						switch keyEvent.Keysym.Sym {
						case sdl.K_n:
							balls = append(balls, MakeBall())
						case sdl.K_SPACE:
							balls = append(balls, MakeBall())
						case sdl.K_ESCAPE:
							return
						}

					}
				}
			}

			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()

			renderer.SetScale(ballThickness, ballThickness)

			for _, ball := range balls {
				vertsToDraw := getCirclePoints(ball)
				drawVertices(renderer, vertsToDraw, ball.Color)

				tickBallPosition(ball, float32(elapsed.Seconds()))
				fmt.Println("elapsed time: ", elapsed.Milliseconds(), "ms")
				// }
			}
			renderer.SetScale(1, 1)

			text := "Balls: " + fmt.Sprint(len(balls)) + ", Press n or SPACE to add ball"
			surface, _ := font.RenderUTF8Blended(text, sdl.Color{255, 255, 255, 255})
			texture, _ := renderer.CreateTextureFromSurface(surface)
			surface.Free()
			renderer.Copy(texture, nil, &sdl.Rect{0, 0, 500, 50})
			texture.Destroy()

			renderer.Present()

		}
	}
	font.Close()
}

func MakeBall() *Ball {

	// generate random position and velocity and radius
	var randPosition Vector2 = Vector2{X: float32(rand.Intn(screenWidth)), Y: float32(rand.Intn(screenHeight))}
	var randVelocity Vector2 = Vector2{X: randRange(-300, 300), Y: randRange(-300, 300)}
	var randRadius float32 = randRange(10, 50)
	// also randomize color
	var randColor sdl.Color = sdl.Color{R: uint8(randRange(0, 255)), G: uint8(randRange(0, 255)), B: uint8(randRange(0, 255)), A: 255}

	return &Ball{
		Position: randPosition,
		Radius:   randRadius,
		Velocity: randVelocity,
		Color:    randColor,
	}

}

func randRange(min, max int) float32 {
	return float32(rand.Intn(max-min) + min)
}

func tickBallPosition(ball *Ball, delta float32) {
	var ballPosition Vector2 = ball.Position
	var ballVelocity Vector2 = ball.Velocity

	var ballRadius float32 = ball.Radius

	ball.Position.X += ballVelocity.X * delta
	ball.Position.Y += ballVelocity.Y * delta

	// fmt.Println("ball position: ", ballPosition, " ball velocity: ", ballVelocity)

	if ballPosition.X > screenWidth-ballRadius {
		ball.Position.X = screenWidth - ballRadius
		ball.Velocity.X = -ballVelocity.X
	}

	if ballPosition.X < ballRadius {
		ball.Position.X = ballRadius
		ball.Velocity.X = -ballVelocity.X

	}

	if ballPosition.Y > screenHeight-ballRadius {
		ball.Position.Y = screenHeight - ballRadius
		ball.Velocity.Y = -ballVelocity.Y
	}

	if ballPosition.Y < ballRadius {
		ball.Position.Y = ballRadius
		ball.Velocity.Y = -ballVelocity.Y
	}
}

func getCirclePoints(ball *Ball) []Vector2 {

	var center Vector2 = ball.Position
	var radius float32 = ball.Radius

	// use the midpoint circle algorithm
	x := radius
	y := float32(0)
	err := float32(0)

	verts := make([]Vector2, 0)
	for x >= y {
		verts = append(verts, Vector2{X: center.X + x, Y: center.Y + y})
		verts = append(verts, Vector2{X: center.X + y, Y: center.Y + x})
		verts = append(verts, Vector2{X: center.X - y, Y: center.Y + x})
		verts = append(verts, Vector2{X: center.X - x, Y: center.Y + y})
		verts = append(verts, Vector2{X: center.X - x, Y: center.Y - y})
		verts = append(verts, Vector2{X: center.X - y, Y: center.Y - x})
		verts = append(verts, Vector2{X: center.X + y, Y: center.Y - x})
		verts = append(verts, Vector2{X: center.X + x, Y: center.Y - y})

		if err <= 0 {
			y += 1
			err += 2*y + 1
		}

		if err > 0 {
			x -= 1
			err -= 2*x + 1
		}
	}

	// adjust for scale
	// https://stackoverflow.com/questions/21560384/how-to-specify-width-or-point-size-in-sdl-2-0-draw-points-lines-or-rect

	for i := 0; i < len(verts); i++ {
		verts[i].X /= ballThickness
		verts[i].Y /= ballThickness
	}

	return verts
}

func drawVertices(renderer *sdl.Renderer, vertices []Vector2, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)

	for i := 0; i < len(vertices); i++ {
		v := TransformPoint(vertices[i].X, vertices[i].Y)
		renderer.DrawPoint(v.X, v.Y)
		// fmt.Println("drawing point ", vertices[i])
	}
}

const offsetX = 0
const offsetY = 0
const SCALE = 1

func TransformPoint(x float32, y float32) *sdl.Point {
	return &sdl.Point{X: int32(offsetX + x*SCALE), Y: int32(offsetY + (y * SCALE))}
}
