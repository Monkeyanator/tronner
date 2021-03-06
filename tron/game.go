package tron

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	gridWidth    = 200
	gridHeight   = 120
	tickInterval = 20 * time.Millisecond
)

// gid is a unqiue ID assigned to a player per-game.
type gid uint8

// game represents a running instance of a Tron game.
type Game struct {
	id string

	players     map[gid]*Player
	playersLock sync.RWMutex

	// game state
	grid     *grid
	finished bool

	broadcast chan event
}

func NewGame() *Game {
	g := &Game{
		id:          uuid.NewString(),
		players:     make(map[gid]*Player, 0),
		playersLock: sync.RWMutex{},
		grid:        newGrid(gridWidth, gridHeight),
		finished:    false,
		broadcast:   make(chan event, 512),
	}

	go g.broadcastLoop()
	return g
}

// broadcastLoop sends event to connected players.
func (g *Game) broadcastLoop() {
	for {
		select {
		case e := <-g.broadcast:
			g.playersLock.Lock()
			for _, p := range g.players {
				p.conn.send(event(e))
			}
			g.playersLock.Unlock()
		}
	}
}

func (g *Game) Add(p *Player) {
	g.playersLock.Lock()
	defer g.playersLock.Unlock()
	p.gid = g.nextGID()
	g.players[p.gid] = p
	log.Printf("Player %s has joined game %s with GID %d", p.id, g.id, p.gid)
	go g.readEvents(p)
}

func (g *Game) Remove(p *Player) {
	g.playersLock.Lock()
	defer g.playersLock.Unlock()
	log.Printf("Player %s has left game %s", g.id, p.id)
	delete(g.players, p.gid)
}

func (g *Game) readEvents(p *Player) {
	for e := range p.conn.recv() {
		g.handleEvent(p, e)
	}
	// if we break out of loop the websocket connection was terminated
	// and we should remove the player from the game.
	g.Remove(p)
}

func (g *Game) handleEvent(p *Player, e event) {
	g.playersLock.Lock()
	defer g.playersLock.Unlock()
	switch e.Kind {
	case EventUp:
		if p.direction != DOWN {
			p.direction = UP
		}
	case EventDown:
		if p.direction != UP {
			p.direction = DOWN
		}
	case EventLeft:
		if p.direction != RIGHT {
			p.direction = LEFT
		}
	case EventRight:
		if p.direction != LEFT {
			p.direction = RIGHT
		}
	case EventBoost:
		p.boostTicks = timeToTicks(time.Second * 2)
	default:
		log.Printf("Cannot handle unknown event %s", e)
	}
}

func (g *Game) initialize() {
	g.grid.initialize()
	g.arrangePlayers()
	for _, p := range g.players {
		p.alive = true
		p.boostTicks = 0
	}
	g.finished = false
	g.broadcast <- event{
		Kind: EventBegin,
	}
}

func (g *Game) Run() {
	g.initialize()
	ticker := time.NewTicker(tickInterval)
	for {
		select {
		case <-ticker.C:
			events := g.tick()
			// consider delivering events generated from same tick as a unit...
			// might run into weird ordering issues otherwise
			for _, e := range events {
				g.broadcast <- e
			}
			if g.finished {
				ticker.Stop()
				<-time.After(time.Second)
				g.initialize()
				ticker.Reset(tickInterval)
			}
		}
	}
}

func (g *Game) tick() []event {
	events := []event{}
	remaining := 0
	for _, p := range g.players {
		if !p.alive {
			continue
		}
		var (
			// unit direction change for x, y
			dx, dy = 0, 0
			// number of units to move in that direction this tick
			xd = 1
		)
		if p.boostTicks > 0 {
			xd += 1
			p.boostTicks--
		}
		switch p.direction {
		case UP:
			dy = -1
		case DOWN:
			dy = 1
		case LEFT:
			dx = -1
		case RIGHT:
			dx = 1
		}
		var newX, newY int
		for d := 1; d <= xd; d++ {
			newX = p.x + dx*d
			newY = p.y + dy*d
			dest := g.grid.get(newX, newY)
			if dest != EMPTY {
				g.grid.blast(newX, newY, 2)
				p.alive = false
				events = append(events, event{
					Kind: EventDeath,
					Data: map[string]interface{}{
						"gid": p.gid,
					},
				})
				break
			}
			g.grid.set(newX, newY, uint8(p.gid))
		}
		if p.alive {
			remaining++
			p.x = newX
			p.y = newY
		}
	}

	// Once working change to 1
	if remaining == 0 {
		g.finished = true
	}

	events = append(events, event{
		Kind: EventStateUpdate,
		Data: map[string]interface{}{
			"grid": g.grid.serialize(),
		}})
	return events
}

/* Provides initial arrangement for the player bikes. Ought to be like:
_______________
|             |
|->         <-|
|             |
|->         <-|
|             |
|->         <-|
|_____________|
*/
func (g *Game) arrangePlayers() {
	g.playersLock.Lock()
	defer g.playersLock.Unlock()
	offset := 20
	for gid, p := range g.players {
		var x, y int
		if gid%2 == 0 {
			p.direction = RIGHT
			x = offset
		} else {
			p.direction = LEFT
			x = gridWidth - offset
		}
		if len(g.players) <= 2 {
			y = gridHeight / 2 * (int(gid)/2 + 1)
		} else {
			y = (gridHeight/len(g.players) + 1/2) * (int(gid)/2 + 1)
		}
		p.x = x
		p.y = y
	}
}

func (g *Game) nextGID() gid {
	i := 0
	for {
		if _, ok := g.players[gid(i)]; !ok {
			return gid(i)
		}
		i++
	}
}

func timeToTicks(t time.Duration) uint {
	return uint(t / tickInterval)
}
