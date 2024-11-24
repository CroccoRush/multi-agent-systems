package message

// Request представляет запрос на получение значений
type Request struct {
	from int
}

func NewRequest(from int) *Request {
	return &Request{from: from}
}

func (r *Request) From() int {
	return r.from
}

// Finish представляет флаг о завершении подсчёта среднего
type Finish struct{}
