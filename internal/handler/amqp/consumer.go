package amqp

import (
	"context"
	"encoding/json"
	"log"
	"metertronik/internal/domain/entity"
	"metertronik/internal/service"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	svc *service.IngestService
	cfg *ConsumerConfig
}

type ConsumerConfig struct {
	QueueName     string
	RoutingKey    string
	Exchange      string
	PrefetchCount int
	RetryDelay    time.Duration
	LogInterval   time.Duration
}

func NewConsumer(svc *service.IngestService, cfg *ConsumerConfig) *Consumer {
	return &Consumer{
		svc: svc,
		cfg: cfg,
	}
}

func (c *Consumer) StartConsuming(ctx context.Context, connStr string) error {
	retryDelay := c.cfg.RetryDelay

	for {
		err := c.consumeWithReconnect(ctx, connStr)
		if err != nil {
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(retryDelay)
		}
	}
}

func (c *Consumer) consumeWithReconnect(ctx context.Context, connStr string) error {
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return err
	}
	defer conn.Close()

	notifyClose := make(chan *amqp.Error, 1)
	conn.NotifyClose(notifyClose)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	notifyChanClose := make(chan *amqp.Error, 1)
	ch.NotifyClose(notifyChanClose)

	var q amqp.Queue

	q, err = ch.QueueDeclare(
		c.cfg.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		if amqpErr, ok := err.(*amqp.Error); ok && amqpErr.Code == 406 {
			log.Printf("   Error Code: %d", amqpErr.Code)
			log.Printf("   Error Reason: %s", amqpErr.Reason)

			q, err = ch.QueueDeclare(
				c.cfg.QueueName,
				true,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				if amqpErr2, ok := err.(*amqp.Error); ok {
					log.Printf("   Error Code: %d", amqpErr2.Code)
					log.Printf("   Error Reason: %s", amqpErr2.Reason)
				}
				return err
			}
		} else {
			if amqpErr, ok := err.(*amqp.Error); ok {
				log.Printf("   Error Code: %d", amqpErr.Code)
				log.Printf("   Error Reason: %s", amqpErr.Reason)
				log.Printf("   Error Server: %v", amqpErr.Server)
			}
			return err
		}
	}

	err = ch.QueueBind(
		q.Name,
		c.cfg.RoutingKey,
		c.cfg.Exchange,
		false,
		nil,
	)
	if err != nil {
		if amqpErr, ok := err.(*amqp.Error); ok {
			log.Printf("   Error Code: %d", amqpErr.Code)
			log.Printf("   Error Reason: %s", amqpErr.Reason)
		}
		return err
	}

	err = ch.Qos(
		c.cfg.PrefetchCount,
		0,
		false,
	)
	if err != nil {
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
		if amqpErr, ok := err.(*amqp.Error); ok {
			log.Printf("   Error Code: %d", amqpErr.Code)
			log.Printf("   Error Reason: %s", amqpErr.Reason)
		}
		return err
	}

	done := make(chan error, 1)

	go processMessages(c, ctx, msgs, done)

	select {
	case err := <-notifyClose:
		if err != nil {
		}
		return err
	case err := <-notifyChanClose:
		if err != nil {
		}
		return err
	case err := <-done:
		if err != nil {
		} else {
		}
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func processMessages(c *Consumer, ctx context.Context, msgs <-chan amqp.Delivery, done chan<- error) {
	defer close(done)

	messageCount := 0
	lastMessageTime := time.Now()
	ticker := time.NewTicker(c.cfg.LogInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			done <- ctx.Err()
			return
		case <-ticker.C:
			elapsed := time.Since(lastMessageTime)
			log.Printf("Still waiting for messages... (Last message: %d, Elapsed since last: %v)", messageCount, elapsed)
		case d, ok := <-msgs:
			if !ok {
				log.Printf("Message channel closed. Processed total %d messages", messageCount)
				log.Printf("Channel closed without error - will reconnect")
				done <- nil
				return
			}

			messageCount++
			lastMessageTime = time.Now()
			log.Printf("Received message #%d", messageCount)
			log.Printf("Message delivery tag: %d, Exchange: %s, RoutingKey: %s", d.DeliveryTag, d.Exchange, d.RoutingKey)

			var data entity.RealTimeElectricity

			log.Printf("Unmarshaling message body (size: %d bytes)...", len(d.Body))
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				log.Printf("   Message body: %s", string(d.Body))
				continue
			}
			log.Printf("Message unmarshaled successfully. DeviceID: %s", data.DeviceID)

			log.Printf("Processing electricity data for device: %s...", data.DeviceID)
			if err := c.svc.ProcessRealTimeElectricity(ctx, &data); err != nil {
				log.Printf("Error processing electricity data: %v", err)
			}

			if messageCount == 20 {
				log.Printf("===== Reached message #20, continuing to wait for message #21... =====")
			}
		}
	}
}
