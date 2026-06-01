package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/magnifimind/sentinel/internal/model"
	"github.com/magnifimind/sentinel/internal/store"
)

type Consumer struct {
	readers []*kafkago.Reader
	store   *store.Store
	logger  *slog.Logger
}

func NewConsumer(brokers, topics []string, group string, s *store.Store, logger *slog.Logger) *Consumer {
	var readers []*kafkago.Reader
	for _, topic := range topics {
		r := kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        group,
			MinBytes:       1e3,  // 1KB
			MaxBytes:       10e6, // 10MB
			CommitInterval: time.Second,
			StartOffset:    kafkago.LastOffset,
		})
		readers = append(readers, r)
	}
	return &Consumer{readers: readers, store: s, logger: logger}
}

// Run starts consuming from all topics. Blocks until ctx is cancelled.
func (c *Consumer) Run(ctx context.Context) {
	for _, r := range c.readers {
		go c.consume(ctx, r)
	}
	<-ctx.Done()
	for _, r := range c.readers {
		r.Close()
	}
}

func (c *Consumer) consume(ctx context.Context, r *kafkago.Reader) {
	topic := r.Config().Topic
	c.logger.Info("kafka consumer started", "topic", topic)

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.logger.Error("kafka read error", "topic", topic, "error", err)
			time.Sleep(time.Second)
			continue
		}

		var event model.KafkaEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.logger.Error("kafka unmarshal error", "topic", topic, "error", err, "offset", msg.Offset)
			continue
		}
		event.Topic = topic

		if err := c.handleEvent(ctx, event); err != nil {
			c.logger.Error("event handling failed", "topic", topic, "type", event.Type, "error", err)
		}
	}
}

func (c *Consumer) handleEvent(ctx context.Context, event model.KafkaEvent) error {
	switch event.Topic {
	case "agent.decisions":
		return c.handleDecision(ctx, event)
	case "agent.actions", "agent.outcomes":
		return c.handleActionOrOutcome(ctx, event)
	default:
		c.logger.Warn("unknown topic", "topic", event.Topic)
		return nil
	}
}

func (c *Consumer) handleDecision(ctx context.Context, event model.KafkaEvent) error {
	var d model.AgentDecision
	if err := json.Unmarshal(event.Payload, &d); err != nil {
		return err
	}
	if d.Timestamp.IsZero() {
		d.Timestamp = event.Timestamp
	}
	if d.AgentID == "" {
		d.AgentID = event.AgentID
	}
	return c.store.InsertDecision(ctx, &d)
}

func (c *Consumer) handleActionOrOutcome(ctx context.Context, event model.KafkaEvent) error {
	d := model.AgentDecision{
		Timestamp:    event.Timestamp,
		AgentID:      event.AgentID,
		DecisionType: event.Type,
		Context:      event.Payload,
		Action:       event.Type,
		PolicyResult: "recorded",
		Success:      true,
	}
	return c.store.InsertDecision(ctx, &d)
}
