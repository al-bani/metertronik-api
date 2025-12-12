package amqp

import (
	"context"
	"encoding/json"
	"log"
	"metertronik/internal/domain/entity"
	"metertronik/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	svc *service.IngestService
}

func NewConsumer(svc *service.IngestService) *Consumer {
	return &Consumer{svc: svc}
}

func (c *Consumer) StartConsuming(ctx context.Context, connStr string) error {
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"electricity_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,
		"electricity_metrics",
		"amq.topic",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Printf("✅ Consumer started, waiting for messages...")

	forever := make(chan struct{})

	go processMessages(c, ctx, msgs)
	<-forever

	return nil
}

func processMessages(c *Consumer, ctx context.Context, msgs <-chan amqp.Delivery) {
	log.Println("Process messages started")

	for d := range msgs {

		var data entity.RealTimeElectricity

		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			log.Printf("❌ Failed to unmarshal message: %v", err)
			log.Printf("   Message body: %s", string(d.Body))
			continue
		}

		if err := c.svc.ProcessRealTimeElectricity(ctx, &data); err != nil {
			log.Printf("❌ Error processing electricity data: %v", err)
		} else {
			log.Printf("✅ Successfully processed electricity data for device: %s", data.DeviceID)
		}
	}
}
