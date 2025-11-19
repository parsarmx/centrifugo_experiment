package broker

import (
	"context"
	"fmt"
	"golang_template/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Rabbit interface {
	Close() error
	Channel() *amqp.Channel

	// might delete later
	ConsumeMessage(ctx context.Context) error
	HealthCheck(ctx context.Context) error
}

type rabbit struct {
	cfg     *config.RabbitConfig
	logger  *zap.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbit(cfg *config.RabbitConfig, logger *zap.Logger) (Rabbit, error) {
	conn, err := amqp.Dial(config.GetRabbitMQUrl(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open up a channel")
	}

	return &rabbit{cfg: cfg, logger: logger, channel: ch, conn: conn}, nil
}

func (r *rabbit) Channel() *amqp.Channel {
	return r.channel
}

func (c *rabbit) Close() error {
	var errs []error

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close RabbitMQ channel: %w", err))
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close RabbitMQ connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("closing rabbitmq: %v", errs)
	}
	return nil
}

func (r *rabbit) HealthCheck(ctx context.Context) error {
	ch := r.Channel()

	q, err := ch.QueueDeclare(
		"health-check", // name
		false,          // durable
		false,          // delete when unused
		true,           // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue")
	}

	body := "health-check"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message")
	}

	r.logger.Info("health-check sent to RabbitMQ")

	return nil
}

func (r *rabbit) ConsumeMessage(ctx context.Context) error {
	ch := r.Channel()

	q, err := ch.QueueDeclare(
		"main", // name
		false,  // durable
		false,  // delete when unused
		true,   // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			r.logger.Info(fmt.Sprintf("Received a message: %s", d.Body))
		}
	}()

	r.logger.Info("Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}
