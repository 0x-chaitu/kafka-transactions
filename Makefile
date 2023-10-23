.PHONY: sample-data
## sample-data: generates sample data
sample-data:
	@ if [ -z "$(TOTAL)" ]; then echo >&2 please set total via the variable TOTAL; exit 2; fi
	@ if [ -z "$(FILE_NAME)" ]; then echo >&2 please set file name via the variable FILE_NAME; exit 2; fi
	@ rm -f "${SAMPLE_DATA_FOLDER}/${FILE_NAME}"
	@ echo "generating file ${SAMPLE_DATA_FOLDER}/${FILE_NAME}..."
	@ go run jsongenerator/jsongenerator.go --llmin 10000 --llmax 30000 --ulmin 100 --ulmax 3000 -t=$(TOTAL) -p=0.7 -f="${SAMPLE_DATA_FOLDER}/${FILE_NAME}"
	@ echo "file ${SAMPLE_DATA_FOLDER}/${FILE_NAME} was generated." 


.PHONY: zookeeper
## zookeeper: starts zookeeper
zookeeper: 
	@ ${KAFKA_HOME}/bin/zookeeper-server-start.sh ${KAFKA_HOME}/config/zookeeper.properties

.PHONY: kafka
## kafka: starts kafka
kafka:
	@ ${KAFKA_HOME}/bin/kafka-server-start.sh ${KAFKA_HOME}/config/server.properties

.PHONY: producer
## producer: starts producer
producer:
	@ if [ -z "$(FILE_NAME)" ]; then echo >&2 please set file name via the variable FILE_NAME; exit 2; fi
	@ go run producer/producer.go -f=$(FILE_NAME)

.PHONY: consumer
## consumer: starts consumer
consumer:
	@ go run consumer/consumer.go

.PHONY: kafka-consumer-publish
## kafka-consumer-publish: Kafka's tool to read data from standard input and publish it to Kafka
kafka-consumer-publish:
	@ if [ -z "$(FILE_NAME)" ]; then echo >&2 please set file name via the variable FILE_NAME; exit 2; fi
	@ cat $(FILE_NAME) | ${KAFKA_HOME}/bin/kafka-console-producer.sh --topic "transactions" --bootstrap-server "localhost:9092"