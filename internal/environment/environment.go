package environment

// Environments представляет среду для агентов
type Environments struct {
	agentCount      int
	average         float64
	agentNeighbours map[int][]AgentLink
}

type AgentLink struct {
	ID      int
	Channel *chan interface{}
}

func NewEnvironments(
	agentCount int, average float64, agentNeighbours map[int][]AgentLink,
) Environments {
	return Environments{
		agentCount:      agentCount,
		average:         average,
		agentNeighbours: agentNeighbours,
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
	return env.agentNeighbours[agentID]
}

func (env *Environments) AgentCount() int {
	return env.agentCount
}
