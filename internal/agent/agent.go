package agent

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	env "multi-agent-systems/t1/internal/environment"
	"multi-agent-systems/t1/internal/message"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Agent представляет агента
type Agent struct {
	id         int
	channel    chan interface{}
	neighbors  map[int]chan interface{}
	agentCount int
	values     map[int]float64
}

// NewAgent создает нового агента
func NewAgent(id int, number *float64) *Agent {
	values := make(map[int]float64)
	if number != nil {
		values[id] = *number
	} else {
		values[id] = rand.Float64() * 100
	}
	return &Agent{
		id:         id,
		channel:    make(chan interface{}, 100),
		neighbors:  make(map[int]chan interface{}),
		agentCount: 0,
		values:     values,
	}
}

// OwnValue возвращает собственное число
func (a *Agent) OwnValue() float64 {
	return a.values[a.id]
}

// Channel возвращает свой канал для связи
func (a *Agent) Channel() *chan interface{} {
	return &a.channel
}

// Run запускает агента в отдельной горутине
func (a *Agent) Run(wg *sync.WaitGroup, env *env.Environments) {
	for _, neighbour := range env.Neighbours(a.id) {
		a.neighbors[neighbour.ID] = *neighbour.Channel
	}
	a.agentCount = env.AgentCount()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer wg.Done()
	for {
		select {
		case msg := <-a.channel:
			switch m := msg.(type) {
			case map[int]float64:
				finish := a.handleMessage(m)
				if finish {
					return
				}
			case message.Request:
				a.handleRequest(m)
			case message.Finish:
				a.handleFinish() // Передаём сигнал о завершении всем остальным
				return
			}
		case <-ticker.C:
			//if a.id == 0 {
			//	fmt.Println("tick")
			//}
			a.SendRequests()
		}
	}
}

// SendRequests отправляет запросы на получение значений всем соседям
func (a *Agent) SendRequests() {
	for _, channel := range a.neighbors {
		channel <- *message.NewRequest(a.id)
	}
}

// handleRequest обрабатывает запрос на получение значений
func (a *Agent) handleRequest(req message.Request) {
	values := make(map[int]float64)
	for id, value := range a.values {
		values[id] = value
	}
	a.neighbors[req.From()] <- values
}

// handleMessage обрабатывает сообщение
func (a *Agent) handleMessage(message map[int]float64) bool {
	for id, value := range message {
		a.values[id] = value
	}
	if a.id == 0 && len(a.values) == a.agentCount {
		a.handleFinish() // Отправляем сигнал о завершении
		a.calculateAverage()
		return true
	}
	return false
}

// handleFinish обрабатывает запрос на завершение вычислений
func (a *Agent) handleFinish() {
	for _, neighbor := range a.neighbors {
		neighbor <- message.Finish{}
	}
}

// calculateAverage вычисляет среднее значение и выводит его
func (a *Agent) calculateAverage() {
	sum := 0.0
	for _, value := range a.values {
		sum += value
	}
	average := sum / float64(a.agentCount)
	fmt.Printf("Agent %d: Average = %f\n", a.id, average)
}
