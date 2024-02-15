// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")
var house = make(map[string]*Hub)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/{room}", serveHome)
	r.HandleFunc("/ws/{room}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		roomID := vars["room"]

		mutex.Lock()
		room, ok := house[roomID]

		//fmt.Println("Sleep 10 seconds")
		//time.Sleep(time.Second * 10)

		var hub *Hub
		if ok {
			fmt.Println("Found room")
			hub = room
		} else {
			fmt.Println("Create room")
			hub = newHub(roomID)
			house[roomID] = hub
			go hub.run()
		}
		serveWs(hub, writer, request)
	})
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           r,
	}
	panic(server.ListenAndServe())
}
