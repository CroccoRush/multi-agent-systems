package agent

import (
	"fmt"
	"multi-agent-systems/t1/internal/message"
)

func (a *Agent) changeLeader(m message.Voting) {
	m.Voices++
	a.iAmLeader = false
	a.myLeader = m.Leader
	a.myLeaderVoices = m.Voices
	for _, neighbor := range a.neighbors {
		neighbor <- m
	}
}

// handleVoting обрабатывает полученное сообщение о каком-то лидере
// Если он лучше, собственного, то агент поменяет мнение и будет продвигать нового лидер.
func (a *Agent) handleVoting(m message.Voting) {
	if a.myLeader == m.Leader {
		return
	}
	if a.myLeaderVoices < m.Voices || (a.myLeaderVoices == m.Voices && a.myLeader < m.Leader) {
		a.changeLeader(m)
		if a.env.Debug() {
			fmt.Printf("VOICE FOR LEADER: Agent %d: score %d\n", a.myLeader, a.myLeaderVoices)
		}
	}
}

// initVoting заявляет себя как лидера
func (a *Agent) initVoting() {
	a.iAmLeader = true
	a.myLeader = a.id
	a.myLeaderVoices = 1
	a.ticksToWin = int(float64(a.tick) * a.alpha)
	a.tick = 0
	if a.env.Debug() {
		fmt.Printf("LEADER CANDIDATE: Agent %d: ticksToWin %d\n", a.id, a.ticksToWin)
	}
	m := message.Voting{
		Leader: a.myLeader,
		Voices: a.myLeaderVoices,
	}
	for _, neighbor := range a.neighbors {
		neighbor <- m
	}
}
