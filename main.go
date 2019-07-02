package main

import (
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
		subSq.Fill(env.playerCol)
	}
	// Append sub in full
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(env.offset), float64(env.offset))
	emptySq.DrawImage(subSq, op)

	return emptySq
}

func buildMap(size int, start, end [2]int) Env {
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
		playerCol: color.RGBA{0, 0, 0, 255},
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

	//Draw all squares and populate sqList
	for y := 0; y < size; y++ {
		env.sqList[y] = make([]*square, size)
		for x := 0; x < size; x++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*env.sqW), float64(y*env.sqW))
			if start[0] == x && start[1] == y {
				env.grid.DrawImage(startSq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x * env.sqW, y * env.sqW},
					state: 1,
				}
			} else if end[0] == x && end[1] == y {
				env.grid.DrawImage(endSq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x * env.sqW, y * env.sqW},
					state: 2,
				}
			} else {
				env.grid.DrawImage(emptySq, op)
				env.sqList[y][x] = &square{
					pos:   image.Point{x * env.sqW, y * env.sqW},
					state: 0,
				}
			}
		}
	}
	return env
}

func parseArgs() (int, [2]int, [2]int) {
	//Parse Args
	if len(os.Args[1:]) < 3 {
		fmt.Printf("Missing Argument\nUsage : Astart_go size start end\nwith :\n\tsize : int (map width)\n\tstart : int,int (ex:6,8)\n\tend : int,int (ex:31,30)\n")
		os.Exit(1)
	}
	size, err := strconv.Atoi(os.Args[1])
	if err != nil || size < 8 {
		fmt.Println("Invalid Argument")
		os.Exit(1)
	}
	startStr := os.Args[2]
	resStart := strings.Split(startStr, ",")
	var startTab [2]int
	sr0, err0 := strconv.Atoi(resStart[0])
	if err0 != nil {
		fmt.Println("Invalid Argument")
		os.Exit(1)
	}
	startTab[0] = sr0
	sr1, err1 := strconv.Atoi(resStart[1])
	if err1 != nil {
		fmt.Println("Invalid Argument")
		os.Exit(1)
	}
	startTab[1] = sr1
	endStr := os.Args[3]
	resEnd := strings.Split(endStr, ",")
	var endTab [2]int
	er0, err0 := strconv.Atoi(resEnd[0])
	if err0 != nil {
		fmt.Println("Invalid Argument")
		os.Exit(1)
	}
	endTab[0] = er0
	er1, err1 := strconv.Atoi(resEnd[1])
	if err1 != nil {
		fmt.Println("Invalid Argument")
		os.Exit(1)
	}
	endTab[1] = er1
	return size, startTab, endTab
}

// Update screen 60 time / s
func (env *Env) update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	// print map
	gridOp := &ebiten.DrawImageOptions{}
	gridOp.GeoM.Translate(0, 0)
	screen.DrawImage(env.grid, gridOp)

	// Move
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyUp):
		if env.player.Y-env.sqW >= 0 {
			env.player.Y -= env.sqW
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyDown):
		if env.player.Y+env.sqW < winH {
			env.player.Y += env.sqW
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		if env.player.X+env.sqW < winW {
			env.player.X += env.sqW
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		if env.player.X-env.sqW >= 0 {
			env.player.X -= env.sqW
		}
	}
	// Draw player
	playerSq := env.buildSquare(3)
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(float64(env.player.X), float64(env.player.Y))
	screen.DrawImage(playerSq, playerOp)
	return nil
}

func main() {
	size, startTab, endTab := parseArgs()
	env := buildMap(size, startTab, endTab)
	// Creates main window
	if err := ebiten.Run(env.update, winW, winH, 1, "Astar Go"); err != nil {
		log.Fatal(err)
	}
}
