// producer.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

func main() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	r := gin.Default()

	// 生产者
	r.POST("/send", func(c *gin.Context) {
		var json struct {
			Message string `json:"message"`
		}
		if err := c.BindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		topic := "test-topic" // Create a string variable for the topic name
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(json.Message),
		}, nil)

		c.JSON(http.StatusOK, gin.H{"status": "message sent"})
	})

	// 消费者
	r.GET("/consume", func(c *gin.Context) {
		consume, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": "localhost:9092",    // Kafka集群地址
			"group.id":          "my-consumer-group", // 自定义消费者组, 用于标识消费者, my-consumer-group为消费者组的名称
			"auto.offset.reset": "earliest",          // 从最早的消息开始消费
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
			os.Exit(1)
		}

		// Subscribe to the 'test-topic' topic
		err = consume.Subscribe("test-topic", nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to subscribe to topic: %s\n", err)
			os.Exit(1)
		}

		// Setup a channel for handling Ctrl+C interrupt
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

		// Wait for messages or interruption
		run := true
		for run == true {
			select {
			case <-sigchan:
				run = false
			default:
				ev := consume.Poll(100)
				if ev == nil {
					continue
				}

				switch e := ev.(type) {
				case *kafka.Message:
					fmt.Printf("Received message: %s on topic %s\n", string(e.Value), *e.TopicPartition.Topic)
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "Error: %v\n", e)
					run = e.Code() != kafka.ErrAllBrokersDown
				default:
					fmt.Printf("Ignored: %v\n", e)
				}
			}
		}

		// Close the consumer
		fmt.Println("Closing consumer")
		consume.Close()

		c.JSON(http.StatusOK, gin.H{"status": "message consumed"})
	})

	if err := r.Run(":8086"); err != nil {
		log.Fatalf("Failed to run server: %s", err)
	}
}
