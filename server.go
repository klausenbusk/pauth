package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

//go:embed web/index.html web/script.js
var f embed.FS

type authServer struct {
	cors     string
	serveMux http.ServeMux

	clientsMutex sync.Mutex
	clients      map[uuid.UUID](chan string)
}

func newServer(cors string) *authServer {
	s := &authServer{
		cors:    cors,
		clients: make(map[uuid.UUID](chan string)),
	}

	webFs, _ := fs.Sub(f, "web")
	s.serveMux.Handle("/", http.FileServer(http.FS(webFs)))
	s.serveMux.HandleFunc("/ws", s.wsHandler)
	s.serveMux.HandleFunc("/uuid", s.uuidHandler)
	return s
}

func (s *authServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *authServer) ListenAndServe(addr string) {
	httpServer := &http.Server{
		Handler:      s,
		Addr:         addr,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	fmt.Printf("Listening on: %s\n", addr)
	log.Fatal(httpServer.ListenAndServe())
}

func (s *authServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	typ, buf, err := c.Read(r.Context())
	if err != nil {
		return
	}
	if typ == websocket.MessageBinary {
		c.Close(websocket.StatusUnsupportedData, "")
		return
	}

	id, err := uuid.Parse(string(buf))
	if err != nil {
		c.Close(4000, "Invalid UUID")
		return
	}

	ch := make(chan string)
	s.clientsMutex.Lock()
	s.clients[id] = ch
	s.clientsMutex.Unlock()

	go func() {
		for {
			typ, buf, err := c.Read(r.Context())
			if err != nil {
				break
			}
			if typ == websocket.MessageBinary {
				c.Close(websocket.StatusUnsupportedData, "")
				break
			}

			arr := strings.SplitN(string(buf), ",", 2)
			if len(arr) != 2 {
				c.Close(4001, "Invalid message")
				break
			}

			destId, err := uuid.Parse(arr[0])
			if err != nil {
				c.Close(4000, "Invalid UUID")
				break
			}

			s.clientsMutex.Lock()
			s.clients[destId] <- fmt.Sprintf("%s,%s", id, arr[1])
			s.clientsMutex.Unlock()
		}
	}()

	for {
		select {
		case b := <-ch:
			c.Write(r.Context(), websocket.MessageText, []byte(b))
		case <-r.Context().Done():
			return
		}
	}
}

func (s *authServer) uuidHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, id)
}
