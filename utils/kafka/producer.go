/* !!
 * File: producer.go
 * File Created: Thursday, 5th March 2020 2:37:15 pm
 * Author: KimEricko™ (phamkim.pr@gmail.com)
 * -----
 * Last Modified: Thursday, 5th March 2020 2:37:16 pm
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
	"4vndevs/utils/json_debug"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
)

// ClientKafka type;
type ClientKafka struct {
	config            *ConfigKafka
	producer          sarama.SyncProducer
	kafkaClientConfig *sarama.Config
}

const (
	defaultNumPartitions     = int32(10)
	defaultReplicationFactor = int16(3)
)

var clientKafka *ClientKafka

// InstallKafkaClient func;
func InstallKafkaClient(config *ConfigKafka) {
	if config.NumPartitions == 0 {
		config.NumPartitions = defaultNumPartitions
	}
	if config.ReplicationFactor == 0 {
		config.ReplicationFactor = defaultReplicationFactor
	}

	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = 5
	conf.Producer.Return.Successes = true
	conf.Consumer.Return.Errors = true
	conf.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	version, err := sarama.ParseKafkaVersion("2.1.1")
	if err != nil {
		errDesc := fmt.Errorf("ClientKafka::InstallKafkaClient - Error parsing Kafka version: %+v", err)
		glog.V(1).Infoln(errDesc)
		return
	}
	conf.Version = version

	producer, err := sarama.NewSyncProducer(config.Addr, conf)
	if err != nil {
		glog.V(1).Infoln(err)
		return
	}

	clientKafka = &ClientKafka{
		kafkaClientConfig: conf,
		config:            config,
		producer:          producer,
	}
}

// GetKafkaClientInstance func;
func GetKafkaClientInstance() *ClientKafka {
	return clientKafka
}

// CreateTopic func;
func (c *ClientKafka) CreateTopic(topicName string) {
	if c == nil || c.config == nil {
		errDesc := fmt.Errorf("ClientKafka::CreateTopic - Need InstallKafkaClient first")
		glog.V(1).Infoln(errDesc)
		return
	}

	topicDetail := &sarama.TopicDetail{}
	topicDetail.NumPartitions = c.config.NumPartitions
	topicDetail.ReplicationFactor = c.config.ReplicationFactor
	topicDetail.ConfigEntries = make(map[string]*string)

	topicDetails := make(map[string]*sarama.TopicDetail)
	topicDetails[topicName] = topicDetail

	request := sarama.CreateTopicsRequest{
		Timeout:      time.Second * 15,
		TopicDetails: topicDetails,
	}

	// Send request to Broker
	broker := sarama.NewBroker(c.config.Addr[0])
	broker.Open(c.kafkaClientConfig)

	response, err := broker.CreateTopics(&request)

	// handle errors if any
	if err != nil {
		log.Printf("ClientKafka::CreateTopic - Error: %+v", err)
	}
	t := response.TopicErrors
	for key, val := range t {
		log.Printf("ClientKafka::CreateTopic - Key is %s", key)
		log.Printf("ClientKafka::CreateTopic - Value is %#v", val.Err.Error())
		if val.Err.Error() != "" {
		}
		log.Printf("ClientKafka::CreateTopic - Value3 is %#v", val.ErrMsg)
	}
	log.Printf("ClientKafka::CreateTopic - The response is %#v", response)
}

// ProducerPushMessage func;
func (c *ClientKafka) ProducerPushMessage(topic string, messageObj MessageKafka) (partition int32, offset int64, err error) {
	if c == nil || c.producer == nil {
		return 0, 0, fmt.Errorf("ClientKafka::ProducerPushMessage - Not found any producer")
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(json_debug.JSONDebugDataString(messageObj)),
	}

	return c.producer.SendMessage(msg)
}
