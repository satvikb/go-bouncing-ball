// main.go
package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Vector2 struct {
	X float32
	Y float32
}

const (
	screenWidth  = 800
	screenHeight = 600
)

var ballPosition Vector2
var ballVelocity Vector2

const ballRadius = 10

func main() {
	ballPosition = Vector2{X: 400, Y: 300}
	ballVelocity = Vector2{X: 10, Y: 100}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
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
						case sdl.K_p:

						}
					}
				}
			}

			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()

			vertsToDraw := getCirclePoints()

			drawVertices(renderer, vertsToDraw, sdl.Color{255, 255, 255, 255})

			renderer.Present()

			tickBallPosition(float32(elapsed.Seconds()))
			fmt.Println("elapsed time: ", elapsed.Milliseconds())
			// }

		}
	}
}

func tickBallPosition(delta float32) {

	ballPosition.X += ballVelocity.X * delta
	ballPosition.Y += ballVelocity.Y * delta

	if ballPosition.X > screenWidth-ballRadius {
		ballPosition.X = screenWidth - ballRadius
		ballVelocity.X = -ballVelocity.X
	}

	if ballPosition.X < ballRadius {
		ballPosition.X = ballRadius
		ballVelocity.X = -ballVelocity.X

	}

	if ballPosition.Y > screenHeight-ballRadius {
		ballPosition.Y = screenHeight - ballRadius
		ballVelocity.Y = -ballVelocity.Y
	}

	if ballPosition.Y < ballRadius {
		ballPosition.Y = ballRadius
		ballVelocity.Y = -ballVelocity.Y
	}

}

func getCirclePoints() []Vector2 {
	// use the midpoint circle algorithm

	var center Vector2 = ballPosition
	const radius float32 = ballRadius

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
