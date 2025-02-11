package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Graph struct {
	adjacencyList map[int][]int
}

func (g *Graph) addEdge(v, u int) {
	g.adjacencyList[v] = append(g.adjacencyList[v], u)
}

func readGraph() Graph {
	g := Graph{adjacencyList: make(map[int][]int)}
	var n, m int

	fmt.Print("Введіть кількість вершин та ребер: ")
	fmt.Scanf("%d %d\n", &n, &m)

	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < m; i++ {
		fmt.Print("Введіть ребро (v u): ")
		line, _ := reader.ReadString('\n')
		edge := strings.Fields(line)
		v, _ := strconv.Atoi(edge[0])
		u, _ := strconv.Atoi(edge[1])
		g.addEdge(v, u)
	}

	return g
}

func dependencyIndex(g Graph, setM, setN []int) int {
	index := 0
	for _, v := range setM {
		for _, u := range g.adjacencyList[v] {
			for _, n := range setN {
				if u == n {
					index++
				}
			}
		}
	}
	return index
}

func main() {
	graph := readGraph()
	sets := make(map[string][]int)
	var setsCount int
	fmt.Print("Введіть кількість множин: ")
	fmt.Scanf("%d\n", &setsCount)

	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < setsCount; i++ {
		fmt.Print("Введіть назву множини та вершини (назва v1 v2 ...): ")
		line, _ := reader.ReadString('\n')
		parts := strings.Fields(line)
		name := parts[0]
		var vertices []int
		for _, vStr := range parts[1:] {
			v, _ := strconv.Atoi(vStr)
			vertices = append(vertices, v)
		}
		sets[name] = vertices
	}

	weightedGraph := make(map[string]map[string]int)
	for nameM, setM := range sets {
		weightedGraph[nameM] = make(map[string]int)
		for nameN, setN := range sets {
			if nameM != nameN {
				index := dependencyIndex(graph, setM, setN)
				if index > 0 {
					weightedGraph[nameM][nameN] = index
				}
			}
		}
	}

	fmt.Println("\nЗважений граф:")
	for nameM, neighbors := range weightedGraph {
		for nameN, weight := range neighbors {
			fmt.Printf("%s -> %s [вага: %d]\n", nameM, nameN, weight)
		}
	}
}
