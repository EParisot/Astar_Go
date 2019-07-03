package main

import (
	"image"
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
	if env.checkMove(currNode.pos.X-1, currNode.pos.Y) {
		pos := image.Point{currNode.pos.X - 1, currNode.pos.Y}
		cost := currNode.cost + 1
		heuristic := currNode.cost + env.getDist(pos, env.end)
		neighbors = append(neighbors, &node{pos, cost, heuristic})
	}
	if env.checkMove(currNode.pos.X, currNode.pos.Y-1) {
		pos := image.Point{currNode.pos.X, currNode.pos.Y - 1}
		cost := currNode.cost + 1
		heuristic := currNode.cost + env.getDist(pos, env.end)
		neighbors = append(neighbors, &node{pos, cost, heuristic})
	}
	if env.checkMove(currNode.pos.X+1, currNode.pos.Y) {
		pos := image.Point{currNode.pos.X + 1, currNode.pos.Y}
		cost := currNode.cost + 1
		heuristic := currNode.cost + env.getDist(pos, env.end)
		neighbors = append(neighbors, &node{pos, cost, heuristic})
	}
	if env.checkMove(currNode.pos.X, currNode.pos.Y+1) {
		pos := image.Point{currNode.pos.X, currNode.pos.Y + 1}
		cost := currNode.cost + 1
		heuristic := currNode.cost + env.getDist(pos, env.end)
		neighbors = append(neighbors, &node{pos, cost, heuristic})
	}
	return neighbors
}

func (env *Env) getDist(x, y image.Point) int {
	var dist int
	// Manhattan Distance Calculation
	return dist
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
		// take first
		currNode := openList[0]
		if env.checkEnd(currNode.pos.X, currNode.pos.Y) {
			//TODO Run Track
			return
		}
		// Eval neighbors
		neighbors := env.getNeighboorsList(currNode)
		for _, neighbor := range neighbors {
			for _, node := range closedList {
				if neighbor.pos.X == node.pos.X && neighbor.pos.Y == node.pos.Y {
					continue
				}
			}
			for _, node := range openList {
				if (neighbor.pos.X == node.pos.X || neighbor.pos.Y == node.pos.Y) &&
					neighbor.cost >= node.cost {
					continue
				}
			}
			openList = append(openList, neighbor)
		}
		closedList = append(closedList, currNode)
	}
}
