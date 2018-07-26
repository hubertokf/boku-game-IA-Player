package main

import (
	"fmt"
	"time"
	"log"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
	"strconv"
	_ "net/http/pprof"
)

var protocol = "http"
// var host = "172.27.12.173"
var host = "localhost"
var port = "8080"

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

func main() {
	// player, _ = strconv.Atoi(os.Args[1])
	player := 1
	start := time.Now()

	go http.ListenAndServe(":8000", http.DefaultServeMux)

	for {

		if player == get_player() {
			movements := get_movements()
			if len(movements) > 2 {
				board := get_board()
				tree := NewNode()
				tree.SetNodeData(board)

				tree.Evaluate(3,1)
				// fmt.Print("best")
				// fmt.Println(best.Data)
				// fmt.Println(*tree.Score)
				e := tree.GetBestChildNode()
				// fmt.Println(*e.Score)
				// fmt.Println(e.Data)
				// fmt.Println(e.Movement)
				res := make_move(1, e.Movement)
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
