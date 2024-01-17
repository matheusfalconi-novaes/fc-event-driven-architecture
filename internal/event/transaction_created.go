package event

import "time"

type TransactionCreated struct {
	Name     string
	Payload  interface{}
	DateTime time.Time
}

func NewTransactionCreated() *TransactionCreated {
	return &TransactionCreated{
		Name:     "TransactionCreated",
		DateTime: time.Now(),
	}
}

func (e *TransactionCreated) GetName() string {
	return e.Name
}

func (e *TransactionCreated) GetDateTime() time.Time {
	return e.DateTime
}

func (e *TransactionCreated) GetPayload() interface{} {
	return e.Payload
}

func (e *TransactionCreated) SetPayload(payload interface{}) {
	e.Payload = payload
}
