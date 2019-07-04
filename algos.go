package main

import (
	"fmt"
	"image"
	"math"
	"sort"
)

func (env *Env) botPlayer(algo string) {
	// Select Algo
	//TODO append algos here
	switch {
	case algo == "Astar":
		go env.aStar()
	}
}

type node struct {
	pos       image.Point
	cost      int
	heuristic int
}

func (env *Env) getNeighboorsList(currNode *node) []*node {
	var neighbors []*node
	if env.checkMove(currNode.pos, LEFT) {
		pos := image.Point{currNode.pos.X - 1, currNode.pos.Y}
		cost := currNode.cost + 1
		heuristic := cost + env.getDist(pos, env.end)
		neighbors = append(neighbors, &node{pos, cost, heuristic})
	}
	if env.checkMove(currNode.pos, UP) {
		pos := image.Point{currNode.pos.X, currNode.pos.Y - 1}
		cost := currNode.cost + 1
		heuristic := cost + env.getDist(pos, env.end)
		neighbors = append(neighbors, &node{pos, cost, heuristic})
	}
	if env.checkMove(currNode.pos, RIGHT) {
		pos := image.Point{currNode.pos.X + 1, currNode.pos.Y}
		cost := currNode.cost + 1
		heuristic := cost + env.getDist(pos, env.end)
		neighbors = append(neighbors, &node{pos, cost, heuristic})
	}
	if env.checkMove(currNode.pos, DOWN) {
		pos := image.Point{currNode.pos.X, currNode.pos.Y + 1}
		cost := currNode.cost + 1
		heuristic := cost + env.getDist(pos, env.end)
		neighbors = append(neighbors, &node{pos, cost, heuristic})
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
	openList = append(openList, &node{env.start, 0, 0})
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
			//TODO Run Track
			for _, node := range closedList {
				fmt.Println(node.cost, " , ", node.heuristic)
			}
			return
		}
		// Eval neighbors
		neighbors := env.getNeighboorsList(currNode)
		for _, neighbor := range neighbors {
			// check if neighbor exists in closedList then continue
			res := env.isPresent(neighbor, closedList)
			if res >= 0 {
				continue
			}
			// check if neighbor exists in openList with lower cost, then continue
			res = env.isPresent(neighbor, openList)
			if res >= 0 {
				if openList[res].cost < neighbor.cost {
					continue
				}
			}
			openList = append(openList, neighbor)
		}
		if env.isPresent(currNode, closedList) == -1 {
			closedList = append(closedList, currNode)
		}
	}
	fmt.Println("Astar Ended without solution")
}
