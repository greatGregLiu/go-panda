package live

import (
	"encoding/json"
	"time"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Node struct {
	Name    string                 `json:"-"`
	Type    string                 `json:"type,omitempty"`
	Sources []string               `json:"sources,omitempty"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

type Nodes map[string]Node

func (nodes Nodes) hackNames() {
	for name, node := range nodes {
		node.Name = name
		nodes[name] = node
	}
}

func (nodes Nodes) MarshalJSON() ([]byte, error) {
	nodes.hackNames()
	p, err := json.Marshal((map[string]Node)(nodes))
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (nodes *Nodes) UnmarshalJSON(p []byte) error {
	if err := json.Unmarshal(p, (*map[string]Node)(nodes)); err != nil {
		return err
	}
	nodes.hackNames()
	return nil
}

type Profile struct {
	ProfileID string     `json:"profile_id"`
	AccountID string     `json:"account_id"`
	Duration  int        `json:"duration,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Nodes     Nodes      `json:"nodes"`
}

type Stream struct {
	AccountID   string            `json:"account_id"`
	StreamID    string            `json:"stream_id"`
	ProfileID   string            `json:"profile_id"`
	Endpoints   map[string]string `json:"endpoints,omitempty"`
	CreatedAt   *time.Time        `json:"created_at"`
	ScheduledAt *time.Time        `json:"scheduled_at,omitempty"`
	StartedAt   *time.Time        `json:"started_at,omitempty"`
	EndedAt     *time.Time        `json:"ended_at,omitempty"`
	Error       *Error            `json:"error,omitempty"`
	Duration    int               `json:"duration"`
	CPU         int               `json:"cpu,omitempty"`
	Status      State             `json:"status"`
}

type State uint8

const (
	StateInvalid    = State(0)
	StateNew        = State(1)
	StateQueued     = State(2)
	StatePending    = State(3)
	StateInProgress = State(4)
	StateEnded      = State(5)
	StateError      = State(6)
	StateReady      = State(7)
)

var stateNames = map[State]string{
	StateInvalid:    "invalid",
	StateNew:        "new",
	StateQueued:     "queued",
	StatePending:    "pending",
	StateInProgress: "in progress",
	StateEnded:      "ended",
	StateError:      "error",
	StateReady:      "ready",
}

func (s State) String() string {
	return stateNames[s]
}

type postResp struct {
	ProfileID string `json:"profile_id"`
	StreamID  string `json:"stream_id"`
}
