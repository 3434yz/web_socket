// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")

// var house = make(map[string]*Hub)
var house sync.Map
var roomMutexes = make(map[string]*sync.Mutex)
var mutexForRoomMutexes = new(sync.Mutex)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func echo(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	roomID := vars["room"]

	mutexForRoomMutexes.Lock()
	roomMutex, ok := roomMutexes[roomID]
	if ok {
		roomMutex.Lock()
	} else {
		roomMutexes[roomID] = new(sync.Mutex)
		roomMutexes[roomID].Lock()
	}
	mutexForRoomMutexes.Unlock()

	//room, ok := house[roomID]
	room, ok := house.Load(roomID)
	//fmt.Println("Sleep 10 seconds")
	//time.Sleep(time.Second * 10)

	var hub *Hub
	if ok {
		fmt.Println("Found room")
		hub = room.(*Hub)
	} else {
		fmt.Println("Create room")
		hub = newHub(roomID)
		//house[roomID] = hub
		house.Store(roomID, hub)
		go hub.run()
	}
	serveWs(hub, writer, request)
}

type Auth struct {
	Next http.Handler
}

func NewAuth(handler http.Handler) *Auth {
	return &Auth{Next: handler}
}

func (a *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.Next == nil {
		a.Next = http.DefaultServeMux
	}
	//auth := r.Header.Get("Authen")
	//if auth == "" {
	//	w.WriteHeader(http.StatusUnauthorized)
	//} else {
	//	a.Next.ServeHTTP(w, r)
	//}
	fmt.Println("------------------------------------------------------")
}

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/{room}", serveHome)
	r.HandleFunc("/ws/{room}", echo)
	r.Handle("/ws/{room}", NewAuth(r))
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           r,
	}
	panic(server.ListenAndServe())
}
