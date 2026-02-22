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

package flow

import (
	"context"
	"log"

	"GoClawd/demo3/pkg/memory"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

func Generate(ctx context.Context, llm model.ToolCallingChatModel, mem *memory.Memory, in []*schema.Message) *schema.Message {
	// 1. Get history from memory
	history := mem.Get()

	// 2. Combine history and new messages
	fullMessages := append(history, in...)

	// 3. Generate response
	result, err := llm.Generate(ctx, fullMessages)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}

	// 4. Save to memory
	mem.Add(in...)
	mem.Add(result)

	return result
}

func Stream(ctx context.Context, llm model.ToolCallingChatModel, mem *memory.Memory, in []*schema.Message) *schema.StreamReader[*schema.Message] {
	// 1. Get history from memory
	history := mem.Get()

	// 2. Combine history and new messages
	fullMessages := append(history, in...)

	// 3. Stream response
	result, err := llm.Stream(ctx, fullMessages)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}

	return result
}
