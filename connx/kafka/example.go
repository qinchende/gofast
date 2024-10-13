package kafka

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

// Kafka集群的broker地址列表
var brokers = []string{"IP:9001", "IP:9002", "IP:9003"}
var topicName = "tpc-sdxtst1"
var groupName = " Group-Sdx1"

type cgHandler struct{}

func (cgHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (cgHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (cgHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		doMessage(msg)
		session.MarkMessage(msg, "")
	}
	return nil
}

func doMessage(msg *sarama.ConsumerMessage) {
	fmt.Printf("ts: %v, topic(%s)/part(%d)/off(%d); value: %s\n",
		msg.Timestamp, msg.Topic, msg.Partition, msg.Offset, msg.Value)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 消费者组
func ConsumerGroup(gpName string) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_6_0_0
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(brokers, gpName, cfg)
	if err != nil {
		log.Fatalf("Failed to NewConsumerGroup: ", err)
	}
	defer func() { _ = group.Close() }()

	// 监听可能的错误并打印
	go func() {
		// group.Close() 执行的时候，error通道会关闭，本gor自动退出，不会泄露
		for err = range group.Errors() {
			fmt.Println("Err: ", err)
		}
	}()

	fmt.Println("Consumer start")
	for {
		if err = group.Consume(context.Background(), []string{topicName}, cgHandler{}); err != nil {
			fmt.Println("group.Consume err: ", err)
			return
		}
	}
}

// 消费所有分区
func ConsumerAllPartition() {
	// 用默认配置创建一个消费者
	// cfg := sarama.NewConfig()
	consumer, err1 := sarama.NewConsumer(brokers, nil)
	if err1 != nil {
		log.Fatalf("Failed to create consumer: %v", err1)
	}
	// 查询有几个分区
	partis, err2 := consumer.Partitions(topicName)
	if err2 != nil {
		log.Fatalf("Failed to get list of partitions: %v", err2)
	}
	fmt.Println("parts: ", partis)

	ctx, cancel := context.WithCancel(context.Background())
	// 遍历每个分区，对每个分区开启异步消费
	for _, partNum := range partis {
		cpItem, err3 := consumer.ConsumePartition(topicName, partNum, sarama.OffsetNewest)
		defer cpItem.AsyncClose()
		if err3 != nil {
			log.Fatalf("Failed to consumer for partition %d: %v", partNum, err3)
			return
		}

		go func(ctx context.Context, part int32, pc sarama.PartitionConsumer) {
			count := 0
		msgLoop:
			for {
				select {
				case msg := <-pc.Messages():
					doMessage(msg)
					count++
				case <-ctx.Done():
					break msgLoop
				}
			}
			log.Printf("Part %d, Consumed: %d\n", part, count)
		}(ctx, partNum, cpItem)
	}

	// 应用关闭时候，退出所有资源
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals               // 阻塞直到应用退出
	cancel()                // 关闭所有协程
	time.Sleep(time.Second) // 等待打印协程退出日志，这不是必须的
	log.Println("ConsumerAllPartition exit")
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 同步生产者
func SyncProducer() {
	cfg := sarama.NewConfig()                              // 创建一个配置对象
	cfg.Producer.RequiredAcks = sarama.WaitForLocal        // Leader成功后确认
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner // 随机分区器
	cfg.Producer.Return.Successes = true                   // 消息成功发送时返回

	// 使用broker地址和配置创建一个同步Producer
	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("Failed to NewSyncProducer: %v", err)
	}
	defer func() {
		if err = producer.Close(); err != nil {
			log.Fatalf("Failed to close producer: %v", err)
		}
	}()

	// 要发送的消息
	message := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder("Message payload here"),
	}
	// 发送消息，返回所属分区和偏移量
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	// 打印消息发送详情
	fmt.Printf("Msg saved in topic(%s)/part(%d)/off(%d)\n", message.Topic, partition, offset)
}

// 异步生产者Goroutines
func AsyncProducer() {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("Failed to NewAsyncProducer: %v", err)
		return
	}
	defer func() {
		if err = producer.Close(); err != nil {
			log.Fatalf("Failed to close producer: %v", err)
		}
	}()

	// 监听捕获关闭信号，方便正常退出
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var wg sync.WaitGroup
	var msgCount, sucCount, errCount int

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			sucCount++
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err = range producer.Errors() {
			log.Println(err)
			errCount++
		}
	}()

	msg := &sarama.ProducerMessage{Topic: topicName}
producerLoop:
	for {
		msg.Value = sarama.StringEncoder(fmt.Sprintf("the-msg-%d", msgCount))

		select {
		case producer.Input() <- msg:
			msgCount++
			log.Println(msg.Value)
			time.Sleep(1 * time.Second)
		case <-signals:
			producer.AsyncClose()
			log.Println("os.Interrupt, AsyncClose, then exit")
			break producerLoop
		}
	}

	// 直到退出，打印统计数据
	wg.Wait()
	log.Printf("Msg-Count: %d; Suc-Count: %d; Err-Count: %d\n", msgCount, sucCount, errCount)
}

// 异步生产者Select
func AsyncProducerSelect() {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("Failed to NewAsyncProducer: %v", err)
		return
	}
	defer func() {
		if err = producer.Close(); err != nil {
			log.Fatalf("Failed to close producer: %v", err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var msgCount, sucCount, errCount int
	msg := &sarama.ProducerMessage{Topic: topicName}

producerLoop:
	for {
		msg.Value = sarama.StringEncoder(fmt.Sprintf("the-msg-%d", msgCount))

		select {
		case producer.Input() <- msg:
			msgCount++
			sucCount++
			log.Println(msg.Value)
			time.Sleep(1 * time.Second)
		case err = <-producer.Errors():
			log.Println("Failed to produce message", err)
			errCount++
			sucCount--
		case <-signals:
			log.Println("os.Interrupt, AsyncClose, then exit")
			break producerLoop
		}
	}

	log.Printf("Msg-Count: %d; Suc-Count: %d; Err-Count: %d\n", msgCount, sucCount, errCount)
}

func RunExample() {
	//ConsumerAllPartition()
	//ConsumerGroup("sdx")

	//SyncProducer()
	//AsyncProducer()
	AsyncProducerSelect()

	log.Println("RunExample exit, bye bye ...")
}
