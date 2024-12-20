package main

import (
	"context"
	"log/slog"
)

type LambdaHandler struct {
	openAIClient   *OpenAIClient
	telegramClient *TelegramClient
	telegramChatID string
}

func NewLambdaHandler(openAIClient *OpenAIClient, telegramClient *TelegramClient, telegramChatID string) *LambdaHandler {
	return &LambdaHandler{
		openAIClient:   openAIClient,
		telegramClient: telegramClient,
		telegramChatID: telegramChatID,
	}
}

func (h *LambdaHandler) HandleRequest(ctx context.Context) {
	slog.InfoContext(ctx, "handle request")

	msg, err := h.openAIClient.GetNextMessage(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get next message", "error", err)
	}
	slog.InfoContext(ctx, "got message")

	if len(msg.Choices) == 0 {
		slog.InfoContext(ctx, "no choices")
		return
	}

	if err = h.telegramClient.SendMessage(ctx, h.telegramChatID, msg.Choices[0].Message.Content); err != nil {
		slog.ErrorContext(ctx, "failed to send message", "error", err)
	}

	slog.InfoContext(ctx, "request handled")
}
