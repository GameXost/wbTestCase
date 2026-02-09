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
	client  *kgo.Client
	service *service.Service
}

func NewConsumer(brokers []string, topic, group string, srv *service.Service) (*Consumer, error) {
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
		client:  client,
		service: srv,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	defer c.client.Close()
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
			for {
				err := c.handleMessage(ctx, record)
				if err == nil {
					break
				}
				log.Printf("critical error handling message (topic: %s, partition: %d, offset: %d): %v",
					record.Topic, record.Partition, record.Offset, err)

				if ctx.Err() != nil {
					return ctx.Err()
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
		return nil
	}

	if order.OrderUId == "" {
		log.Println("order without UID")
		return nil
	}

	err = c.service.CreateOrder(ctx, &order)
	if err != nil {
		if errors.Is(err, errHandle.ErrValidation) {
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
