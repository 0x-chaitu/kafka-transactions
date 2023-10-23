package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"transactions/config"
	"transactions/mongodb"
	"transactions/task"
	kafkaworker "transactions/task/worker/kafka"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	bootstrapServersKey   = "bootstrap.servers"
	groupIdKey            = "group.id"
	autoOffsetResetKey    = "auto.offset.reset"
	autoOffsetReset       = "earliest"
	enablePartitionEofKey = "enable.partition.eof"
)

func run() error {
	const envFile = ".env"
	defer log.Println("main: completed")
	ctx := context.Background()

	cfg, err := config.Read(envFile)
	if err != nil {
		return errors.Unwrap(err)
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		bootstrapServersKey:   cfg.KafkaBrokerHost,
		groupIdKey:            cfg.KafkaGroupId,
		autoOffsetResetKey:    autoOffsetReset,
		enablePartitionEofKey: false,
	})

	if err != nil {
		return errors.Join(err, errors.New("connecting to broker"+cfg.KafkaBrokerHost))
	}
	if err := consumer.SubscribeTopics([]string{cfg.KafkaTopic}, nil); err != nil {
		return errors.Join(err, errors.New("subscribing to kafka topic"+cfg.KafkaTopic))
	}

	db, err := mongodb.Connect(ctx, cfg.MongodbHostName, cfg.MongodbDatabase, cfg.MongodbPort)
	if err != nil {
		return errors.Unwrap(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	maxGoRoutines := runtime.GOMAXPROCS(0)
	pool := task.New(ctx, maxGoRoutines)

	go func() {
		defer close(shutdown)
		defer close(serverErrors)
		for {
			select {
			case <-shutdown:
				log.Println("run: Start shutdown")
				if err := consumer.Close(); err != nil {
					serverErrors <- errors.Unwrap(err)
				}
				return
			default:
				msg, err := consumer.ReadMessage(-1)
				if err != nil {
					serverErrors <- err
				} else {
					kw := &kafkaworker.Worker{
						Msg: msg,
						Db:  db,
					}
					pool.Do(kw)
				}
			}
		}
	}()

	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdown:
		log.Printf("run: %v start shutdown", sig)
		if err := consumer.Close(); err != nil {
			return errors.Unwrap(err)
		}
		return nil
	}
}

func main() {
	const logFileName = "logs/consumer.txt"
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
