package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const debug = false

var heights [][]int
var dimenW, dimenH int

func main() {
	mapFile, err := os.Open("data/map.txt")
	reader := bufio.NewReader(mapFile)

	if err != nil {
		logrus.Error(err)
		return
	}

	line, err := reader.ReadString('\n')
	fmt.Sscanf(line, "%d %d", &dimenW, &dimenH)
	heights = make([][]int, dimenH)

	for l := 0; l < dimenH; l++ {
		line, err := reader.ReadString('\n')
		parseLine(l, line)

		if err == io.EOF {
			break
		}
	}

	now := time.Now()
	bestPath := 0
	bestSteepness := 0
	for y, row := range heights {
		for x := range row {
			path, steepness := BFS(&Node{First: x, Second: y})

			if bestPath < path {
				bestPath = path
				bestSteepness = steepness
				logrus.Infof("new best path=%d steepness=%d", bestPath, bestSteepness)
			} else if bestPath == path && bestSteepness < steepness {
				bestPath = path
				bestSteepness = steepness
				logrus.Infof("new best path=%d steepness=%d", bestPath, bestSteepness)
			}
		}
	}

	logrus.Infof("best path=%d steepness=%d", bestPath, bestSteepness)
	logrus.Infof("time elapsed %s", time.Since(now))
}

type Node struct {
	First  int
	Second int
	Path   int
}

func BFS(root *Node) (path, deltaHeight int) {
	queue := []*Node{root}
	var pair *Node

	height := heights[root.Second][root.First]
	deltaHeight = -1
	root.Path = 0
	path = 0
	bestPath := 0

	for {
		if len(queue) == 0 {
			return
		}

		pair, queue = queue[0], queue[1:]
		path = pair.Path + 1

		if debug {
			logrus.Infof("pair=%#v", pair)
			logrus.Infof("len(q)=%d", len(queue))
			PrintQueue(queue)
		}

		if pair.First-1 >= 0 && heights[pair.Second][pair.First] > heights[pair.Second][pair.First-1] {
			curHeight := heights[pair.Second][pair.First-1]
			queue = append(queue, &Node{First: pair.First - 1, Second: pair.Second, Path: path})

			if bestPath < path {
				bestPath = path
				deltaHeight = height - curHeight
			} else if bestPath == path {
				if height-curHeight > deltaHeight {
					deltaHeight = height - curHeight
				}
			}
		}

		if pair.First+1 <= dimenW-1 && heights[pair.Second][pair.First] > heights[pair.Second][pair.First+1] {
			curHeight := heights[pair.Second][pair.First+1]
			queue = append(queue, &Node{First: pair.First + 1, Second: pair.Second, Path: path})

			if bestPath < path {
				bestPath = path
				deltaHeight = height - curHeight
			} else if bestPath == path {
				if height-curHeight > deltaHeight {
					deltaHeight = height - curHeight
				}
			}
		}

		if pair.Second-1 >= 0 && heights[pair.Second][pair.First] > heights[pair.Second-1][pair.First] {
			curHeight := heights[pair.Second-1][pair.First]
			queue = append(queue, &Node{First: pair.First, Second: pair.Second - 1, Path: path})

			if bestPath < path {
				bestPath = path
				deltaHeight = height - curHeight
			} else if bestPath == path {
				if height-curHeight > deltaHeight {
					deltaHeight = height - curHeight
				}
			}
		}

		if pair.Second+1 <= dimenH-1 && heights[pair.Second][pair.First] > heights[pair.Second+1][pair.First] {
			curHeight := heights[pair.Second+1][pair.First]
			queue = append(queue, &Node{First: pair.First, Second: pair.Second + 1, Path: path})

			if bestPath < path {
				bestPath = path
				deltaHeight = height - curHeight
			} else if bestPath == path {
				if height-curHeight > deltaHeight {
					deltaHeight = height - curHeight
				}
			}
		}

		if debug {
			PrintQueue(queue)
			logrus.Infof("len(q)=%d\n--------------", len(queue))
		}
	}
}

func PrintQueue(queue []*Node) {
	if len(queue) == 0 {
		logrus.Infof("[]")
		return
	}

	for _, p := range queue {
		logrus.Infof("p=%#v", p)
	}
}

func parseLine(l int, line string) {
	heights[l] = make([]int, dimenW)

	inputs := strings.Split(line, " ")
	for i := 0; i < dimenW; i++ {
		height, _ := strconv.Atoi(strings.TrimSpace(inputs[i]))
		heights[l][i] = height
	}
}
