package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"time"
)

var Address = []string{":9092"}

func main() {
	syncProducer(Address)
	//asyncProducer1(Address)
}

//同步消息模式
func syncProducer(address []string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	defer p.Close()

	c, err := sarama.NewConsumer(address, config)
	if err != nil {
		log.Println("error,", err)
	}
	defer c.Close()


	topic := "my_topic"
	srcValue := "sync: this is a message. index=%d"
	//
	//go func() {
	//	partitionConsumer, err := c.ConsumePartition(topic, 0, 0)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	defer func() {
	//		if err := partitionConsumer.Close(); err != nil {
	//			log.Fatalln(err)
	//		}
	//	}()
	//
	//	// Trap SIGINT to trigger a shutdown.
	//	signals := make(chan os.Signal, 1)
	//	signal.Notify(signals, os.Interrupt)
	//
	//	consumed := 0
	//ConsumerLoop:
	//	for {
	//		select {
	//		case msg := <-partitionConsumer.Messages():
	//			log.Printf("Consumed message offset %d\n", msg.Offset)
	//			consumed++
	//		case <-signals:
	//			break ConsumerLoop
	//		}
	//	}
	//
	//	log.Printf("Consumed: %d\n", consumed)
	//}()

	for i := 0; i < 10; i++ {
		value := fmt.Sprintf(srcValue, i)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(value),
		}
		part, offset, err := p.SendMessage(msg)
		if err != nil {
			log.Printf("send message(%s) err=%s \n", value, err)
		} else {
			fmt.Fprintf(os.Stdout, value+"发送成功，partition=%d, offset=%d \n", part, offset)
		}
		time.Sleep(2 * time.Second)

		//r, err := c.ConsumePartition(topic, 0, offset)
		//fmt.Println("------->", r, err)

		//topics, err := c.Topics()
		//fmt.Println(topics, err)
	}
}
