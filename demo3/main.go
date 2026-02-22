/*
 * Copyright 2024 CloudWeGo Authors
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

package main

import (
	"context"
	"log"
	"sync"

	"GoClawd/demo3/pkg/chatmodel"
	"GoClawd/demo3/pkg/flow"
	"GoClawd/demo3/pkg/memory"
	"GoClawd/demo3/pkg/prompt"
	"GoClawd/demo3/pkg/tools"

	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	wg := sync.WaitGroup{}

	// 创建 memory
	log.Printf("===create memory===\n")
	mem := memory.NewMemory()

	// 使用模版创建messages
	log.Printf("===create messages===\n")
	messages := prompt.CreateMessagesFromTemplate()

	// 创建llm
	log.Printf("===create llm===\n")
	cm := chatmodel.CreateOpenAIChatModel(ctx)

	// 绑定工具到 ChatModel
	log.Printf("===bind tools to llm===\n")
	toolInfos := tools.GetAllToolInfos(ctx)
	var err error
	cm, err = cm.WithTools(toolInfos)
	if err != nil {
		log.Fatalf("bind tools failed: %v", err)
	}

	// 获取工具实例
	allTools := tools.GetAllTools(ctx)

	log.Printf("===first turn: llm generate===\n")
	result := flow.Generate(ctx, cm, mem, messages)
	log.Printf("result: %+v\n\n", result.Content)

	// 第二轮对话
	log.Printf("===second turn: create messages===\n")
	secondMessages := []*schema.Message{
		schema.UserMessage("我还是觉得很难，能再多鼓励我一下吗？"),
	}
	result2 := flow.Generate(ctx, cm, mem, secondMessages)
	log.Printf("result2: %+v\n\n", result2.Content)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("===third turn: llm stream generate===\n")
		thirdMessages := []*schema.Message{
			schema.UserMessage("谢谢你！"),
		}
		streamResult := flow.Stream(ctx, cm, mem, thirdMessages)
		flow.ReportStream(streamResult)
	}()

	wg.Wait()

	// 第四轮对话：工具调用示例
	log.Printf("\n===fourth turn: tool calling example===\n")
	toolMessages := []*schema.Message{
		schema.UserMessage("帮我查一下北京的天气，然后告诉我现在几点了"),
	}
	toolResult := flow.GenerateWithTools(ctx, cm, mem, toolMessages, allTools)
	log.Printf("tool result: %+v\n\n", toolResult.Content)

	// 第五轮对话：工具调用流式示例
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("\n===fifth turn: tool calling stream example===\n")
		toolStreamMessages := []*schema.Message{
			schema.UserMessage("帮我计算一下 25 * 4 等于多少"),
		}
		toolStreamResult := flow.StreamWithTools(ctx, cm, mem, toolStreamMessages, allTools)
		flow.ReportStream(toolStreamResult)
	}()
	wg.Wait()
}
