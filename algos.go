package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func (env *Env) botPlayer(algo string) {
	// wait for grphics
	for {
		if env.started == true {
			break
		}
		time.Sleep(time.Second)
	}
	// Select Algo
	//TODO append algos here
	if algo == "Astar" {
		env.aStar()
	} else {
		fmt.Println("Error Unknown Algorith")
		os.Exit(1)
	}
}

type node struct {
	pos       image.Point
	cost      int
	heuristic int
}

func (env *Env) appendNeighbor(neighbors []*node, currNode *node, pos image.Point) []*node {
	cost := currNode.cost + 1
	heuristic := cost + env.getDist(pos, env.end)
	neighbors = append(neighbors, &node{pos, cost, heuristic})
	return neighbors
}

func (env *Env) getNeighboorsList(currNode *node) []*node {
	var neighbors []*node
	if env.checkMove(currNode.pos, LEFT) {
		pos := image.Point{currNode.pos.X - 1, currNode.pos.Y}
		neighbors = env.appendNeighbor(neighbors, currNode, pos)
	}
	if env.checkMove(currNode.pos, UP) {
		pos := image.Point{currNode.pos.X, currNode.pos.Y - 1}
		neighbors = env.appendNeighbor(neighbors, currNode, pos)
	}
	if env.checkMove(currNode.pos, RIGHT) {
		pos := image.Point{currNode.pos.X + 1, currNode.pos.Y}
		neighbors = env.appendNeighbor(neighbors, currNode, pos)
	}
	if env.checkMove(currNode.pos, DOWN) {
		pos := image.Point{currNode.pos.X, currNode.pos.Y + 1}
		neighbors = env.appendNeighbor(neighbors, currNode, pos)
	}
	return neighbors
}

func (env *Env) getDist(src, dst image.Point) int {
	// Manhattan Distance Calculation
	dist := math.Abs(float64(dst.X-src.X)) + math.Abs(float64(dst.Y-src.Y))
	return int(dist)
}

func (env *Env) isPresent(elem *node, list []*node) int {
	for i, node := range list {
		if elem.pos.X == node.pos.X && elem.pos.Y == node.pos.Y {
			return i
		}
	}
	return -1
}

func (env *Env) aStar() {
	var closedList []*node
	var openList []*node
	// Append start node
	openList = append(openList, &node{env.start, 0, env.getDist(env.start, env.end)})
	for len(openList) != 0 {
		// Sort slice
		sort.Slice(openList, func(i, j int) bool {
			return openList[i].heuristic < openList[j].heuristic
		})
		// Unstack first
		currNode := openList[0]
		openList[0] = nil
		openList = openList[1:]
		if env.checkEnd(currNode.pos.X, currNode.pos.Y) {
			closedList = append(closedList, currNode)
			env.drawMap(closedList)
			env.moveBot(closedList)
			return
		}
		// Eval neighbors
		neighbors := env.getNeighboorsList(currNode)
		validNeighbors := false
		for _, neighbor := range neighbors {
			// check if neighbor exists in closedList then continue
			res := env.isPresent(neighbor, closedList)
			if res >= 0 {
				continue
			}
			// check if neighbor exists in openList with lower cost, then continue
			res = env.isPresent(neighbor, openList)
			if res >= 0 {
				if openList[res].cost <= neighbor.cost {
					continue
				}
			}
			openList = append(openList, neighbor)
			validNeighbors = true
		}
		if validNeighbors && env.isPresent(currNode, closedList) == -1 {
			closedList = append(closedList, currNode)
			env.drawMap(closedList)
		}
	}
	ebitenutil.DebugPrint(env.grid, "GAME OVER\nNo Solution")
	return
}

func (env *Env) drawMap(closedList []*node) {
	costMax := 0
	heurMax := 0
	for _, node := range closedList {
		if node.cost > costMax {
			costMax = node.cost
		}
		if node.heuristic > heurMax {
			heurMax = node.cost
		}
	}
	for _, node := range closedList {
		// Color squares from nodes
		costR := uint8(100 * float64(node.cost) / float64(costMax))
		heurR := uint8(100 * float64(node.heuristic-node.cost) / float64(heurMax))
		sqCol := color.RGBA{heurR, costR, 0, 100}
		sq := env.buildSquare(sqCol)
		sqOp := &ebiten.DrawImageOptions{}
		sqOp.GeoM.Translate(float64(node.pos.X*env.sqW), float64(node.pos.Y*env.sqW))
		env.grid.DrawImage(sq, sqOp)
	}
}

func (env *Env) moveBot(closedList []*node) {
	var finalPath []*node
	var currNode *node
	for i := len(closedList) - 1; i >= 0; i-- {
		if len(finalPath) == 0 {
			finalPath = append(finalPath, closedList[i])
			currNode = closedList[i]
		}
		if closedList[i].cost == currNode.cost-1 &&
			env.checkNextStep(closedList[i], currNode) {
			finalPath = append(finalPath, closedList[i])
			currNode = closedList[i]
		}
	}
	for i := len(finalPath) - 1; i >= 0; i-- {
		// print trail
		if env.player != env.start {
			sqCol := color.RGBA{255, 153, 0, 100}
			sq := env.buildSquare(sqCol)
			sqOp := &ebiten.DrawImageOptions{}
			sqOp.GeoM.Translate(float64(env.player.X*env.sqW),
				float64(env.player.Y*env.sqW))
			env.grid.DrawImage(sq, sqOp)
		}
		// exec step
		env.player.X = finalPath[i].pos.X
		env.player.Y = finalPath[i].pos.Y
		env.score++
		time.Sleep(DELAY)
	}
}

func (env *Env) checkNextStep(node *node, currNode *node) bool {
	if math.Abs(float64(currNode.pos.X-node.pos.X)) == 1 &&
		math.Abs(float64(currNode.pos.Y-node.pos.Y)) == 0 {
		return true
	} else if math.Abs(float64(currNode.pos.X-node.pos.X)) == 0 &&
		math.Abs(float64(currNode.pos.Y-node.pos.Y)) == 1 {
		return true
	} else {
		return false
	}
}
