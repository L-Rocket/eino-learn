/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package memory

import (
	"sync"

	"github.com/cloudwego/eino/schema"
)

type Memory struct {
	mu       sync.Mutex
	messages []*schema.Message
}

func NewMemory() *Memory {
	return &Memory{
		messages: make([]*schema.Message, 0),
	}
}

func (m *Memory) Add(msgs ...*schema.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, msgs...)
}

func (m *Memory) Get() []*schema.Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.messages
}
