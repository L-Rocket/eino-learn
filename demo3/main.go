package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

func main() {
	ctx := context.Background()
	model := openai.NewChatModel("gpt-4o-mini")
	resp, _ := model.Invoke(ctx, "一句话解释什么是黑洞")
	fmt.Println(resp)

}
