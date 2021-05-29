package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Monkeyanator/tronner/tron"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

// Server contains state required to run game.
type Server struct {
	httpServer *http.Server
	c          Config

	games map[string]*tron.Game
	mu    sync.RWMutex
}

type Config struct {
	// Port to run Tron server on.
	Port uint
}

func New(c Config) *Server {
	r := mux.NewRouter()
	httpServer := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", c.Port),
	}
	s := &Server{
		httpServer: httpServer,
		c:          c,
		games:      make(map[string]*tron.Game),
		mu:         sync.RWMutex{},
	}

	// Setup route handlers.
	r.HandleFunc("/game/ws/{id:[0-9]+}", s.gameHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	return s
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// gameHandler serves when client tries joining existing game ID.
func (s *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad ID in route: %s", r.URL.Path)
		return
	}
	s.mu.RLock()
	game, ok := s.games[id]
	s.mu.RUnlock()

	// For now, create a game if not exists
	// should be separate path or method on this path
	if !ok {
		s.mu.Lock()
		game = tron.NewGame()
		log.Printf("Creating game ID %s", id)
		s.games[id] = game
		s.mu.Unlock()
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading websocket: %v", err)
		writeError(w, err)
		return
	}
	p := tron.NewPlayer(tron.NewWebsocketConnection(c))
	game.Add(p)
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
