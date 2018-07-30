package main

import (
	"fmt"
	"strconv"
	_ "time"
	_ "log"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
	_ "os"
)


var protocol = "http"
// var host = "172.27.12.173"
var host = "localhost"
var port = "8080"

type Coord struct {
	x int
	y int
}

type Node struct {
	Score *int
	parent *Node
	children []*Node
	Data [][]int
	isOpponent bool
	Movement Coord
}

func (node *Node) AddTerminal(score int, data [][]int, movement Coord) {
	node.add(&score, data, movement)
}

func (node *Node) Add(data [][]int, movement Coord) {
	node.add(nil, data, movement)
}

func (node *Node) add(score *int, data [][]int, movement Coord) {
	childNode := Node{parent: node, Score: score, Data: data, Movement: movement}

	childNode.isOpponent = !node.isOpponent
	node.children = append(node.children, &childNode)
}

func (node *Node) neighbors(pos Coord) []Coord {
	column := pos.x
	line := pos.y

	var l []Coord

	if line < len(node.Data[column])-1 { // DOWN
		position := Coord{column, line + 1}
		l = append(l, position)
	}

	if column <= len(node.Data)-1 && column != 0 { // DIAGONAL L/D
		if line != len(node.Data[column])-1 && column < 6 {
			position := Coord{column - 1, line}
			l = append(l, position)
		} else if column >= 6 {
			position := Coord{column - 1, line + 1}
			l = append(l, position)
		}
	}

	if column <= len(node.Data)-1 && column != 0 { // DIAGONAL L/U
		if column < 6 && line != 0 {
			position := Coord{column - 1, line - 1}
			l = append(l, position)
		} else if column >= 6 {
			position := Coord{column - 1, line}
			l = append(l, position)
		}
	}

	if line != 0 {
		position := Coord{column, line - 1} // UP
		l = append(l, position)
	}

	if column < len(node.Data)-1 { // DIAGONAL R/U
		if column < 5 {
			position := Coord{column + 1, line}
			l = append(l, position)
		} else if column >= 5 && line != 0 {
			position := Coord{column + 1, line - 1}
			l = append(l, position)
		}
	}

	if column < len(node.Data)-1 { // DIAGONAL R/D
		if column < 5 {
			position := Coord{column + 1, line + 1}
			l = append(l, position)
		} else if column >= 5 && line != len(node.Data[column+1]) {
			position := Coord{column + 1, line}
			l = append(l, position)
		}
	}

	return l
}

func (node *Node) find_vertical(sequence string) bool{
	for _, column := range node.Data {
		fmt.Println(column)
		var s string
		for _, cell := range column {
			s += strconv.Itoa(cell)
			if strings.Contains(s, sequence){
				return true
			}
		}
	}
	return false
}

func (node *Node) find_up_diagonal(sequence string) bool{
	diags := []Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 5}, {2, 6}, {3, 7}, {4, 8}, {5, 9}}

	i := 6
	for _, coords := range diags {
		var s string
		dec := 1
		for cel := coords.x; cel < i; cel++ {
			if cel < 6{
				s += strconv.Itoa(node.Data[cel][coords.y])
			}else{
				s += strconv.Itoa(node.Data[cel][coords.y-dec])
				dec++
			}
			if strings.Contains(s, sequence){
				return true
			}

		}
		if i <= 10{
			i++
		}
	}
	return false
}

func (node *Node) find_down_diagonal(sequence string) bool{
	diags := []Coord{{5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}}
	for _, coords := range diags {
		var s string
		for &coords != nil{
			column := coords.x
			line := coords.y
			state := node.Data[column - 1][line - 1]
			s += string(state)

			if strings.Contains(s, sequence){
				return true
			}
			
			coords = node.neighbors(coords)[4]
		}
	}
	return false
}

func (node *Node) final_State() bool {
	return false
}

func (node *Node) heuristic(position Coord) int {
	counter := 0
	v := node.neighbors(position)

	for k := 0; k < len(v); k++ {
		if node.Data[v[k].x][v[k].y] == 0 {
			counter = counter + 10
		}
		if node.Data[v[k].x][v[k].y] == 1 { // VERIFICAR o 1
			counter = counter + 20
		}
		if node.Data[v[k].x][v[k].y] == 2 {
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
		// fmt.Println(node.children[k])
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

func copy_board(board [][]int) [][]int {
	duplicate := make([][]int, len(board))
	for i := range board {
		duplicate[i] = make([]int, len(board[i]))
		copy(duplicate[i], board[i])
	}
	
	// n := len(board)
	// m := len(board[4])
	// duplicate := make([][]int, n)
	// data := make([]int, n*m)
	// for i := range board {
	// 	start := i*m
	// 	end := start + m
	// 	duplicate[i] = data[start:end:end]
	// 	copy(duplicate[i], board[i])
	// }
	return duplicate
}

func (node *Node) generateChilds(player int, eval bool) {
	data_size := len(node.Data)
	for col := 0; col < data_size; col++ {
		col_size := len(node.Data[col])
		for cell := 0; cell < col_size; cell++ {
			if node.Data[col][cell] == 0 {
				copy := copy_board(node.Data)
				copy[col][cell] = player
				
				if eval == true {
					score := node.heuristic(Coord{col,cell})
					node.AddTerminal(score, copy, Coord{col,cell})
				}else{
					node.Add(copy, Coord{col,cell})
				}

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

func make_move(player int, movement Coord) string {
	s := fmt.Sprintf("move?player=%d&coluna=%d&linha=%d", player, movement.x+1,movement.y+1)
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
	// player := 1
	// if len(os.Args) > 1 {
	// 	player, _ = strconv.Atoi(os.Args[1])
	// }
	
	// go http.ListenAndServe(":8000", http.DefaultServeMux)
	
	board := get_board()
	tree := Node{Data:board}

	fmt.Println("vertical: ", tree.find_vertical("1111"))
	fmt.Println("up diagonal: ", tree.find_up_diagonal("1111"))
	fmt.Println("down diagonal: ", tree.find_down_diagonal("1111"))


	// for {
		
	// 	if player == get_player() {
	// 		start := time.Now()

	// 		movements := get_movements()
	// 		if len(movements) > 2 {

	// 			board := get_board()
	// 			tree := Node{Data:board}

	// 			tree.Evaluate(3, player, false)
	// 			e := tree.GetBestChildNode()

	// 			fmt.Println(e.Data)
	// 			fmt.Println(e.Movement)
	// 			fmt.Println(e.neighbors(e.Movement))
				
	// 			tree.heuristic(e.Movement)
				
	// 			res := make_move(player, e.Movement)
	// 			fmt.Println(res)
	// 		} else {
	// 			fmt.Println("...")
	// 		}

	// 		elapsed := time.Since(start)
	// 		log.Printf("Time %s", elapsed)

	// 		time.Sleep(2 * time.Second)
	// 	}
	// }
	
}



// [0 0 0 0 0]
// [0 0 0 0 0 0]
// [0 0 0 0 0 0 0]
// [0 0 0 0 0 0 0 0]
// [0 0 0 0 0 0 0 0 0]
// [0 0 0 0 0 0 0 0 0 0]
// [0 0 0 0 0 0 0 0 0]
// [0 0 0 0 0 0 0 0]
// [0 0 0 0 0 0 0]
// [0 0 0 0 0 0]
// [0 0 0 0 0]

//      [0 0 0 0 0]
//     [0 0 0 0 0 0]
//    [0 0 0 0 0 0 0]
//   [0 0 0 0 0 0 0 0]
//  [0 0 0 0 0 0 0 0 0]
// [0 0 0 0 0 0 0 0 0 0]
//  [0 0 0 0 0 0 0 0 0]
//   [0 0 0 0 0 0 0 0]
//    [0 0 0 0 0 0 0]
//     [0 0 0 0 0 0]
//      [0 0 0 0 0]