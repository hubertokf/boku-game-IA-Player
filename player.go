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
	"bytes"
	_ "sync"
	"math/rand"
	_ "net/http/pprof"
)

var protocol = "http"
// var host = "172.27.12.173"
var host = "localhost"
var port = "8080"
var start = true

type Coord struct {
	x int
	y int
}

type Node struct {
	Score int
	parent *Node
	children []*Node
	Data [][]int
	isOpponent bool
	Movement Coord
}

func (node *Node) AddTerminal(score int, data [][]int, movement Coord) *Node {
	return node.add(score, data, movement)
}

func (node *Node) Add(data [][]int, movement Coord) *Node {
	return node.add(-1, data, movement)
}

func (node *Node) add(score int, data [][]int, movement Coord) *Node {
	childNode := Node{parent: node, Score: score, Data: data, Movement: movement}

	childNode.isOpponent = !node.isOpponent
	node.children = append(node.children, &childNode)

	return &childNode
}

func (node *Node) find_vertical(sequence string) bool{
	
	for _, column := range node.Data {
		var buffer bytes.Buffer
		for _, cell := range column {
			buffer.WriteString(strconv.Itoa(cell))
			if strings.Contains(buffer.String(), sequence){
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
		var buffer bytes.Buffer
		dec := 1
		for cel := coords.x; cel < i; cel++ {
			if cel < 6{
				buffer.WriteString(strconv.Itoa(node.Data[cel][coords.y]))
			}else{
				buffer.WriteString(strconv.Itoa(node.Data[cel][coords.y-dec]))
				dec++
			}
			if strings.Contains(buffer.String(), sequence){
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
		var buffer bytes.Buffer
		i := k

		if coords.x == 0{
			lim--
			k++
		}

		for cel := coords.x; cel <= lim; cel++ {
			
			if cel < 5{
				buffer.WriteString(strconv.Itoa(node.Data[cel][i]))
				i++
			}else{
				buffer.WriteString(strconv.Itoa(node.Data[cel][i]))
			}
			if strings.Contains(buffer.String(), sequence){
				return true
			}
		}
	}
	return false
}

func (node *Node) is_final_state(player int) bool {
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

func (node *Node) block_oponnent_gap(player int) int {
	counter:=0
	var sequences []string
	if player == 1{
		sequences = []string{"12222","21222","22122","22212","22221"}
	}else{
		sequences = []string{"21111","12111","11211","11121","11112"}
	}
	for _, sequence := range sequences {
		if node.find_vertical(sequence){
			counter += 3000
		}
		if node.find_up_diagonal(sequence){
			counter += 3000
		}
		if node.find_down_diagonal(sequence){
			counter += 3000
		}
	}
	return counter
}

func (node *Node) block_oponnent_state(player int) int {
	var sequence string
	var sequence2 string
	counter:=0
	if player == 1{
		sequence = "2221"
		sequence2 = "122221"
	}else{
		sequence = "1112"
		sequence2 = "211112"
	}
	if node.find_vertical(sequence){
		counter += 1000
	}
	if node.find_up_diagonal(sequence){
		counter += 1000
	}
	if node.find_down_diagonal(sequence){
		counter += 1000
	}
	if node.find_vertical(sequence2){
		counter += 1500
	}
	if node.find_up_diagonal(sequence2){
		counter += 1500
	}
	if node.find_down_diagonal(sequence2){
		counter += 1500
	}
	return counter
}

func (node *Node) two_streak(player int) int {
	var sequence string
	counter:=0
	if player == 1{
		sequence = "11"
	}else{
		sequence = "22"
	}
	if node.find_vertical(sequence){
		counter += 100
	}
	if node.find_up_diagonal(sequence){
		counter += 100
	}
	if node.find_down_diagonal(sequence){
		counter += 100
	}
	return counter
}

func (node *Node) tree_streak(player int) int {
	var sequence string
	counter:=0
	if player == 1{
		sequence = "111"
	}else{
		sequence = "222"
	}
	if node.find_vertical(sequence){
		counter += 300
	}
	if node.find_up_diagonal(sequence){
		counter += 300
	}
	if node.find_down_diagonal(sequence){
		counter += 300
	}
	return counter
}

func (node *Node) four_streak(player int) int {
	var sequence string
	counter:=0
	if player == 1{
		sequence = "1111"
	}else{
		sequence = "2222"
	}
	if node.find_vertical(sequence){
		counter += 400
	}
	if node.find_up_diagonal(sequence){
		counter += 400
	}
	if node.find_down_diagonal(sequence){
		counter += 400
	}
	return counter
}

func (node *Node) five_streak(player int) int {
	var sequence string
	counter:=0
	if player == 1{
		sequence = "11111"
	}else{
		sequence = "22222"
	}
	if node.find_vertical(sequence){
		counter += 10000
	}
	if node.find_up_diagonal(sequence){
		counter += 10000
	}
	if node.find_down_diagonal(sequence){
		counter += 10000
	}
	return counter
}

func (node *Node) atk_gap1(player int) int {
	counter:=0
	var sequences []string
	if player == 1{
		sequences = []string{"101"}
	}else{
		sequences = []string{"202"}
	}
	for _, sequence := range sequences {
		if node.find_vertical(sequence){
			counter += 400
		}
		if node.find_up_diagonal(sequence){
			counter += 400
		}
		if node.find_down_diagonal(sequence){
			counter += 400
		}
	}
	return counter
}

func (node *Node) atk_gap2(player int) int {
	counter:=0
	var sequences []string
	if player == 1{
		sequences = []string{"1011"}
	}else{
		sequences = []string{"2022"}
	}
	for _, sequence := range sequences {
		if node.find_vertical(sequence){
			counter += 450
		}
		if node.find_up_diagonal(sequence){
			counter += 450
		}
		if node.find_down_diagonal(sequence){
			counter += 450
		}
	}
	return counter
}

func (node *Node) atk_gap3(player int) int {
	counter:=0
	var sequences []string
	if player == 1{
		sequences = []string{"10111"}
	}else{
		sequences = []string{"20222"}
	}
	for _, sequence := range sequences {
		if node.find_vertical(sequence){
			counter += 550
		}
		if node.find_up_diagonal(sequence){
			counter += 550
		}
		if node.find_down_diagonal(sequence){
			counter += 550
		}
	}
	return counter
}

func (node *Node) heuristic(player int) int {
	counter:=0
	counter += node.atk_gap1(changePlayer(player))
	counter += node.atk_gap2(changePlayer(player))
	counter += node.atk_gap3(changePlayer(player))
	counter += node.block_oponnent_gap(changePlayer(player))
	counter += node.block_oponnent_state(changePlayer(player))
	counter += node.five_streak(changePlayer(player))
	counter += node.four_streak(changePlayer(player))
	counter += node.tree_streak(changePlayer(player))
	counter += node.two_streak(changePlayer(player))

	node.Score = counter
	return counter
}

func copy_board(board [][]int) [][]int {
	duplicate := make([][]int, len(board))
	for i := range board {
		duplicate[i] = make([]int, len(board[i]))
		copy(duplicate[i], board[i])
	}
	return duplicate
}

// Max returns the larger of x or y.
func max(x, y int) int {
    if x < y {
        return y
    }
    return x
}

// Min returns the smaller of x or y.
func min(x, y int) int {
    if x > y {
        return y
    }
    return x
}

func changePlayer(player int) int{
	if player == 1{
		return 2
	}
	return 1
}

func (node *Node) generateChilds(player int) {
	for k_col, col := range node.Data {
		for k_cell, cell := range col {
			if cell == 0 {
				copy := copy_board(node.Data)
				copy[k_col][k_cell] = player

				node.Add(copy, Coord{k_col,k_cell})
			}
		}
	}
}

// func (node *Node) minimax(depth int, maximizingPlayer bool, player int) {
// 	node.generateChilds(player)
// 	for _, cn := range node.children {
// 		if node.is_final_state(player) || depth == 0{
// 			cn.heuristic(player)
// 		}else{
// 			cn.minimax(depth-1, maximizingPlayer, changePlayer(player))
// 		}

// 		if cn.parent.Score == -1 {
// 			cn.parent.Score = cn.Score
// 		} else if cn.isOpponent && cn.Score > cn.parent.Score {
// 			cn.parent.Score = cn.Score
// 		} else if !cn.isOpponent && cn.Score < cn.parent.Score {
// 			cn.parent.Score = cn.Score
// 		}
// 	}
// }


func (node *Node) minimax(depth int, maximizingPlayer bool, player int) int {
	const MaxUint = ^uint(0) 
	const MaxInt = int(MaxUint >> 1) 
	const MinInt = -MaxInt - 1
	if depth == 0 || node.is_final_state(player){
		return node.heuristic(changePlayer(player))
	}

	node.generateChilds(player)

	if maximizingPlayer{
		bestValue := MinInt

		for _ ,child := range node.children {
			v := child.minimax(depth-1, false, changePlayer(player))
			bestValue = max(bestValue, v)
			child.parent.Score = bestValue
			child = nil
		}
		return bestValue
	}else{
		bestValue := MaxInt
		for _, child := range node.children {
			v := child.minimax(depth-1, true, changePlayer(player))
			bestValue = min(bestValue, v)
			child.parent.Score = bestValue
			child = nil
		}
		return bestValue
	}
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

func make_move(player int, movement Coord) []int {
	var response []int
	s := fmt.Sprintf("move?format=json&player=%d&coluna=%d&linha=%d", player, movement.x+1,movement.y+1)
	res := http_get(s)

	dec := json.NewDecoder(strings.NewReader(res))
	dec.Decode(&response)

	return response
}

func (node *Node) GetBestChildNode() *Node {
	children_size := len(node.children)
	child := node.children[(children_size/2+9)]
	for k := 1; k < children_size; k++ {
		if node.children[k].Score > child.Score {
			child = node.children[k]
		}
	}

	return child
}


func main() {
	player := 1
	plies := 3
	if len(os.Args) > 1 {
		player, _ = strconv.Atoi(os.Args[1])
	}
	if len(os.Args) > 2 {
		plies, _ = strconv.Atoi(os.Args[2])
	}
	// restart_board()
	
	go http.ListenAndServe(":8000", http.DefaultServeMux)
	
	for {
		if player == get_player() {
			start := time.Now()

			movements := get_movements()
			if len(movements) > 2 {

				board := get_board()
				tree := Node{Data:board, isOpponent: true}

				tree.minimax(plies, false, player)

				e := tree.GetBestChildNode()

				res := make_move(player, e.Movement)
				fmt.Println(res)

				if res[0] == -5{
					board := get_board()
					tree := Node{Data:board, isOpponent: true}
					tree.Data[e.Movement.x][e.Movement.y] = 9

					tree.minimax(plies, false, player)

					e := tree.GetBestChildNode()

					res := make_move(player, e.Movement)
					fmt.Println(res)
				}

				if res[0] == -6{
					s1 := rand.NewSource(time.Now().UnixNano())
					r1 := rand.New(s1)
					c := Coord{movements[r1.Intn(2)][0]-1,movements[r1.Intn(2)][1]-1}
					res := make_move(player, c)
					fmt.Println(res)
				}
			} else {
				s1 := rand.NewSource(time.Now().UnixNano())
				r1 := rand.New(s1)
				c := Coord{movements[r1.Intn(2)][0]-1,movements[r1.Intn(2)][1]-1}
				res := make_move(player, c)
				fmt.Println(res)
			}

			elapsed := time.Since(start)
			log.Printf("Time %s", elapsed)

			time.Sleep(2 * time.Second)
		}
	}
	
}
