package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	winW = 640
	winH = 640
)

type Env struct {
	sqList [][]*square
	grid   *ebiten.Image

	player image.Point

	sqW    int
	offset int

	start image.Point
	end   image.Point

	bgCol     color.Color
	lnCol     color.Color
	startCol  color.Color
	endCol    color.Color
	wallCol   color.Color
	playerCol color.Color
}

type square struct {
	pos   image.Point
	state int
}

func (env *Env) buildSquare(state int) *ebiten.Image {
	// Creates square
	emptySq, err := ebiten.NewImage(env.sqW, env.sqW, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	emptySq.Fill(env.lnCol)
	// Creates sub square
	subSq, err := ebiten.NewImage(env.sqW-2*env.offset, env.sqW-2*env.offset, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case state == 0:
		subSq.Fill(env.bgCol)
	case state == 1:
		subSq.Fill(env.startCol)
	case state == 2:
		subSq.Fill(env.endCol)
	case state == 3:
		subSq.Fill(env.wallCol)
	case state == 4:
		subSq.Fill(env.playerCol)
	}
	// Append sub in full
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(env.offset), float64(env.offset))
	emptySq.DrawImage(subSq, op)

	return emptySq
}

func buildMap(size int, start, end []int, readenMap [][]int) Env {
	env := Env{
		sqList:    make([][]*square, size),
		player:    image.Point{start[0], start[1]},
		sqW:       int(winW / size),
		offset:    1,
		start:     image.Point{start[0], start[1]},
		end:       image.Point{end[0], end[1]},
		bgCol:     color.RGBA{255, 255, 255, 255},
		lnCol:     color.RGBA{244, 236, 215, 255},
		startCol:  color.RGBA{178, 76, 99, 155},
		endCol:    color.RGBA{50, 232, 117, 155},
		wallCol:   color.RGBA{0, 0, 0, 155},
		playerCol: color.RGBA{0, 0, 0, 200},
	}
	// Main grid
	var err error
	env.grid, err = ebiten.NewImage(winW, winH, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	emptySq := env.buildSquare(0)
	startSq := env.buildSquare(1)
	endSq := env.buildSquare(2)
	wallSq := env.buildSquare(3)

	//Draw all squares and populate sqList
	for y := 0; y < size; y++ {
		env.sqList[y] = make([]*square, size)
		for x := 0; x < size; x++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*env.sqW), float64(y*env.sqW))
			switch {
			case start[0] == x && start[1] == y:
				env.grid.DrawImage(startSq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x * env.sqW, y * env.sqW},
					state: 1,
				}
			case end[0] == x && end[1] == y:
				env.grid.DrawImage(endSq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x * env.sqW, y * env.sqW},
					state: 2,
				}
			case readenMap[y][x] == 0:
				env.grid.DrawImage(emptySq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x * env.sqW, y * env.sqW},
					state: 0,
				}
			case readenMap[y][x] == 3:
				env.grid.DrawImage(wallSq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x * env.sqW, y * env.sqW},
					state: 3,
				}
			}
		}
	}
	return env
}

func parseArgs() (int, []int, []int, [][]int) {
	//Parse Args
	if len(os.Args[1:]) < 1 {
		fmt.Printf("Missing Argument\nUsage : Astart_go map_file\n")
		os.Exit(1)
	}
	// Read map file
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	firstLine, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("Error Reading map file")
	}
	params := strings.Split(string(firstLine), ";")
	// Read size
	size, err := strconv.Atoi(params[0])
	if err != nil || size < 8 {
		fmt.Println("Invalid Argument")
		os.Exit(1)
	}
	// Read map
	startTab := make([]int, 2)
	endTab := make([]int, 2)
	readenMap := make([][]int, size)
	for y := 0; y < size; y++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("Error Reading map file")
		}
		readenMap[y] = make([]int, size)
		for x := 0; x < size; x++ {
			switch {
			case line[x] == '.':
				readenMap[y][x] = 0
			case line[x] == 's':
				readenMap[y][x] = 1
				startTab[0] = x
				startTab[1] = y
			case line[x] == 'e':
				readenMap[y][x] = 2
				endTab[0] = x
				endTab[1] = y
			case line[x] == '#':
				readenMap[y][x] = 3
			}
		}
	}
	return size, startTab, endTab, readenMap
}

func (env *Env) checkMove(x, y int) bool {
	if env.sqList[y][x].state == 3 {
		return false
	}
	return true
}

func (env *Env) checkEnd(x, y int) bool {
	if env.sqList[y][x].state == 2 {
		return true
	}
	return false
}

func (env *Env) movePlayer() {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		if env.player.Y-1 >= 0 && env.checkMove(env.player.X, env.player.Y-1) {
			env.player.Y--
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		if env.player.Y+1 < winH/env.sqW && env.checkMove(env.player.X, env.player.Y+1) {
			env.player.Y++
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		if env.player.X+1 < winW/env.sqW && env.checkMove(env.player.X+1, env.player.Y) {
			env.player.X++
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		if env.player.X-1 >= 0 && env.checkMove(env.player.X-1, env.player.Y) {
			env.player.X--
		}
	}
}

// Update screen 60 time / s
func (env *Env) update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	// Print map
	gridOp := &ebiten.DrawImageOptions{}
	gridOp.GeoM.Translate(0, 0)
	screen.DrawImage(env.grid, gridOp)
	// Move
	env.movePlayer()
	// Print player
	playerSq := env.buildSquare(4)
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(float64(env.player.X*env.sqW), float64(env.player.Y*env.sqW))
	screen.DrawImage(playerSq, playerOp)
	return nil
}

func main() {
	size, startTab, endTab, readenMap := parseArgs()
	env := buildMap(size, startTab, endTab, readenMap)
	// Creates main window
	if err := ebiten.Run(env.update, winW, winH, 1, "Astar Go"); err != nil {
		log.Fatal(err)
	}
}
