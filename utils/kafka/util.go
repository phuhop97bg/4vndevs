/* !!
 * File: utils.go
 * File Created: Thursday, 5th March 2020 2:37:35 pm
 * Author: KimEricko™ (phamkim.pr@gmail.com)
 * -----
 * Last Modified: Thursday, 5th March 2020 2:37:35 pm
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

// ConfigKafka type;
type ConfigKafka struct {
	Addr              []string `json:"addr"`
	ReplicationFactor int16    `json:"replicationFactor"`
	NumPartitions     int32    `json:"numPartitions"`
}

// MessageKafka func;
type MessageKafka struct {
	Event      string `json:"event"`
	ObjectJSON string `json:"object_json"`
}
