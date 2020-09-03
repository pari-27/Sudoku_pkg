package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

//for the key of the grid
var gridSize int = 9

var chances int = 0

//var result chan<- bool = make(chan bool)

type Grid struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func init_grid(sudokuGrid map[Cell]int) {

	// 4 smaller grid of 2*2
	// a map for each cell location as key and the value of the cell
	chances = 3
	//initial grid

	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			sudokuGrid[Cell{i, j}] = 0
			//fmt.Println(i, j)
		}
	}
	rand.Seed(123456789123456789)

}

func generate_grid(sudokuGrid map[Cell]int, level int) bool {

	count := 0
	for count < level {
		row := rand.Intn(gridSize)
		col := rand.Intn(gridSize)
		number := rand.Intn(gridSize) + 1
		//fmt.Println(row, col, number)
		if safe_grid(sudokuGrid, row, col, number) {
			sudokuGrid[Cell{row, col}] = number
			//render_grid(sudokuGrid)
			if fitGrid(sudokuGrid) {

				//fmt.Println(count)
				count = count + 1
			} else {
				sudokuGrid[Cell{row, col}] = 0
			}
		}
	}
	return true
}

func getUnassignedLocation(sudokuGrid map[Cell]int) Cell {
	//gridSize := int(math.Sqrt(float64(len(sudokuGrid))))
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			if sudokuGrid[Cell{i, j}] == 0 {
				return Cell{i, j}
			}
		}

	}
	return Cell{-1, -1}

}
func check_input(sudokuGrid map[Cell]int, row int, col int, value int) string {

	if ans_grid[Cell{row, col}] == value {
		sudokuGrid[Cell{row, col}] = value
		return "valid"
	} else {
		chances = chances - 1
		if chances == 0 {
			return "loss"
		}
		log.Println(chances)
		return "invalid"
	}

}
func check_win(sudokuGrid map[Cell]int) bool {
	loc := getUnassignedLocation(sudokuGrid)
	if loc.x == -1 && loc.y == -1 {
		return true
	}

	return false
}
func get_testgrid(sudokuGrid map[Cell]int) map[Cell]int {
	test_grid := make(map[Cell]int)
	for i, v := range sudokuGrid {
		test_grid[i] = v

	}
	return test_grid
}
func get_ansGrid(sudukoGrid map[Cell]int) {

	if fitGrid(sudukoGrid) {
		render_grid(ans_grid)
	}
}

func get_JSONuserGrid(sudukoGrid map[Cell]int) string {
	//var str string = "'"
	var json_data []Grid
	for i, v := range sudukoGrid {
		cell := strconv.Itoa(i.x) + strconv.Itoa(i.y)
		// str = str + `{"cell":` + cell + `,"value":` + strconv.Itoa(v) + `}`
		elem := Grid{Key: cell, Value: v}
		json_data = append(json_data, elem)

	}

	// fmt.Println(json_data)
	json_UserGrid, err := json.Marshal(json_data)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(json_UserGrid)

}
func fitGrid(sudokuGrid map[Cell]int) bool {

	test_grid := get_testgrid(sudokuGrid)

	unassignedLocation := getUnassignedLocation(test_grid)
	row, col := unassignedLocation.x, unassignedLocation.y
	if row == -1 && col == -1 {
		ans_grid = get_testgrid(test_grid)
		return true
	}
	var check bool
	// gridSize := int(math.Sqrt(float64(len(sudokuGrid))))

	for k := 1; k <= gridSize; k++ {

		if safe_grid(test_grid, row, col, k) {
			test_grid[Cell{row, col}] = k
			check = fitGrid(test_grid)
			if check == true {
				return true
			}

		}
		test_grid[Cell{row, col}] = 0

	}

	return false

}

func safe_grid(sudokuGrid map[Cell]int, p int, q int, counter int) bool {

	//for the full grid size
	//gridSize := int(math.Sqrt(float64(len(sudokuGrid))))

	empty_flag, row_flag, column_flag, box_flag := true, true, true, true
	//check if the cell is blank
	if (sudokuGrid[Cell{p, q}] != 0) {
		//set flag empty_flag to false
		empty_flag = false
	}

	//check for the same integer in the row and column
	for i := 0; i < gridSize; i++ {
		if i != q && sudokuGrid[Cell{p, i}] == counter {
			row_flag = false
		}
		if i != p && sudokuGrid[Cell{i, q}] == counter {
			column_flag = false
		}

	}
	block := int(math.Sqrt(float64(gridSize)))
	row_index := int(p / block)
	column_index := int(q / block)
	start_index := Cell{row_index * block, column_index * block}
	end_index := Cell{(row_index * block) + (block - 1), (column_index * block) + (block - 1)}

	//check for the same integer the same box
	for i := start_index.x; i <= end_index.x; i++ {
		for j := start_index.y; j <= end_index.y; j++ {
			if sudokuGrid[Cell{i, j}] == counter {
				box_flag = false
			}

		}
	}

	//check for safe grid
	if empty_flag && row_flag && column_flag && box_flag {
		return true
	} else {
		return false
	}

}

func render_grid(sudokuGrid map[Cell]int) {
	//gridSize := int(math.Sqrt(float64(len(sudokuGrid))))

	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {

			str := strings.Join([]string{"[", strconv.Itoa(sudokuGrid[Cell{i, j}]), "]"}, " ")
			if (j+1)%3 == 0 && j+1 < gridSize {
				str = strings.Join([]string{"[", strconv.Itoa(sudokuGrid[Cell{i, j}]), "]", "  |  "}, " ")
			}
			fmt.Printf("%s", str)
		}
		fmt.Println()
	}

}
