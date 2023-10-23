package randomtransaction

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"transactions/randomdata"
	"transactions/transaction"
)

var (
	openFile        = os.OpenFile
	jsonMarshall    = json.Marshal
	fileWriteString = func(file *os.File, s string) (n int, err error) {
		return file.WriteString(s)
	}

	printToLog = func(v ...any) {
		log.Println(v...)
	}
)

type Worker struct {
	FilePath  string
	MinAmount float32
	MaxAmount float32
}

func (w *Worker) Work(ctx context.Context) {
	t := generateRandomTransaction(w.MinAmount, w.MaxAmount)
	file, err := openFile(w.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		printToLog(err)
		return
	}
	defer file.Close()
	jsonData, err := jsonMarshall(t)
	if err != nil {
		printToLog(err)
		return
	}
	_, err = fileWriteString(file, string(jsonData)+"\n")
	if err != nil {
		printToLog("error writing to file: ", err)
	}
}

func generateRandomTransaction(minAmount, maxAmount float32) *transaction.Transaction {
	const withdrawal = "withdrawal"
	t := &transaction.Transaction{
		TransactionId:     randomdata.TransactionId(),
		TransactionType:   withdrawal,
		AccountNumber:     randomdata.AccountNumber(),
		TransactionAmount: randomdata.TransactionAmount(minAmount, maxAmount),
		TransactionTime:   randomdata.TransactionTime(),
		Location:          randomdata.Location(),
	}
	return t
}
