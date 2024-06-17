// server.go
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	// The event loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)
		for i := 0; i < 10; i++ {
			go func() {
				err = conn.WriteMessage(messageType, message)
				if err != nil {
					log.Println("Error during message writing:", err)
					return
				}
			}()
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index Page")
}

func main() {
	//http.HandleFunc("/socket", socketHandler)
	//http.HandleFunc("/", home)
	//log.Fatal(http.ListenAndServe("localhost:8000", nil))

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println("当前协程数量", runtime.NumGoroutine())
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/socket", socketHandler)
	r.HandleFunc("/", home)

	server := &http.Server{
		Addr:              "localhost:8080",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           r,
	}

	go func() {
		server.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
