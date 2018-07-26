// Package minimax implements the minimax algorithm
// Minimax (sometimes MinMax or MM[1]) is a decision rule used in decision theory,
// game theory, statistics and philosophy for minimizing the possible loss for
// a worst case (maximum loss) scenario
// See for more details: https://en.wikipedia.org/wiki/Minimax
package main

import (
	"fmt"
	"strconv"
	_ "github.com/mohae/deepcopy"
	"runtime"
)

// Node represents an element in the decision tree
type Node struct {
	// Score is available when supplied by an evaluation function or when calculated
	Score      *int
	parent     *Node
	children   []*Node
	isOpponent bool
	Movement [2]int

	// Data field can be used to store additional information by the consumer of the
	// algorithm
	Data [][]int
}

// New returns a new minimax structure
func NewNode() Node {
	n := Node{isOpponent: false}
	return n
}

// Set the data of the node
func (node *Node) SetNodeData(d [][]int) *Node {
	node.Data = d

	return node
}

func (node *Node) GetNodeData() [][]int {
	return node.Data
}

// GetBestChildNode returns the first child node with the matching score
func (node *Node) GetBestChildNode() *Node {
	for _, cn := range node.children {
		if cn.Score == node.Score {
			return cn
		}
	}

	return nil
}

// Evaluate runs through the tree and caculates the score from the terminal nodes
// all the the way up to the root node
func (node *Node) Evaluate(plies int, player int) {
	eval := false
	if plies == 0{
		eval = true
	}
	node.generateChilds(player, eval)
	// PrintMemUsage()
	for _, cn := range node.children {
		// fmt.Println(cn.Data)
		if plies != 0 {
			if player == 1{
				player = 2
			}else{
				player = 1
			}
			cn.Evaluate(plies-1, player)
		}

		if cn.parent.Score == nil {
			cn.parent.Score = cn.Score
		} else if cn.isOpponent && *cn.Score > *cn.parent.Score {
			cn.parent.Score = cn.Score
		} else if !cn.isOpponent && *cn.Score < *cn.parent.Score {
			cn.parent.Score = cn.Score
		}
	}
}

// Print the node for debugging purposes
func (node *Node) Print(level int) {
	var padding = ""
	for j := 0; j < level; j++ {
		padding += " "
	}

	var s = ""
	if node.Score != nil {
		s = strconv.Itoa(*node.Score)
	}

	fmt.Println(padding, node.isOpponent, node.Data, "["+s+"]")

	for _, cn := range node.children {
		level += 2
		cn.Print(level)
		level -= 2
	}
}

// AddTerminal adds a terminal node (or leave node).  These nodes
// should contain a score and no children
func (node *Node) AddTerminal(score int, data [][]int, movement [2]int) *Node {
	return node.add(&score, data, movement)
}

// Add a new node to structure, this node should have children and
// an unknown score
func (node *Node) Add(data [][]int, movement [2]int) *Node {
	return node.add(nil, data, movement)
}

func (node *Node) add(score *int, data [][]int, movement [2]int) *Node {
	childNode := Node{parent: node, Score: score, Data: data, Movement: movement}

	childNode.isOpponent = !node.isOpponent
	node.children = append(node.children, &childNode)
	return &childNode
}

func (node *Node) isTerminal() bool {
	return len(node.children) == 0
}

func (node *Node) neighbors(pos []int) [][]int {
	column := pos[0]
	line := pos[1]

	var position []int
	var l [][]int

	if line < len(node.Data[column])-1 { // DOWN
		position = []int{column, line + 1}
		l = append(l, position)
	}

	if column <= len(node.Data)-1 && column != 0 { // DIAGONAL L/D
		if line != len(node.Data[column])-1 && column < 6 {
			position = []int{column - 1, line}
			l = append(l, position)
		} else if column >= 6 {
			position = []int{column - 1, line + 1}
			l = append(l, position)
		}
	}

	if column <= len(node.Data)-1 && column != 0 { // DIAGONAL L/U
		if column < 6 && line != 0 {
			position = []int{column - 1, line - 1}
			l = append(l, position)
		} else if column >= 6 {
			position = []int{column - 1, line}
			l = append(l, position)
		}
	}

	if line != 0 {
		position = []int{column, line - 1} // UP
		l = append(l, position)
	}

	if column < len(node.Data)-1 { // DIAGONAL R/U
		if column < 5 {
			position = []int{column + 1, line}
			l = append(l, position)
		} else if column >= 5 && line != 0 {
			position = []int{column + 1, line - 1}
			l = append(l, position)
		}
	}

	if column < len(node.Data)-1 { // DIAGONAL R/D
		if column < 5 {
			position = []int{column + 1, line + 1}
			l = append(l, position)
		} else if column >= 5 && line != len(node.Data[column+1]) {
			position = []int{column + 1, line}
			l = append(l, position)
		}
	}

	return l
}

func (node *Node) heuristic(position []int) int {
	counter := 0
	v := node.neighbors(position)

	for _, vizinho := range v {
		if node.Data[vizinho[0]][vizinho[1]] == 0 {
			counter = counter + 10
		}
		if node.Data[vizinho[0]][vizinho[1]] == 1 { // VERIFICAR o 1
			counter = counter + 20
		}
		if node.Data[vizinho[0]][vizinho[1]] == 2 {
			counter = counter - 50
		}
	}
	// node.Score = counter
	return counter
}

func copy_board(original_board [][]int) [][]int {

	var board [][]int = make([][]int, len(original_board))

	for x, list := range original_board {
		board[x] = make([]int, len(list))
		for y, value := range list {
			board[x][y] = value
			// fmt.Println(x, y, value)
		}
	}
	return board
}

func PrintMemUsage() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
        fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
        fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
        fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}

func (node *Node) generateChilds(player int, eval bool) *Node {
	for k_col, col := range node.Data {
		for k_cell, cell := range col {
			if cell == 0 {
				node.Data[k_col][k_cell] = player
				board := copy_board(node.Data)
				// board := deepcopy.Copy(node.Data)
				node.Data[k_col][k_cell] = 0

				// PrintMemUsage()

				if eval == true {
					// evaluate state
					score := node.heuristic([]int{k_col,k_cell})
					// add child node to node child list
					node.AddTerminal(score, board, [2]int{k_col,k_cell})
				}else{
					// add child node to node child list
					node.Add(board, [2]int{k_col,k_cell})
				}

			}
		}
	}

	return node
}