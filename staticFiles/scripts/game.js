var con;
var chances=3;

function generateGrid() {
    document.getElementById("chances").innerHTML=chances;
    table = document.getElementById("sudoku");
    table.style.borderTop = "thick solid black"
    table.style.borderLeft = "thick solid black"
    for (i = 0; i < 9; i++) {
        row = document.createElement("tr");
        table.appendChild(row);
        for (j = 0; j < 9; j++) {

            col = document.createElement("td");
            col.id = [i, j].join('');
            //console.log(col.id)
            row.appendChild(col);
            col.innerHTML = "a";
            if (((i + 1) % 3 == 0) && ((j + 1) % 3 == 0)) {
                document.getElementById(col.id).style.borderBottom = "thick solid black";
                document.getElementById(col.id).style.borderRight = "thick solid black";
            }
            if ((i + 1) % 3 == 0) {
                document.getElementById(col.id).style.borderBottom = "thick solid black";
            }
            if ((j + 1) % 3 == 0) {
                document.getElementById(col.id).style.borderRight = "thick solid black";
            }
        }
    }
    //alert("gridgenerated")
}
window.onload = function () {
    generateGrid();
    if (window["WebSocket"]) {
        loadGame(0)
    } else {
        console.log("Your Browser does not websocket")
        alert("Your Browser does not websocket")
    }
}

function loadGame(level) {
    con = new WebSocket("ws://" + document.location.host + "/ws");
    console.log("load game", level)
    con.close = function () {
        console.log("Connection with Server closed")
    }
    con.onopen = () => con.send(level)
    con.onmessage = function (server_data) {

        if (JSON.parse(server_data.data)) {

            var obj = JSON.parse(server_data.data);
            // console.log(obj[0],obj[1].key)
            loadGrid(obj)

        }
    }




}

function loadGrid(data) {

    data.forEach(element => {
        if (element.value == 0) {
            console.log(element.key)
            document.getElementById(element.key).innerHTML = "<input type='text' class='inVal' value='' onchange='checkInput(\"" + element.key  + "\",value)'/>";
        }
        else {
            document.getElementById(element.key).innerHTML = element.value
        }

    });

    document.getElementById("sudoku").style.visibility = "visible";
}

function checkInput(key, value) {
    console.log(key)
    if (value == '') {
        document.getElementById(key).style.backgroundColor = "white"
        document.getElementById(key).children[0].style.backgroundColor = "white"

    }
    else {
        if (parseInt(value) >= 1 && parseInt(value) <= 9) {
            const req_data = new Int8Array(2)
            req_data[0] = parseInt(key);
            req_data[1] = parseInt(value)
            con.send(req_data)
            con.onmessage = function (server_data) {
                if (server_data.data == 'invalid') {
                    chances=chances-1;
                    document.getElementById("chances").innerHTML=chances;
                    document.getElementById(key).children[0].style.backgroundColor = "red"
                    document.getElementById(key).style.backgroundColor = "red"
                }
                else {
                    cell = document.getElementById(key);
                    cell.children[0].value = value;
                    document.getElementById(key).children[0].disabled= "true"

                }

                if(server_data.data == 'Won'){
                    var modal = document.getElementById("myModal")
                    var span = document.getElementsByClassName("close")[0]
                    modal.style.display = "block";
                    document.getElementById("message").innerHTML="Congratulations!!!! You Won"
                    

                    span.onclick = function() {
                        modal.style.display = "none";
                      }
                }

                if(server_data.data=='loss'){
                    chances=chances-1;
                    document.getElementById("chances").innerHTML=chances;
                    var modal = document.getElementById("myModal");
                    var span = document.getElementsByClassName("close")[0];
                    modal.style.display = "block";
                    document.getElementById("message").innerHTML="Sorry Try Again!!!"
                    

                    span.onclick = function() {
                        modal.style.display = "none";
                      }

                }

            }


        }

    }

}

function changeDifficulty() {
    var level = document.getElementById("level").value

    loadGame(level)

}