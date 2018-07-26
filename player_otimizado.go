package main

import (
	"fmt"
	"strconv"
	"time"
	"log"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
	"os"
)


var protocol = "http"
// var host = "172.27.12.173"
var host = "localhost"
var port = "8080"

type Node struct {
	Score *int
	parent *Node
	children []*Node
	Data [][]int
	isOpponent bool
	Movement [2]int
}

func (node *Node) AddTerminal(score int, data [][]int, movement [2]int) {
	node.add(&score, data, movement)
}

func (node *Node) Add(data [][]int, movement [2]int) {
	node.add(nil, data, movement)
}

func (node *Node) add(score *int, data [][]int, movement [2]int) {
	childNode := Node{parent: node, Score: score, Data: data, Movement: movement}

	childNode.isOpponent = !node.isOpponent
	node.children = append(node.children, &childNode)
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

	for k := 0; k < len(v); k++ {
		if node.Data[v[k][0]][v[k][1]] == 0 {
			counter = counter + 10
		}
		if node.Data[v[k][0]][v[k][1]] == 1 { // VERIFICAR o 1
			counter = counter + 20
		}
		if node.Data[v[k][0]][v[k][1]] == 2 {
			counter = counter - 50
		}
	}
	// node.Score = counter
	return counter
}

func (node *Node) Evaluate(plies int, player int, remove bool) {
	eval := false
	if plies == 0{
		eval = true
	}
	node.generateChilds(player,eval)
	children_size := len(node.children)
	for k := 0; k < children_size; k++ {
		if plies != 0 {
			if player == 1{
				player = 2
			}else{
				player = 1
			}
			node.children[k].Evaluate(plies-1, player, true)
		}

		if node.children[k].parent.Score == nil {
			node.children[k].parent.Score = node.children[k].Score
		} else if node.children[k].isOpponent && *node.children[k].Score > *node.children[k].parent.Score {
			node.children[k].parent.Score = node.children[k].Score
		} else if !node.children[k].isOpponent && *node.children[k].Score < *node.children[k].parent.Score {
			node.children[k].parent.Score = node.children[k].Score
		}

		if remove == true{
			node.children[k] = nil
		}
	}
}

func (node *Node) generateChilds(player int, eval bool) {
	data_size := len(node.Data)
	for col := 0; col < data_size; col++ {
		col_size := len(node.Data[col])
		for cell := 0; cell < col_size; cell++ {
			if node.Data[col][cell] == 0 {
				node.Data[col][cell] = player
				
				if eval == true {
					score := node.heuristic([]int{col,cell})
					node.AddTerminal(score, node.Data, [2]int{col,cell})
				}else{
					node.Add(node.Data, [2]int{col,cell})
				}

				node.Data[col][cell] = 0
			}
		}
	}

	// return node
}

func http_get(dir string) string {
	resp, _ := http.Get(protocol + "://" + host + ":" + port + "/" + dir)
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	return string(bodyBytes)
}

func get_player() int {
	res := http_get("jogador")

	playerNumber, _ := strconv.Atoi(res)

	return playerNumber
}

func get_board() [][]int {
	var board [][]int

	res := http_get("tabuleiro")

	dec := json.NewDecoder(strings.NewReader(res))
	dec.Decode(&board)

	return board
}

func get_movements() [][]int {
	var movements [][]int

	res := http_get("movimentos?format=json")

	dec := json.NewDecoder(strings.NewReader(res))
	dec.Decode(&movements)

	return movements
}

func make_move(player int, movement [2]int) string {
	s := fmt.Sprintf("move?player=%d&coluna=%d&linha=%d", player, movement[0]+1,movement[1]+1)
	res := http_get(s)

	return res
}

func (node *Node) GetBestChildNode() *Node {
	children_size := len(node.children)
	for k := 0; k < children_size; k++ {
		if node.children[k].Score == node.Score {
			return node.children[k]
		}
	}

	return nil
}

func main() {
	player := 1
	if len(os.Args) > 1 {
		player, _ = strconv.Atoi(os.Args[1])
	}
	start := time.Now()

	// go http.ListenAndServe(":8000", http.DefaultServeMux)

	for {

		if player == get_player() {
			movements := get_movements()
			if len(movements) > 2 {
				board := get_board()
				tree := Node{Data:board}

				tree.Evaluate(2, player, false)
				e := tree.GetBestChildNode()

				res := make_move(player, e.Movement)
				fmt.Println(res)
			} else {
				fmt.Println("...")
			}
			time.Sleep(2 * time.Second)
		}
	}
	
	elapsed := time.Since(start)
	log.Printf("Time %s", elapsed)
}
