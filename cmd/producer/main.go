package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GameXost/wbTestCase/internal/generator"
	"github.com/GameXost/wbTestCase/models"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var count = 10000

type Producer struct {
	client *kgo.Client
	topic  string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	)
	if err != nil {
		return nil, err
	}
	return &Producer{
		client: client,
		topic:  topic,
	}, nil
}

func (p *Producer) Close() {
	p.client.Close()
}

func (p *Producer) PublishOrder(ctx context.Context, order *models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Topic: p.topic,
		Key:   []byte(order.OrderUId),
		Value: data,
	}

	results := p.client.ProduceSync(ctx, record)
	if err = results.FirstErr(); err != nil {
		return fmt.Errorf("kafka produce error: %w", err)
	}
	return nil
}

func main() {

	ctx := context.Background()
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	producer, err := NewProducer(brokers, "orders")
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	defer producer.Close()

	log.Println("Producer started sending messages...")

	for i := 1; i <= count; i++ {
		orderUID := strconv.Itoa(rand.Int())
		var order *models.Order
		if rand.Intn(10) < 8 {
			order = generator.ValidOrder(orderUID)
		} else {
			order = generator.InvalidOrder(orderUID)
		}
		if err = producer.PublishOrder(ctx, order); err != nil {
			log.Printf("Failed to send order %s: %v", order.OrderUId, err)
			continue
		}
		log.Printf("order: %s", order.OrderUId)

	}

	log.Println("Donon")
}
