package config

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"os"
)

var Producer *kafka.Producer
//This var is for delivery reports in a file
var FileForKafka *os.File

func init(){
	//Connects Producer To Kafka
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "<server IP>:9092"})
	if err != nil {
		log.Fatal(err, "Kafka Producer Did Not Start")
	}

	Producer = p

	//Creates or Opens a text file called kafkaDelivery.txt
	//This section is for delivery reports in a file
	f, err := os.OpenFile("kafkaDelivery.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err, "File for Kafka Delivery did not open")
	}

	FileForKafka = f

}
