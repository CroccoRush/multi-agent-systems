package message

// Finish представляет флаг о завершении подсчёта среднего
type Finish struct{}

// Voting является сообщением для голосования за лидера для отправки среднего
type Voting struct {
	Leader int
	Voices int
}
