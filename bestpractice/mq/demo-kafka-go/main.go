package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
	"time"
)

func main() {
	// 生产消息
	ProduceMsg()

	// 消费消息
	ConsumeMsg()

	// By default kafka has the auto.create.topics.enable='true' (KAFKA_AUTO_CREATE_TOPICS_ENABLE='true' in the wurstmeister/kafka kafka docker image)
	// 默认情况下，自动创建topic是开启的，此时，可以使用下面的方式去自动创建和使用topic
	CreateTopic()

	// 如果无法自动创建topic，此时，可以使用下面的方式去创建topic
	CreateTopic2()

	// 等等。。。
}

func CreateTopic2() {
	// to create topics when auto.create.topics.enable='false'
	topic := "my-topic"

	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}

func CreateTopic() {
	// to create topics when auto.create.topics.enable='true' ，这个true是默认值
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "my-topic", 0)
	if err != nil {
		panic(err.Error())
	}

	conn.Close()

}

func ConsumeMsg() {
	// to consume messages
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}

func ProduceMsg() {
	// to produce messages
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	// 针对此connection的写超时，全局的。不设置，则默认不会有写超时。
	if err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		log.Fatal("setWriteDeadline error:", err)
	}

	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err = conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
