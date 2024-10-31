package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	brokers := []string{"kafka:29092"}
	topic := "test"

	client, err := sarama.NewConsumerGroup(brokers, "test-group", config)
	if err != nil {
		log.Fatalf("unable to create kafka consumer group: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			err := client.Consume(ctx, []string{topic}, &consumerHandler{})
			if err != nil {
				log.Printf("consume error: %v", err)
			}

			select {
			case <-signals:
				cancel()
				return
			default:
			}
		}
	}()

	wg.Wait()
}

type consumerHandler struct{}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Received message: key=%s, value=%s, partition=%d, offset=%d\n", string(msg.Key), string(msg.Value), msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
	}
	return nil
}
