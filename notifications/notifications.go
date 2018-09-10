package notifications

import (
	"log"
	"encoding/json"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/Shopify/sarama"
)

type kafkaPayload struct {
  Type string `json:"type"`
  Message string `json:"message"`
	Tag string `json:"tag"`

  encoded []byte
	err error
}

func NewNotificationsProducer(brokerList []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer: ", err)
	}

	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

	return producer
}

// PDAlert ...
func PDAlert(action string, serviceName string, integrationKey string, message string, tag string) error {
	event := pagerduty.Event{
		Type:        action,
		ServiceKey:  integrationKey,
		Description: message,
		IncidentKey: tag + serviceName,
	}
	resp, err := pagerduty.CreateEvent(event)
	if err != nil {
		log.Println(resp)
		return err
	}
	return nil
}

func KafkaAlert(action string, serviceName string, message string, tag string) error {
	event := sarama.ProducerMessage{
		Topic: "logging_nomad_events",
		Key: sarama.StringEncoder(tag),
		Value: sarama.StringEncoder(message),
	}

	select {
	case producer.Input() <- event:
		enqueued++
	}
}
