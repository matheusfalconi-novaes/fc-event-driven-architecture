package event

import "time"

type BalanceUpdated struct {
	Name     string
	Payload  interface{}
	DateTime time.Time
}

func NewBalanceUpdated() *BalanceUpdated {
	return &BalanceUpdated{
		Name:     "BalanceUpdated",
		DateTime: time.Now(),
	}
}

func (e *BalanceUpdated) GetName() string {
	return e.Name
}

func (e *BalanceUpdated) GetPayload() interface{} {
	return e.Payload
}

func (e *BalanceUpdated) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *BalanceUpdated) GetDateTime() time.Time {
	return e.DateTime
}
