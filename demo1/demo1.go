package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func main() {
	url := "https://api.siliconflow.cn/v1/chat/completions"
	apiKey := os.Getenv("SILICONFLOW_API_KEY")

	requestBody := ChatRequest{
		// Model: "deepseek-ai/DeepSeek-R1-0528-Qwen3-8B",
		Model: "Qwen/Qwen3-8B",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Give me a hello world example in Go."},
		},
	}
	jsonData, _ := json.Marshal(requestBody)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		panic(err)
	}

	var result ChatResponse
	for {
		// set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		// send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close() // close body, otherwise it will cause resource leak

		if err := json.Unmarshal(body, &result); err != nil {
			time.Sleep(1 * time.Second) // wait for 1 second before retrying
			fmt.Println("Error unmarshalling response:", err)
		} else {
			break
		}

	}

	if len(result.Choices) > 0 {
		fmt.Println("Response:", result.Choices[0].Message.Content)
	} else {
		fmt.Println("No response received.")
	}
}
