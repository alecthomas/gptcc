package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	openai "github.com/sashabaranov/go-openai"
)

var cli struct {
	Timeout time.Duration `default:"10s" help:"Timeout for the request."`
	Token   string        `env:"OPENAI_API_KEY" help:"OpenAI API token." required:""`
	Message []string      `arg:"" help:"Commit message to convert to Conventional Commits."`
}

func main() {
	kctx := kong.Parse(&cli)
	client := openai.NewClient(cli.Token)
	ctx, cancel := context.WithTimeout(context.Background(), cli.Timeout)
	defer cancel()
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `
					Input is a commit message. Output is a commit message with a Conventional Commits (CC) prefix.
     
					If the input already has a CC prefix just return it, otherwise add it.

					Do not otherwise modify the input unless the prefix looks like a scope is already there.
			   
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
