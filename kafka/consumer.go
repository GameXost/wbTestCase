package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GameXost/wbTestCase/internal/errHandle"
	"github.com/GameXost/wbTestCase/internal/service"
	"github.com/GameXost/wbTestCase/models"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"time"
)

type Consumer struct {
	client   *kgo.Client
	service  *service.Service
	dlqTopic string
}

func NewConsumer(brokers []string, topic, group string, srv *service.Service, dlqTopic string) (*Consumer, error) {
	options := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.DisableAutoCommit(),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()), // в бд защита от дублей, не пропустим необработанные
	}
	client, err := kgo.NewClient(options...)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:   client,
		service:  srv,
		dlqTopic: dlqTopic,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("kafka cons: context cancelled")
			return ctx.Err()
		default:
		}

		fetches := c.client.PollFetches(ctx)
		if err := fetches.Err(); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, kgo.ErrClientClosed) {

				return nil
			}
			log.Printf("kafka error: %v", err)
			time.Sleep(time.Second)
			continue
		}
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			for i := 0; i < 5; i++ {
				err := c.handleMessage(ctx, record)
				if err == nil {
					break
				}
				log.Printf("critical error handling message (topic: %s, partition: %d, offset: %d): %v",
					record.Topic, record.Partition, record.Offset, err)

				if ctx.Err() != nil {
					return ctx.Err()
				}
				if i == 4 {
					c.sendToDLQ(ctx, record)
					break
				}
				time.Sleep(3 * time.Second)
			}
		}

		if err := c.client.CommitUncommittedOffsets(ctx); err != nil {
			log.Printf("commit offset error: %v", err)
		}

	}

}

func (c *Consumer) handleMessage(ctx context.Context, record *kgo.Record) error {
	var order models.Order

	err := json.Unmarshal(record.Value, &order)
	if err != nil {
		log.Printf("invalid json error: %v", err)
		c.sendToDLQ(ctx, record)
		return nil
	}

	if order.OrderUId == "" {
		log.Println("order without UID")
		c.sendToDLQ(ctx, record)
		return nil
	}

	err = c.service.CreateOrder(ctx, &order)
	if err != nil {
		if errors.Is(err, errHandle.ErrValidation) {
			c.sendToDLQ(ctx, record)
			return nil
		}
		return err
	}
	return nil
}

func (c *Consumer) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

func (c *Consumer) sendToDLQ(ctx context.Context, record *kgo.Record) {
	dlqRec := &kgo.Record{
		Topic: c.dlqTopic,
		Key:   record.Key,
		Value: record.Value,
	}
	if err := c.client.ProduceSync(ctx, dlqRec).FirstErr(); err != nil {
		log.Printf("fauked to send to dlq: %v", err)
	}
}
