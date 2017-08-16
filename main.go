package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"io"

	"errors"

	"github.com/sirupsen/logrus"
)

const debug = false

var points [][]int
var width, height int
var sources []*Vertex

func main() {
	mapFile, err := os.Open("data/map.txt")
	reader := bufio.NewReader(mapFile)

	if err != nil {
		logrus.Error(err)
	}

	line, err := reader.ReadString('\n')
	fmt.Sscanf(line, "%d %d", &width, &height)
	points = make([][]int, height)

	vertChan := make(chan Vertex, width*height)

	for l := 0; l < height; l++ {

		line, err := reader.ReadString('\n')
		go parseLine(vertChan, l, line)

		if err == io.EOF {
			break
		}
	}

	graph := &Graph{}
	graph.Vertices = map[string]Vertex{}

	for i := 0; i < width*height; i++ {
		vertex := <-vertChan
		debugF("storing vertex at %s", vertex.XYStr())
		graph.Vertices[vertex.XYStr()] = vertex
	}

	//Calculate edge
	for y, row := range points {
		for x, _ := range row {
			source, err := graph.Vertex(x, y)

			if err == nil {
				getEdges(graph, source)
			} else {
				logrus.Errorf("error=%v", err)
			}
		}
	}

	debugF("g.V=%d g.E=%d g.EL=%#v", len(graph.Vertices), len(graph.Edges), graph.Edges)
	logrus.Infof("sources len=%d", len(sources))
}

func parseLine(channel chan Vertex, l int, line string) {
	points[l] = make([]int, width)

	numbers := strings.Split(line, " ")
	for i, numberStr := range numbers {
		number, _ := strconv.Atoi(strings.TrimSpace(numberStr))
		points[l][i] = number

		channel <- Vertex{
			X:      i,
			Y:      l,
			Height: number,
		}
	}
}

func debugF(format string, args ...interface{}) {
	if debug {
		if len(args) > 0 {
			logrus.Infof(format, args...)
		} else {
			logrus.Info(format)
		}
	}
}

func getEdges(g *Graph, source *Vertex) {
	hasEdge := false
	if source.X-1 >= 0 {
		dest, err := g.Vertex(source.X-1, source.Y)
		if err != nil {
			logrus.Errorf("error getting vertex, error=%v", err)
		}

		if dest != nil && dest.Height < source.Height {
			g.Edges = append(g.Edges, Edge{
				U:      *source,
				V:      *dest,
				Weight: -1,
			})

			hasEdge = true
		}

	}

	if source.X+1 <= width-1 {
		dest, err := g.Vertex(source.X+1, source.Y)
		if err != nil {
			logrus.Errorf("error getting vertex, error=%v", err)
		}

		if dest != nil && dest.Height < source.Height {
			g.Edges = append(g.Edges, Edge{
				U:      *source,
				V:      *dest,
				Weight: -1,
			})

			hasEdge = true
		}
	}

	if source.Y-1 >= 0 {
		dest, err := g.Vertex(source.X, source.Y-1)
		if err != nil {
			logrus.Errorf("error getting vertex, error=%v", err)
		}

		if dest != nil && dest.Height < source.Height {
			g.Edges = append(g.Edges, Edge{
				U:      *source,
				V:      *dest,
				Weight: -1,
			})

			hasEdge = true
		}
	}

	if source.Y+1 <= height-1 {
		dest, err := g.Vertex(source.X, source.Y+1)
		if err != nil {
			logrus.Errorf("error getting vertex, error=%v", err)
		}

		if dest != nil && dest.Height < source.Height {
			g.Edges = append(g.Edges, Edge{
				U:      *source,
				V:      *dest,
				Weight: -1,
			})

			hasEdge = true
		}
	}

	if hasEdge {
		sources = append(sources, source)
	}
}

type Vertex struct {
	X      int
	Y      int
	Height int
}

func (v *Vertex) XYStr() string {
	return fmt.Sprintf("%d_%d", v.X, v.Y)
}

type Edge struct {
	U      Vertex
	V      Vertex
	Weight int
}

type Graph struct {
	Vertices map[string]Vertex
	Edges    []Edge
}

func (g *Graph) Vertex(x int, y int) (*Vertex, error) {
	debugF("fetching vertex at %d_%d", x, y)
	v, ok := g.Vertices[fmt.Sprintf("%d_%d", x, y)]

	if ok {
		return &v, nil
	}

	return nil, errors.New(fmt.Sprintf("Vertex not found for coord(%d, %d)", x, y))
}