package tron

type event struct {
	Kind eventType
	Data map[string]interface{}
}

type eventType string

const (
	// client events
	EventUp    eventType = "UP"
	EventDown  eventType = "DOWN"
	EventLeft  eventType = "LEFT"
	EventRight eventType = "RIGHT"
	EventBoost eventType = "BOOST"

	// server events
	EventStateUpdate eventType = "STATE_UPDATE"
	EventDeath       eventType = "DEATH"
	EventJoin        eventType = "JOIN"
	EventDisconnect  eventType = "DISCONNECT"
	EventBegin       eventType = "BEGIN"
)
