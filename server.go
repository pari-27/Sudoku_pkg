package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
)

var Wait = 50000 * time.Second

type Cell struct {
	x int `json:"x"`
	y int `json:"y"`
}

var ans_grid map[Cell]int = make(map[Cell]int)

func InitRouter() (router *mux.Router) {
	router = mux.NewRouter()
	//fmt.Println("hello from init")

	// No version requirement for /ping
	// router.HandleFunc("/new").Methods(http.MethodGet)
	router.PathPrefix("/staticFiles").Handler(http.StripPrefix("/staticFiles", http.FileServer(http.Dir("./staticFiles/"))))

	router.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	router.HandleFunc("/ws", gameHandler).Methods(http.MethodGet)

	return
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	sudokuGrid := map[Cell]int{}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	_, level, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	DiffLevels := map[string]int{"0": 50, "1": 40, "2": 30, "3": 20}
	emptyCells := DiffLevels[string(level)]

	init_grid(sudokuGrid)
	generate_grid(sudokuGrid, emptyCells)
	get_ansGrid(sudokuGrid)

	conn.SetWriteDeadline(time.Now().Add(Wait))

	json_resp := get_JSONuserGrid(sudokuGrid)
	conn.WriteMessage(websocket.TextMessage, []byte(json_resp))
	for {
		_, inputVal, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		row := int(inputVal[0] / 10)
		col := int(inputVal[0] % 10)
		value := int(inputVal[1])

		//log.Println(row, col, value)
		valid_test := check_input(sudokuGrid, row, col, value)
		log.Println(valid_test)

		conn.WriteMessage(websocket.TextMessage, []byte(valid_test))

		if check_win(sudokuGrid) {
			conn.WriteMessage(websocket.TextMessage, []byte("won"))
			break
		}

	}

	//log.Println(inputVal)

}

//var homeTemplate, err = template.ParseFiles("")

func web_game() {
	//fmt.Println("hello from webgame")
	router := InitRouter()
	server := negroni.Classic()
	server.UseHandler(router)

	server.Run(":3000")
}

func main() {

	web_game()
	//

	// get_ansGrid(sudokuGrid)

	// fmt.Println("---------------------------------")
	// render_grid(sudokuGrid)

}
