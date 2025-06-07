package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	hook "github.com/robotn/gohook"
	"golang.org/x/net/websocket"
)

type Coord struct {
	X, Y int
}

var C Coord

func main() {
	go func() {
		add()
		low()
	}()
	http.HandleFunc("/", admin)
	http.HandleFunc("/client", client)
	http.Handle("/stream", websocket.Handler(stream))
	http.ListenAndServe(":1234", nil)
}

func add() {
	fmt.Println("To start press q")
	hook.Register(hook.KeyDown, []string{"q"}, func(e hook.Event) {
		hook.End()
	})

	s := hook.Start()
	<-hook.Process(s)
}

func low() {
	evChan := hook.Start()
	defer hook.End()
	for ev := range evChan {
		if ev.Kind == 10 {
			C.X = int(ev.X)
			C.Y = int(ev.Y)
		}
	}
}

func admin(w http.ResponseWriter, r *http.Request) {
	adminTemplate.Execute(w, "127.0.0.1:1234")
}

func client(w http.ResponseWriter, r *http.Request) {
	clientTemplate.Execute(w, "ws://127.0.0.1:1234/stream")
}

func stream(ws *websocket.Conn) {
	var err error

	for {
		msg := fmt.Sprint(C.X) + " " + fmt.Sprint(C.Y)
		fmt.Println("coord: ", msg)
		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

var adminTemplate = template.Must(template.New("").Parse(`
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-+0n0xVW2eSR5OomGNYDnhzAbDsOXxcvSN1TPprVMTNDbiYZCxYbOOl7+AMvyTG2x" crossorigin="anonymous">
    <title>Drag</title>
  </head>
  <body>
    <style>
      .block1 {
    width: 100px;
    height: 100px;
    position: absolute;
    left: 50px;
    top: 50px;
    background: red;
    }
    </style>

    <div class="block1" id="square"></div>

    <script type="text/javascript">
        var div = document.getElementById('square');
        var listener = function(e) {
          div.style.left = e.pageX - 50 + "px";
          div.style.top = e.pageY - 50 + "px";
        };

        square.addEventListener('mousedown', e => {
            document.addEventListener('mousemove', listener);
        });

        square.addEventListener('mouseup', e => {
            document.removeEventListener('mousemove', listener);
        });
    </script>
  </body>
</html>
`))

var clientTemplate = template.Must(template.New("").Parse(`
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-+0n0xVW2eSR5OomGNYDnhzAbDsOXxcvSN1TPprVMTNDbiYZCxYbOOl7+AMvyTG2x" crossorigin="anonymous">
    <title>Stream</title>
  </head>
  <body>
    <style>
      .block1 {
    width: 100px;
    height: 100px;
    position: absolute;
    left: 50px;
    top: 50px;
    background: red;
    }
    </style>

    <div class="block1" id="square"></div>

    <script type="text/javascript">
        var sock = null;
        var wsuri = "ws://127.0.0.1:1234/stream";
        
        window.onload = function() {

            console.log("onload");

            sock = new WebSocket(wsuri);

            sock.onopen = function() {
                console.log("connected to " + wsuri);
            }
            
            sock.onmessage = function(e) {
                var div = document.getElementById('square');
                var arr = e.data.split(" ");
                div.style.left = arr[0] + "px";
                div.style.top = arr[1] + "px";
            }
        };
    </script>
  </body>
</html>
`))
