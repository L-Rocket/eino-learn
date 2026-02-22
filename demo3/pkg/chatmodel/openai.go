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

package chatmodel

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

func CreateOpenAIChatModel(ctx context.Context) model.ToolCallingChatModel {
	// key := os.Getenv("HUOSHAN_API_KEY")
	// modelName := "doubao-seed-1-6-lite-251015"
	// baseURL := "https://ark.cn-beijing.volces.com/api/v3/"
	key := os.Getenv("SILICONFLOW_API_KEY")
	modelName := "Qwen/Qwen3-8B"
	baseURL := "https://api.siliconflow.cn/v1"
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  key,
	})
	if err != nil {
		log.Fatalf("create openai chat model failed, err=%v", err)
	}
	return chatModel
}
