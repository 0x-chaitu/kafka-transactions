package kafka

import (
	"context"
	"log"
	"transactions/mongodb"
	"transactions/mongodb/suspicioustransaction"
	"transactions/mongodb/suspicioustransaction/models"
	"transactions/transaction"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	printToLog = func(v ...any) {
		log.Println(v...)
	}
	stInsert = func(ctx context.Context, db *mongodb.Mongodb, sp *models.SuspiciousTransaction) error {
		return suspicioustransaction.Insert(ctx, db, sp)
	}
)

type Worker struct {
	Msg *kafka.Message
	Db  *mongodb.Mongodb
}

func (c *Worker) insertSuspiciousTransaction(ctx context.Context, sp *transaction.Transaction) error {
	spDb := &models.SuspiciousTransaction{
		TransactionId:     sp.TransactionId,
		AccountNumber:     sp.AccountNumber,
		TransactionType:   sp.TransactionType,
		TransactionAmount: sp.TransactionAmount,
		TransactionTime:   sp.TransactionTime,
		Location:          sp.Location,
	}

	return stInsert(ctx, c.Db, spDb)
}

func (w *Worker) Work(ctx context.Context) {
	transaction, err := transaction.New(string(w.Msg.Value))
	if err != nil {
		printToLog(err)
		return
	}
	if transaction.IsSuspicious() {
		printToLog("suspicios transaction")
		if err := w.insertSuspiciousTransaction(ctx, transaction); err != nil {
			printToLog("error when inserting suspicious transaction")
		}
	}
}
