package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"tools-cart/client"
	"tools-cart/tools"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
)

func main() {
	catalog := tools.GetToolsCatalog()

	jsonOutputModel := "ai/qwen2.5:latest"

	systemContentIntroduction := `You have access to the following tools:`
	// make a JSON String from the content of tools
	toolsJson, err := json.Marshal(catalog)
	if err != nil {
		panic("Error marshalling tools catalog to JSON: " + err.Error())
	}
	// Create a prompt from the tools catalog
	toolsContent := "[AVAILABLE_TOOLS]" + string(toolsJson) + "[/AVAILABLE_TOOLS]"

	systemContentInstructions := `If the question of the user matched the description of a tool, the tool will be called.
	To call a tool, respond with a JSON object with the following structure: 
	[
		{
			"name": <name of the called tool>,
			"arguments": {
				<name of the argument>: <value of the argument>
			}
		},
	]
	
	search the name of the tool in the list of tools with the Name field
	`

	userMessage := `
		search the Dune book in books 
		search all books with a limit of 5 found books
		search all electronics with a limit of 3 found electronics

		add 3 iPad Pro 12.9 to the cart
		add 2 macbook air M3 to the cart
		add 5 Sapiens book to the cart and 2 Dune book to the cart
		remove iPad Pro 12.9 from the cart

		view the cart
	`

	engine, err := client.GetDMRClient()
	if err != nil {
		panic("Error getting DMR client: " + err.Error())
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemContentIntroduction + "\n" + toolsContent + "\n" + systemContentInstructions),
			openai.UserMessage(userMessage),
		},
		Model:       jsonOutputModel,
		Temperature: openai.Opt(0.0),
	}
	completion, err := engine.Chat.Completions.New(context.Background(), params)

	if err != nil {
		panic("Error creating chat completion: " + err.Error())
	}
	if len(completion.Choices) == 0 {
		fmt.Println("No choices returned from chat completion")
		os.Exit(0)

	}
	result := completion.Choices[0].Message.Content
	if result == "" {
		fmt.Println("No content returned from chat completion")
		os.Exit(0)
	}

	// Create a new chat completion with the result to create a JSON string
	paramsNext := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("Return all function calls wrapped in a container object with a 'function_calls' key."),
			openai.UserMessage(result),
		},
		Model:       jsonOutputModel,
		Temperature: openai.Opt(0.0),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{
				Type: "json_object",
			},
		},
	}

	completionNext, err := engine.Chat.Completions.New(context.Background(), paramsNext)
	if err != nil {
		panic("Error creating chat completion for next step: " + err.Error())
	}
	if len(completionNext.Choices) == 0 {
		fmt.Println("No choices returned from chat completion")
		os.Exit(0)
	}

	resultNext := completionNext.Choices[0].Message.Content
	if resultNext == "" {
		fmt.Println("No content returned from chat completion")
		os.Exit(0)
	}
	// Print the result of the chat completion
	//fmt.Println("Result of the chat completion:")
	//fmt.Println(resultNext)

	// Create the list of Tool Calls from the result
	type Command struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	type FunctionCalls struct {
		FunctionCalls []Command `json:"function_calls"`
	}

	var commands FunctionCalls

	errJson := json.Unmarshal([]byte(resultNext), &commands)
	if errJson != nil {
		panic("Error unmarshalling result to commands: " + errJson.Error())
	}

	if len(commands.FunctionCalls) == 0 {
		fmt.Println("No commands found in the result")
		os.Exit(0)
	}

	toolCalls := make([]openai.ChatCompletionMessageToolCall, len(commands.FunctionCalls))

	for i, command := range commands.FunctionCalls {
		// transform command.Arguments to  JSON string
		argumentsJson, err := json.Marshal(command.Arguments)
		if err != nil {
			panic(fmt.Sprintf("Error marshalling arguments for command %s: %s", command.Name, err.Error()))
		}

		toolCalls[i] = openai.ChatCompletionMessageToolCall{
			ID: uuid.New().String(), // Generate a unique ID for the tool call
			Function: openai.ChatCompletionMessageToolCallFunction{
				Name:      command.Name,
				Arguments: string(argumentsJson),
			},
			Type: "function",
		}
	}

	fmt.Println("Tool Calls:")
	for _, toolCall := range toolCalls {
		fmt.Printf("Tool Name: %s, Arguments: %s\n", toolCall.Function.Name, toolCall.Function.Arguments)
	}

}
