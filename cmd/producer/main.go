package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GameXost/wbTestCase/models"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

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

func generateOrder() *models.Order {
	orderUID := fmt.Sprintf("order_%d_%d", time.Now().Unix(), rand.Intn(1000))

	return &models.Order{
		OrderUId:    orderUID,
		TrackNumber: fmt.Sprintf("TRACK_%d", rand.Intn(100000)),
		Entry:       "фыв",
		Delivery: models.Delivery{
			OrderUId: orderUID,
			Name:     "TestUser " + fmt.Sprint(rand.Intn(100)),
			Phone:    "фыв",
			Zip:      "фывфыв",
			City:     "окак",
			Address:  "Проспект Мира",
			Region:   "Вернадский",
			Email:    "собакагмэйлру",
		},
		Payment: models.Payment{
			OrderId:      orderUID,
			Transaction:  orderUID,
			RequestId:    "фыв",
			Currency:     "листики",
			Provider:     "кивиживи",
			Amount:       888,
			PaymentDt:    int64(time.Now().Unix()),
			Bank:         "банка",
			DeliveryCost: 99,
			GoodsTotal:   123123,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				OrderUId:    orderUID,
				ChrtId:      int64(rand.Intn(1000000) + 1),
				TrackNumber: fmt.Sprintf("TRACK_ITEM_%d", rand.Intn(10000)),
				Price:       132,
				RID:         fmt.Sprintf("rid_%d", rand.Intn(100000)),
				Name:        "Петруччо",
				Sale:        1,
				Size:        "ё хабло эспаньоло",
				TotalPrice:  111,
				NmId:        2342,
				Brand:       "Ё КОМО МАНЗАНАЗ",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "13",
		DeliveryService:   "ОКАК",
		Shardkey:          "9",
		SmId:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
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
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	producer, err := NewProducer(brokers, "orders")
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	defer producer.Close()

	log.Println("Producer started sending messages...")

	count := 100000
	for i := 1; i <= count; i++ {
		order := generateOrder()

		if err = producer.PublishOrder(ctx, order); err != nil {
			log.Printf("Failed to send order %s: %v", order.OrderUId, err)
			continue
		}

		log.Printf("%d Sent order: %s", i, order.OrderUId)

	}

	log.Println("Done!")
}
