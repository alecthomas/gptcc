package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	openai "github.com/sashabaranov/go-openai"
)

var cli struct {
	Token   string   `env:"OPENAI_API_KEY" help:"OpenAI API token." required:""`
	Message []string `arg:"" help:"Commit message to convert to Conventional Commits."`
}

func main() {
	kctx := kong.Parse(&cli)
	client := openai.NewClient(cli.Token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `
					Input is a commit message. Output is a commit message with a Conventional Commits (CC) prefix.
     
					If the input already has a CC prefix just return it, otherwise add it.
					Do not otherwise modify the input unless the prefix looks like a scope is already there.
					Try to infer the scope from the message.
			   
					Do not output anything other than the possibly modified commit message.
					`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: strings.Join(cli.Message, " "),
				},
			},
		},
	)

	kctx.FatalIfErrorf(err)
	fmt.Println(resp.Choices[0].Message.Content)
}
