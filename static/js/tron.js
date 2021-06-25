const GRID_WIDTH = 200;
const GRID_HEIGHT = 120;
const TILE_SIZE = 5;

var prevBoard = new Array(GRID_WIDTH*GRID_HEIGHT+1).join('~');

window.onload = function () {
    var conn;
    
    const UP_ARROW = 38;
    const DOWN_ARROW = 40;
    const LEFT_ARROW = 37;
    const RIGHT_ARROW = 39;
    const SPACE = 32;

    canvas = document.createElement('canvas');
    ctx = canvas.getContext('2d');
    container = document.getElementById('tron');
    container.append(canvas);
    canvas.width = GRID_WIDTH * TILE_SIZE;
    canvas.height = GRID_HEIGHT * TILE_SIZE;

    document.onkeydown = function(e) {
        switch(e.keyCode) {
            case UP_ARROW:
                conn.send("UP");
                break;
            case DOWN_ARROW:
                conn.send("DOWN");
                break;
            case LEFT_ARROW:
                conn.send("LEFT");
                break;
            case RIGHT_ARROW:
                conn.send("RIGHT");
                break;
            case SPACE:
                conn.send("BOOST");
                break;
        }
    }

    if (window["WebSocket"]) {
        const queryString = window.location.search;
        const urlParams = new URLSearchParams(queryString);
        const gameID = urlParams.get('gameID');
        conn = new WebSocket(`ws://${document.location.host}/game/ws/${gameID}`);
        conn.onclose = function (evt) {
            console.log("Connection closed...")
        };
        conn.onmessage = function (evt) {
            var event = JSON.parse(evt.data);
            switch(event.Kind) {
                case "STATE_UPDATE":
                    drawGrid(ctx, event.Data.grid);
                    prevBoard = event.Data.grid;
                    break;
                case "BEGIN":
                    break;

            }
        };
    } else {
        console.log("Browser does not support WebSockets...")
    }
};

function drawGrid(ctx, board) {
    colors = {
        '9': '#FFFFFF', // WALL
        '8': '#120458', // EMPTY
        // PLAYER COLORS
        '0': '#FF124F',
        '1': '#01DFF2',
        '2': '#FF00A0',
        '3': '#00DDFF',
        '4': '#FE75FE',
        '5': '#00FFEE',
    }
    for(i = 0; i < GRID_HEIGHT; i++) {
        for(j = 0; j < GRID_WIDTH; j++) {
            var index = i * GRID_WIDTH + j;
            var prevItem = prevBoard.charAt(index);
            var item = board.charAt(index);
            // bug where WALLS don't get drawn... I'd blame JS but I think it's on me :)
            if(item === prevItem && item != '9') {
                continue
            }
            ctx.fillStyle = colors[item];
            ctx.fillRect(
                TILE_SIZE * j,
                TILE_SIZE * i,
                TILE_SIZE,
                TILE_SIZE
            );
        }
    }
}
