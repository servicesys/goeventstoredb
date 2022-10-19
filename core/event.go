package core

import (
	"crypto/rand"
	"math/big"
	"time"
)

type Event struct {
	EventID       string
	EventType     EventType
	DomainTenant  string
	AppName       string
	TransactionID string
	EventVersion  string
	TimeStamp     string
	AggregateID   int64
	AggregateType string
	Payload       []byte
	CreatedAt     time.Time
	UserID        string
}

type EventType struct {
	ID        string
	MetaData  []byte
	CreatedAt time.Time
}

func GenerateTransactionID() string {
	//Max random value, a 130-bits integer, i.e 2^130 - 1
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	//Generate cryptographically strong pseudo-random between 0 - max
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		//error handling
	}

	//String representation of n in base 32
	nonce := n.Text(10)
	return nonce
}
