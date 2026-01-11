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
		_ = c.consumeWithReconnect(ctx, connStr)

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
			q, err = ch.QueueDeclare(
				c.cfg.QueueName,
				true,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				return err
			}
		} else {
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
		return err
	}

	done := make(chan error, 1)

	go processMessages(c, ctx, msgs, done)

	select {
	case err := <-notifyClose:
		return err
	case err := <-notifyChanClose:
		return err
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func processMessages(c *Consumer, ctx context.Context, msgs <-chan amqp.Delivery, done chan<- error) {
	defer close(done)

	messageCount := 0

	for {
		select {
		case <-ctx.Done():
			done <- ctx.Err()
			return
		case d, ok := <-msgs:
			if !ok {
				done <- nil
				return
			}

			messageCount++

			log.Printf("\n\n[%s]\nReceiving Data...\nChecking Data...", time.Now().Format("2006-01-02 15:04:05"))

			var data entity.RealTimeElectricity
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("invalid payload, skip")
				continue
			}

			if err := c.svc.ProcessRealTimeElectricity(ctx, &data); err != nil {
				log.Printf("failed processing data")
				continue
			}
		}
	}
}
