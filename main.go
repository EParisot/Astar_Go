package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	winW = 640
	winH = 640
)

//Env Game's environment
type Env struct {
	sqList [][]*square
	grid   *ebiten.Image
	size   int

	autoMode string
	player   image.Point
	score    int
	over     bool

	sqW    int
	offset int

	start image.Point
	end   image.Point
}

type square struct {
	pos   image.Point
	state int
}

func parseArgs() (string, int, []int, []int, [][]int) {
	//Parse Args
	if len(os.Args[1:]) < 1 {
		fmt.Printf("Missing Argument\nUsage : Astar_go map_file [-m algo]\nwith algo : 'Astar', WIP")
		os.Exit(1)
	}
	// Get args
	mode := ""
	for i, arg := range os.Args {
		if arg == "-m" && i+1 < len(os.Args) {
			mode = os.Args[i+1]
		}
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
	// Read size
	size, err := strconv.Atoi(string(firstLine))
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
	return mode, size, startTab, endTab, readenMap
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
	if env.over == false {
		// Print player
		playerSq := env.buildSquare(color.RGBA{255, 153, 0, 200})
		playerOp := &ebiten.DrawImageOptions{}
		playerOp.GeoM.Translate(float64(env.player.X*env.sqW), float64(env.player.Y*env.sqW))
		screen.DrawImage(playerSq, playerOp)
		// Check if finished
		if env.checkEnd(env.player.X, env.player.Y) {
			env.over = true
		}
	} else {
		// Print Score
		scoreStr := strconv.Itoa(env.score)
		scoreMsg := fmt.Sprintf("GAME OVER\nScore : %s", scoreStr)
		ebitenutil.DebugPrint(screen, scoreMsg)
	}
	return nil
}

func main() {
	mode, size, start, end, readenMap := parseArgs()
	env := Env{
		sqList:   make([][]*square, size),
		size:     size,
		autoMode: mode,
		player:   image.Point{start[0], start[1]},
		score:    0,
		over:     false,
		sqW:      int(winW / size),
		offset:   1,
		start:    image.Point{start[0], start[1]},
		end:      image.Point{end[0], end[1]},
	}
	env.buildMap(size, readenMap)
	// Creates main window
	if len(env.autoMode) > 0 {
		go env.botPlayer(env.autoMode)
	} else {
		go env.movePlayer()
	}
	if err := ebiten.Run(env.update, winW, winH, 1, "Astar Go"); err != nil {
		log.Fatal(err)
	}
}
