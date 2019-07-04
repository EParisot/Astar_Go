package main

import (
	"image"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
)

// Moves opcodes
const (
	LEFT  = 0
	UP    = 1
	RIGHT = 2
	DOWN  = 3

	DELAY = time.Second / 2
)

func (env *Env) buildMap(size int, readenMap [][]int) {
	// Buikd main grid
	var err error
	env.grid, err = ebiten.NewImage(winW, winH, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	// Build basic squares
	emptySq := env.buildSquare(color.RGBA{255, 255, 255, 155})
	startSq := env.buildSquare(color.RGBA{153, 0, 0, 155})
	endSq := env.buildSquare(color.RGBA{0, 102, 0, 155})
	wallSq := env.buildSquare(color.RGBA{0, 0, 0, 155})
	// Draw all squares on grid and populate sqList
	for y := 0; y < size; y++ {
		env.sqList[y] = make([]*square, size)
		for x := 0; x < size; x++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*env.sqW), float64(y*env.sqW))
			switch {
			case readenMap[y][x] == 0:
				env.grid.DrawImage(emptySq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x, y},
					state: 0,
				}
			case readenMap[y][x] == 1:
				env.grid.DrawImage(startSq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x, y},
					state: 1,
				}
			case readenMap[y][x] == 2:
				env.grid.DrawImage(endSq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x, y},
					state: 2,
				}
			case readenMap[y][x] == 3:
				env.grid.DrawImage(wallSq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x, y},
					state: 3,
				}
			}
		}
	}
}

func (env *Env) buildSquare(sqColor color.Color) *ebiten.Image {
	// Creates square
	emptySq, err := ebiten.NewImage(env.sqW, env.sqW, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	emptySq.Fill(color.RGBA{244, 236, 215, 255})
	// Creates sub square
	subSq, err := ebiten.NewImage(env.sqW-2*env.offset, env.sqW-2*env.offset, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	subSq.Fill(sqColor) // bg
	subSq.Fill(sqColor) // start
	subSq.Fill(sqColor) // end
	subSq.Fill(sqColor) // wall
	subSq.Fill(sqColor) // player
	// Append sub square in full
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(env.offset), float64(env.offset))
	emptySq.DrawImage(subSq, op)
	return emptySq
}

func (env *Env) movePlayer() {
	for {
		switch {
		case ebiten.IsKeyPressed(ebiten.KeyLeft):
			env.execStep(LEFT, false)
		case ebiten.IsKeyPressed(ebiten.KeyUp):
			env.execStep(UP, false)
		case ebiten.IsKeyPressed(ebiten.KeyRight):
			env.execStep(RIGHT, false)
		case ebiten.IsKeyPressed(ebiten.KeyDown):
			env.execStep(DOWN, false)
		}
		time.Sleep(DELAY / 4)
	}
}

func (env *Env) checkMove(node image.Point, dir int) bool {
	x := node.X
	y := node.Y
	switch {
	case dir == 0:
		x--
	case dir == 1:
		y--
	case dir == 2:
		x++
	case dir == 3:
		y++
	}
	if x < 0 ||
		y < 0 ||
		x >= winW/env.sqW ||
		y >= winH/env.sqW ||
		env.sqList[y][x].state == 3 {
		return false
	}
	return true
}

func (env *Env) execStep(move int, delay bool) {
	if delay {
		time.Sleep(DELAY)
	}
	moved := 0
	switch {
	case move == LEFT && env.checkMove(env.player, LEFT):
		env.player.X--
		moved = 1
	case move == UP && env.checkMove(env.player, UP):
		env.player.Y--
		moved = 1
	case move == RIGHT && env.checkMove(env.player, RIGHT):
		env.player.X++
		moved = 1
	case move == DOWN && env.checkMove(env.player, DOWN):
		env.player.Y++
		moved = 1
	}
	env.score += moved
	if env.checkEnd(env.player.X, env.player.Y) {
		env.over = true
	}
}

func (env *Env) checkEnd(x, y int) bool {
	if env.end.X == x && env.end.Y == y {
		return true
	}
	return false
}
