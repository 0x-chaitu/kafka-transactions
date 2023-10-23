package transaction

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	TransactionId     int       `json:"transaction_id"`
	AccountNumber     int       `json:"account_number"`
	TransactionType   string    `json:"transaction_type"`
	TransactionAmount float32   `json:"transaction_amount"`
	TransactionTime   time.Time `json:"transaction_time"`
	Location          string    `json:"location"`
}

func New(rawTransaction string) (*Transaction, error) {
	t := new(Transaction)
	if err := json.Unmarshal([]byte(rawTransaction), &t); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Transaction) IsSuspicious() bool {
	const suspiciousAmount = float32(10_000)
	return t.TransactionAmount > suspiciousAmount
}
