package environment

import (
	"fmt"
	"math/rand"
)

// Environments представляет среду для агентов
type Environments struct {
	agentCount      int
	average         float64
	agentNeighbours map[int][]AgentLink
	maxDelta        float64
	noise           float64
	debug           bool
}

type AgentLink struct {
	ID          int
	Channel     *chan interface{}
	Reliability float64
}

func NewEnvironments(
	agentCount int, average float64, agentNeighbours map[int][]AgentLink, maxDelta float64, noise float64, debug bool,
) Environments {
	return Environments{
		agentCount:      agentCount,
		average:         average,
		agentNeighbours: agentNeighbours,
		maxDelta:        maxDelta,
		noise:           noise,
		debug:           debug,
	}
}

func (env *Environments) AddAgentLinks(agentID int, link []AgentLink) {
	if env.agentNeighbours == nil {
		env.agentNeighbours = make(map[int][]AgentLink)
	}
	env.agentNeighbours[agentID] = append(
		env.agentNeighbours[agentID], link...,
	)
}

func (env *Environments) SetAverage(average float64) {
	env.average = average
}

func (env *Environments) Average() float64 {
	return env.average
}

func (env *Environments) Neighbours(agentID int) []AgentLink {
	toReturn := make([]AgentLink, 0, len(env.agentNeighbours[agentID]))
	for _, neighbour := range env.agentNeighbours[agentID] {
		if neighbour.Reliability >= rand.Float64() {
			toReturn = append(toReturn, neighbour)
		} else if env.debug {
			fmt.Printf("unlink %d -- %d\n", agentID, neighbour.ID)
		}
	}
	return toReturn
}

func (env *Environments) AgentCount() int {
	return env.agentCount
}

func (env *Environments) MaxDelta() float64 {
	return env.maxDelta
}

func (env *Environments) Noise() float64 {
	return env.noise
}

func (env *Environments) Debug() bool {
	return env.debug
}
