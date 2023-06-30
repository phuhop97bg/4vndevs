/* !!
 * File: consumer.go
 * File Created: Thursday, 5th March 2020 2:37:24 pm
 * Author: KimEricko™ (phamkim.pr@gmail.com)
 * -----
 * Last Modified: Thursday, 5th March 2020 2:37:24 pm
 * Modified By: KimEricko™ (phamkim.pr@gmail.com>)
 * -----
 * Copyright 2018 - 2020 mySoha Platform, VCCorp
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Developer: NhokCrazy199 (phamkim.pr@gmail.com)
 */

package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/golang/glog"

	"github.com/Shopify/sarama"
)

// ConsumerProcessInstance func;
type ConsumerProcessInstance interface {
	ErrorCallback(err error)
	MessageCallback(messageObj MessageKafka)
}

// InstallConsumerGroup func;
// Func need run as goroutine
func (c *ClientKafka) InstallConsumerGroup(group string, topics []string, processingInstance ConsumerProcessInstance) {
	if c == nil || c.config == nil {
		errDesc := fmt.Errorf("ClientKafka::InstallConsumerGroup - Need InstallKafkaClient first")
		glog.V(1).Infoln(errDesc)

		return
	}

	c.kafkaClientConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	c.kafkaClientConfig.Consumer.Return.Errors = true

	// config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	c.kafkaClientConfig.Consumer.Group.Session.Timeout = 20 * time.Second
	c.kafkaClientConfig.Consumer.Group.Heartbeat.Interval = 6 * time.Second
	c.kafkaClientConfig.Consumer.MaxProcessingTime = 500 * time.Millisecond
	// config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange

	c.kafkaClientConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumer := newConsumerInstance(processingInstance)

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(c.config.Addr, group, c.kafkaClientConfig)
	if err != nil {
		glog.V(1).Infof("ClientKafka::InstallConsumerGroup - Error creating consumer group client: %+v", err)
		cancel()
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, topics, consumer); err != nil {
				processingInstance.ErrorCallback(err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("ClientKafka::InstallConsumerGroup - Consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("ClientKafka::InstallConsumerGroup - Terminating: context cancelled")
	case <-sigterm:
		log.Println("ClientKafka::InstallConsumerGroup - Terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		glog.V(1).Infof("ClientKafka::InstallConsumerGroup - Error closing client: %v", err)
	}
}

// Consumer represents a Sarama consumer group consumer
type consumerInstance struct {
	ready              chan bool
	processingInstance ConsumerProcessInstance
}

// NewConsumer func;
func newConsumerInstance(processingInstance ConsumerProcessInstance) *consumerInstance {
	// Mark the consumer as ready
	return &consumerInstance{
		ready:              make(chan bool),
		processingInstance: processingInstance,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *consumerInstance) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *consumerInstance) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *consumerInstance) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		// glog.V(3).Infof("ClientKafka::ConsumeClaim - Message claimed: value = %s, timestamp = %v, topic = %s, part: %d, offset: %d", string(message.Value), message.Timestamp, message.Topic, message.Partition, message.Offset)

		msgObj := MessageKafka{}
		err := json.Unmarshal(message.Value, &msgObj)
		if err != nil {
			consumer.processingInstance.ErrorCallback(err)
		} else {
			consumer.processingInstance.MessageCallback(msgObj)
		}

		session.MarkMessage(message, "")
	}

	return nil
}
