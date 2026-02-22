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
	"fmt"
	"log"

	"GoClawd/demo3/pkg/memory"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func GenerateWithTools(ctx context.Context, llm model.ToolCallingChatModel, mem *memory.Memory, in []*schema.Message, tools []tool.BaseTool) *schema.Message {
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		log.Fatalf("create tools node failed: %v", err)
	}

	history := mem.Get()
	fullMessages := append(history, in...)

	var result *schema.Message
	var currentMessages []*schema.Message = fullMessages

	maxIterations := 5
	for i := 0; i < maxIterations; i++ {
		resp, err := llm.Generate(ctx, currentMessages)
		if err != nil {
			log.Fatalf("llm generate failed: %v", err)
		}

		if len(resp.ToolCalls) == 0 {
			result = resp
			mem.Add(currentMessages...)
			mem.Add(resp)
			break
		}

		fmt.Printf("\n=== Tool Call %d ===\n", i+1)
		for _, tc := range resp.ToolCalls {
			fmt.Printf("Tool: %s, Arguments: %s\n", tc.Function.Name, tc.Function.Arguments)
		}

		toolResults, err := toolsNode.Invoke(ctx, resp)
		if err != nil {
			log.Fatalf("tools node invoke failed: %v", err)
		}

		currentMessages = append(currentMessages, resp)
		currentMessages = append(currentMessages, toolResults...)

		fmt.Printf("Tool Results: %d messages\n", len(toolResults))
	}

	return result
}

func StreamWithTools(ctx context.Context, llm model.ToolCallingChatModel, mem *memory.Memory, in []*schema.Message, tools []tool.BaseTool) *schema.StreamReader[*schema.Message] {
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		log.Fatalf("create tools node failed: %v", err)
	}

	history := mem.Get()
	fullMessages := append(history, in...)

	var currentMessages []*schema.Message = fullMessages
	maxIterations := 5

	resultChan, writer := schema.Pipe[*schema.Message](0)
	streamReader := resultChan

	go func() {
		defer writer.Close()

		for i := 0; i < maxIterations; i++ {
			resp, err := llm.Generate(ctx, currentMessages)
			if err != nil {
				log.Printf("llm generate failed: %v", err)
				return
			}

			if len(resp.ToolCalls) == 0 {
				writer.Send(resp, nil)
				mem.Add(currentMessages...)
				mem.Add(resp)
				break
			}

			fmt.Printf("\n=== Tool Call %d ===\n", i+1)
			for _, tc := range resp.ToolCalls {
				fmt.Printf("Tool: %s, Arguments: %s\n", tc.Function.Name, tc.Function.Arguments)
			}

			toolResults, err := toolsNode.Invoke(ctx, resp)
			if err != nil {
				log.Printf("tools node invoke failed: %v", err)
				return
			}

			currentMessages = append(currentMessages, resp)
			currentMessages = append(currentMessages, toolResults...)

			fmt.Printf("Tool Results: %d messages\n", len(toolResults))
		}
	}()

	return streamReader
}
