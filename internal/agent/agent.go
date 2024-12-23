package agent

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	env "multi-agent-systems/t1/internal/environment"
	"multi-agent-systems/t1/internal/message"
)

// Agent представляет агента
type Agent struct {
	id        int
	env       *env.Environments
	channel   chan interface{}
	neighbors map[int]chan interface{}

	// Поля для подсчёта среднего
	value float64
	sum   float64
	count int
	alpha float64

	// Поля для консенсуса
	myLeader       int
	myLeaderVoices int
	iAmLeader      bool
	tick           int
	ticksToWin     int
}

// NewAgent создает нового агента
func NewAgent(id int, number *float64, alpha float64, environments *env.Environments) *Agent {
	var value float64
	if number != nil {
		value = *number
	} else {
		value = rand.Float64() * 100
	}
	return &Agent{
		id:             id,
		env:            environments,
		channel:        make(chan interface{}, 100),
		neighbors:      make(map[int]chan interface{}),
		value:          value,
		sum:            0,
		count:          0,
		alpha:          alpha,
		myLeader:       -1,
		myLeaderVoices: 0,
		iAmLeader:      false,
		tick:           0,
		ticksToWin:     -1,
	}
}

// OwnValue возвращает собственное число
func (a *Agent) OwnValue() float64 {
	return a.value
}

// Channel возвращает свой канал для связи
func (a *Agent) Channel() *chan interface{} {
	return &a.channel
}

func (a *Agent) UpdateNeighbours() {
	for _, neighbour := range a.env.Neighbours(a.id) {
		a.neighbors[neighbour.ID] = *neighbour.Channel
	}
}

// Run запускает агента в отдельной горутине
func (a *Agent) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	a.UpdateNeighbours()
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case msg := <-a.channel:
			switch m := msg.(type) {
			case float64:
				a.handleMessage(m)
			case message.Voting:
				a.handleVoting(m)
			case message.Finish:
				a.handleFinish() // Передаём сигнал о завершении всем остальным
				return
			}
		case <-ticker.C:
			if a.id == 0 {
				fmt.Println("tick")
			}
			a.UpdateNeighbours()
			a.SendValue()
			a.tick++
			if a.iAmLeader && a.tick == a.ticksToWin {
				if a.env.Debug() {
					fmt.Printf("Agent %d: Average = %f\n", a.id, a.value)
				}
				a.handleFinish()
			}
		}
	}
}

// handleFinish рассылка сигнала к завершению
func (a *Agent) handleFinish() {
	for _, neighbor := range a.neighbors {
		neighbor <- message.Finish{}
	}
}
