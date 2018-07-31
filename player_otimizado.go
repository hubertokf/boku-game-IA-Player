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
	"sync"
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

func (node *Node) AddTerminal(score int, data [][]int, movement Coord) *Node {
	return node.add(&score, data, movement)
}

func (node *Node) Add(data [][]int, movement Coord) *Node {
	return node.add(nil, data, movement)
}

func (node *Node) add(score *int, data [][]int, movement Coord) *Node {
	childNode := Node{parent: node, Score: score, Data: data, Movement: movement}

	childNode.isOpponent = !node.isOpponent
	node.children = append(node.children, &childNode)

	return &childNode
}

func copy_board(board [][]int) [][]int {
	duplicate := make([][]int, len(board))
	for i := range board {
		duplicate[i] = make([]int, len(board[i]))
		copy(duplicate[i], board[i])
	}
	return duplicate
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
	// Posições das diagonais descendentes
	diags := []Coord{{5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}}
	
	// limite do final da matriz
	lim := 10
	k := 0

	// Percorre cada uma das posições das diagonais descendentes - mandinga
	for _, coords := range diags {
		var s string
		i := k

		if coords.x == 0{
			lim--
			k++
		}

		for cel := coords.x; cel <= lim; cel++ {
			
			if cel < 5{
				s += strconv.Itoa(node.Data[cel][i])
				i++
			}else{
				s += strconv.Itoa(node.Data[cel][i])
			}
			if strings.Contains(s, sequence){
				return true
			}
		}
	}
	return false
}

func (node *Node) is_final_State(player int) bool {
	var sequence string
	if player == 1{
		sequence = "11111"
	}else{
		sequence = "22222"
	}
	if node.find_vertical(sequence) || node.find_up_diagonal(sequence) || node.find_down_diagonal(sequence){
		return true
	}
	return false
}

func (node *Node) heuristic(player int) int {
	counter := 0
	for i:=5; i>=4; i-- {
		// fmt.Println(i)
		sequence := strings.Repeat(strconv.Itoa(player), i)
		if node.find_vertical(sequence){
			counter += i*2
		}
		if node.find_up_diagonal(sequence){
			counter += i*2
		}
		if node.find_down_diagonal(sequence){
			counter += i*2
		} 
	}

	node.Score = &counter
	return counter
}

func (node *Node) generateChilds(player int) {
	data_size := len(node.Data)
	for col := 0; col < data_size; col++ {
		col_size := len(node.Data[col])
		for cell := 0; cell < col_size; cell++ {
			if node.Data[col][cell] == 0 {
				copy := copy_board(node.Data)
				copy[col][cell] = player

				node.Add(copy, Coord{col,cell})

			}
		}
	}

	// return node
}

func (node *Node) Evaluate(plies int, player int, remove bool) {

	node.generateChilds(player)

	if player == 1{
		player = 2
	}else{
		player = 1
	}

	children_size := len(node.children)
	var wg sync.WaitGroup
	wg.Add(children_size)

	for k := 0; k < children_size; k++ {
		go func(k int) {
			if plies != 0 {

				node.children[k].Evaluate(plies-1, player, true)
			}else{
				node.children[k].heuristic(player)
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
			defer wg.Done()
		}(k)
	}

	wg.Wait()
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

func restart_board() {
	http_get("reiniciar")
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
	var child *Node
	for k := 0; k < children_size; k++ {
		if *node.children[k].Score > *node.Score {
			child = node.children[k]
		}else{
			child = node.children[children_size/2]
		}
	}

	return child
}

func main() {
	player := 1
	plies :=1
	if len(os.Args) > 1 {
		player, _ = strconv.Atoi(os.Args[1])
	}
	if len(os.Args) > 2 {
		plies, _ = strconv.Atoi(os.Args[2])
	}
	// restart_board()
	
	// go http.ListenAndServe(":8000", http.DefaultServeMux)
	
	for {
		
		if player == get_player() {
			start := time.Now()

			movements := get_movements()
			if len(movements) > 2 {

				board := get_board()
				tree := Node{Data:board, isOpponent: true}

				tree.Evaluate(plies, player, false)
				// fmt.Println(tree)
				e := tree.GetBestChildNode()
				fmt.Println(*e.Score)

				// fmt.Println(*e.Score)
				
				res := make_move(player, e.Movement)
				fmt.Println(res)
			} else {
				fmt.Println("...")
			}

			elapsed := time.Since(start)
			log.Printf("Time %s", elapsed)

			time.Sleep(2 * time.Second)
		}
	}
	
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

