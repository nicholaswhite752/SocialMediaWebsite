package controllers

import (
	"SocialWebsite/config"
	"SocialWebsite/models"
	"encoding/json"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
)

//Sends user and page visited to Kafka Topic
func ProduceUserPageRequest(usernameToSend string, pageToSend string) {

	// Delivery report handler for produced messages
	// This section is for delivery reports

	//This should work
	//go func() {
	//	for e := range config.Producer.Events() {
	//		switch ev := e.(type) {
	//		case *kafka.Message:
	//			if ev.TopicPartition.Error != nil {
	//				//Prints to file if delivery failed
	//				_, err := fmt.Fprintf(config.FileForKafka, "Delivery failed: %v\n", ev.TopicPartition)
	//				if err != nil {
	//					log.Println(err, "Failed to write Kafka Delivery Failed to file")
	//				}
	//			} else {
	//				//Prints to file if delivery success
	//				_, err := fmt.Fprintf(config.FileForKafka, "Delivered message to %v\n", ev.TopicPartition)
	//				if err != nil {
	//					log.Println(err, "Failed to write Kafka Delivery Success to file")
	//				}
	//			}
	//		}
	//	}
	//}()

	// Produce messages to topic (asynchronously)
	topic := "userPageVisits"

	//Creates a JSON of username and pagevisited
	userPageVisitsJson := models.UserPageVisits{
		Username:    usernameToSend,
		PageVisited: pageToSend,
	}

	jsonMessage, err := json.Marshal(userPageVisitsJson)
	if err != nil {
		log.Println(err, "Failed to create JSON for User Page Visit : kafkaProducerHandler.go")
		return
	}

	err = config.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(usernameToSend),
		Value:          jsonMessage,
	}, nil)
	if err != nil {
		log.Println(err, "Failed to Produce Message to Kafka : kafkaProducerHandler.go")
		return
	}

}
