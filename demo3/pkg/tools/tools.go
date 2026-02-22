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

package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type GetWeatherRequest struct {
	City string `json:"city"`
}

type GetWeatherResponse struct {
	City        string `json:"city"`
	Temperature int    `json:"temperature"`
	Condition   string `json:"condition"`
	Humidity    int    `json:"humidity"`
}

func NewGetWeatherTool() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "get_weather",
			Desc: "获取指定城市的天气信息",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"city": {
					Type: "string",
					Desc: "城市名称，例如：北京、上海、广州",
				},
			}),
		},
		func(ctx context.Context, input *GetWeatherRequest) (output *GetWeatherResponse, err error) {
			weatherData := map[string]GetWeatherResponse{
				"北京": {City: "北京", Temperature: 22, Condition: "晴", Humidity: 45},
				"上海": {City: "上海", Temperature: 25, Condition: "多云", Humidity: 65},
				"广州": {City: "广州", Temperature: 28, Condition: "小雨", Humidity: 80},
				"深圳": {City: "深圳", Temperature: 27, Condition: "阴", Humidity: 75},
				"杭州": {City: "杭州", Temperature: 23, Condition: "晴", Humidity: 55},
			}

			result, ok := weatherData[input.City]
			if !ok {
				return &GetWeatherResponse{
					City:        input.City,
					Temperature: 20,
					Condition:   "未知",
					Humidity:    50,
				}, nil
			}

			return &result, nil
		},
	)
}

type GetCurrentTimeRequest struct {
	Timezone string `json:"timezone"`
}

type GetCurrentTimeResponse struct {
	CurrentTime string `json:"current_time"`
	Timezone    string `json:"timezone"`
	Date        string `json:"date"`
}

func NewGetCurrentTimeTool() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "get_current_time",
			Desc: "获取当前时间，支持指定时区",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"timezone": {
					Type: "string",
					Desc: "时区，例如：Asia/Shanghai, America/New_York，默认为 Asia/Shanghai",
				},
			}),
		},
		func(ctx context.Context, input *GetCurrentTimeRequest) (output *GetCurrentTimeResponse, err error) {
			timezone := input.Timezone
			if timezone == "" {
				timezone = "Asia/Shanghai"
			}

			loc, err := time.LoadLocation(timezone)
			if err != nil {
				loc = time.FixedZone("UTC", 0)
			}

			now := time.Now().In(loc)
			return &GetCurrentTimeResponse{
				CurrentTime: now.Format("15:04:05"),
				Timezone:    timezone,
				Date:        now.Format("2006-01-02"),
			}, nil
		},
	)
}

type CalculatorRequest struct {
	Expression string `json:"expression"`
}

type CalculatorResponse struct {
	Result     float64 `json:"result"`
	Expression string  `json:"expression"`
	Error      string  `json:"error,omitempty"`
}

func NewCalculatorTool() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "calculator",
			Desc: "执行简单的数学计算，支持加减乘除运算",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"expression": {
					Type: "string",
					Desc: "数学表达式，例如：1+2, 10*5, 100/4, 20-5",
				},
			}),
		},
		func(ctx context.Context, input *CalculatorRequest) (output *CalculatorResponse, err error) {
			var result float64
			var op rune
			var num1, num2 float64

			n, err := fmt.Sscanf(input.Expression, "%f%c%f", &num1, &op, &num2)
			if err != nil || n != 3 {
				return &CalculatorResponse{
					Expression: input.Expression,
					Error:      "invalid expression format",
				}, nil
			}

			switch op {
			case '+':
				result = num1 + num2
			case '-':
				result = num1 - num2
			case '*':
				result = num1 * num2
			case '/':
				if num2 == 0 {
					return &CalculatorResponse{
						Expression: input.Expression,
						Error:      "division by zero",
					}, nil
				}
				result = num1 / num2
			default:
				return &CalculatorResponse{
					Expression: input.Expression,
					Error:      "unsupported operator",
				}, nil
			}

			return &CalculatorResponse{
				Result:     result,
				Expression: input.Expression,
			}, nil
		},
	)
}

func GetAllTools(ctx context.Context) []tool.BaseTool {
	return []tool.BaseTool{
		NewGetWeatherTool(),
		NewGetCurrentTimeTool(),
		NewCalculatorTool(),
	}
}

func GetAllToolInfos(ctx context.Context) []*schema.ToolInfo {
	tools := GetAllTools(ctx)
	var toolInfos []*schema.ToolInfo
	for _, t := range tools {
		info, err := t.Info(ctx)
		if err != nil {
			continue
		}
		toolInfos = append(toolInfos, info)
	}
	return toolInfos
}
