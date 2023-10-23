package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"transactions/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const bootstrapServersKey = "bootstrap.servers"

func stringPrt(s string) *string {
	return &s
}

func run(cfg *config.Config, transactionFile string) error {
	log.Println("main: Initializing Kafka!")
	defer log.Println("main: finished")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		bootstrapServersKey: cfg.KafkaBrokerHost,
	})
	if err != nil {
		return errors.Unwrap(err)
	}

	defer producer.Close()
	file, err := os.Open(transactionFile)
	if err != nil {
		return errors.Unwrap(err)
	}
	defer file.Close()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	deliverChan := make(chan kafka.Event)
	scanner := bufio.NewScanner(file)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if err := producer.Produce(
				&kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     stringPrt(cfg.KafkaTopic),
						Partition: kafka.PartitionAny,
					},
					Value: []byte(line),
				}, deliverChan,
			); err != nil {
				log.Println(err, cfg.KafkaTopic)
			}
			delivery := <-deliverChan
			m := delivery.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				log.Println("message failed")
			}

			if err := scanner.Err(); err != nil {
				errors.Unwrap(err)
			}
		}
	}()

	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdown:
		log.Println(sig)
		return nil
	}

}

func main() {
	const (
		envFile     = ".env"
		logFileName = "logs/producer.txt"
	)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|
		os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf(`opening file error%v`, logFile)
	}
	cfg, err := config.Read(envFile)
	if err != nil {
		log.Panicln(err)
		os.Exit(1)
	}
	if err := run(cfg, logFileName); err != nil {
		log.Panicln(err)
		os.Exit(1)
	}

}
