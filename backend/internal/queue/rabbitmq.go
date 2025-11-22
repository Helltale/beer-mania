package queue

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Helltale/beer-mania/backend/internal/config"

	"github.com/google/uuid"
)

const (
	ExchangeName = "image_processing_exchange"
	QueueName    = "image_processing"
	DLQName      = "image_processing.dlq"
	RoutingKey   = "image.processing"
)

type RabbitMQQueue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *config.RabbitMQConfig
	logger  *slog.Logger
}

func NewRabbitMQQueue(cfg *config.RabbitMQConfig) (*RabbitMQQueue, error) {
	return NewRabbitMQQueueWithLogger(cfg, slog.Default())
}

func NewRabbitMQQueueWithLogger(cfg *config.RabbitMQConfig, logger *slog.Logger) (*RabbitMQQueue, error) {
	hostPort := net.JoinHostPort(cfg.Host, cfg.Port)
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		hostPort,
		cfg.VHost)

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		if closeErr := conn.Close(); closeErr != nil {
			return nil, fmt.Errorf("failed to open channel: %w (also failed to close connection: %w)", err, closeErr)
		}
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	queue := &RabbitMQQueue{
		conn:    conn,
		channel: channel,
		cfg:     cfg,
		logger:  logger,
	}

	if setupErr := queue.setup(); setupErr != nil {
		var errs []error
		errs = append(errs, setupErr)
		if closeErr := channel.Close(); closeErr != nil {
			errs = append(errs, fmt.Errorf("failed to close channel: %w", closeErr))
		}
		if closeErr := conn.Close(); closeErr != nil {
			errs = append(errs, fmt.Errorf("failed to close connection: %w", closeErr))
		}
		return nil, fmt.Errorf("failed to setup RabbitMQ: %w", errs[0])
	}

	return queue, nil
}

func (q *RabbitMQQueue) setup() error {
	err := q.channel.ExchangeDeclare(
		ExchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	_, err = q.channel.QueueDeclare(
		DLQName, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	_, err = q.channel.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    "", // Use default exchange for DLQ
			"x-dead-letter-routing-key": DLQName,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	err = q.channel.QueueBind(
		QueueName,    // queue name
		RoutingKey,   // routing key
		ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	q.logger.Info("RabbitMQ setup completed",
		"exchange", ExchangeName,
		"queue", QueueName,
		"dlq", DLQName)

	return nil
}

func (q *RabbitMQQueue) PublishTask(ctx context.Context, taskID uuid.UUID, imageID uuid.UUID) error {
	msg := &ProcessingMessage{
		TaskID:  taskID,
		ImageID: imageID,
	}

	body, err := msg.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = q.channel.PublishWithContext(
		ctx,
		ExchangeName, // exchange
		RoutingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	q.logger.InfoContext(ctx, "Published task",
		"task_id", taskID,
		"image_id", imageID)
	return nil
}

func (q *RabbitMQQueue) ConsumeTasks(
	ctx context.Context,
	handler func(taskID uuid.UUID, imageID uuid.UUID) error,
) error {
	if err := q.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := q.channel.Consume(
		QueueName, // queue
		"",        // consumer tag (empty = auto-generate)
		false,     // auto-ack (false = manual ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	q.logger.InfoContext(ctx, "Started consuming from queue", "queue", QueueName)

	go q.processMessages(ctx, msgs, handler)

	return nil
}

func (q *RabbitMQQueue) processMessages(
	ctx context.Context,
	msgs <-chan amqp.Delivery,
	handler func(taskID uuid.UUID, imageID uuid.UUID) error,
) {
	for {
		select {
		case <-ctx.Done():
			q.logger.InfoContext(ctx, "Stopping consumer", "error", ctx.Err())
			return
		case msg, ok := <-msgs:
			if !ok {
				q.logger.InfoContext(ctx, "Message channel closed")
				return
			}

			q.handleMessage(msg, handler)
		}
	}
}

func (q *RabbitMQQueue) handleMessage(
	msg amqp.Delivery,
	handler func(taskID uuid.UUID, imageID uuid.UUID) error,
) {
	processingMsg, unmarshalErr := UnmarshalProcessingMessage(msg.Body)
	if unmarshalErr != nil {
		q.logger.Warn("Failed to unmarshal message, sending to DLQ", "error", unmarshalErr)
		if nackErr := msg.Nack(false, false); nackErr != nil {
			q.logger.Error("Failed to nack invalid message", "error", nackErr)
		}
		return
	}

	if handlerErr := handler(processingMsg.TaskID, processingMsg.ImageID); handlerErr != nil {
		q.logger.Warn("Task processing failed, sending to DLQ",
			"task_id", processingMsg.TaskID,
			"error", handlerErr)
		if nackErr := msg.Nack(false, false); nackErr != nil {
			q.logger.Error("Failed to nack failed message", "error", nackErr)
		}
		return
	}

	if ackErr := msg.Ack(false); ackErr != nil {
		q.logger.Error("Failed to acknowledge message", "error", ackErr)
	} else {
		q.logger.Info("Task processed successfully", "task_id", processingMsg.TaskID)
	}
}

func (q *RabbitMQQueue) Close() error {
	var errs []error

	if q.channel != nil {
		if err := q.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close channel: %w", err))
		}
	}

	if q.conn != nil {
		if err := q.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing RabbitMQ: %v", errs)
	}

	return nil
}
