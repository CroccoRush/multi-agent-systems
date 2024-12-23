package agent

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// SendValue рассылает собственное значение всем соседям
func (a *Agent) SendValue() {
	if a.myLeader != -1 {
		return
	}
	for _, channel := range a.neighbors {
		val := a.value + rand.NormFloat64()*a.env.Noise()
		if a.env.Debug() {
			fmt.Printf("real: %f, noisy: %f\n", a.value, val)
		}
		channel <- val
	}
}

// handleMessage обрабатывает сообщение со значением
func (a *Agent) handleMessage(value float64) {
	delta := value - a.value
	a.sum += delta
	a.count++
	if a.count == len(a.neighbors) {
		ok := a.calculateAverage()
		if ok && math.Abs(delta) < math.Abs(a.env.MaxDelta()) && a.tick > 10 {
			a.initVoting() // Заявляет себя как кандидата к отправке
		}
	}
}

// calculateAverage вычисляет среднее значение и выводит его
func (a *Agent) calculateAverage() bool {
	delta := a.alpha * a.sum
	a.value += delta
	a.sum = 0
	a.count = 0
	return math.Abs(delta) < math.Abs(a.env.MaxDelta())
}
