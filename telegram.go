package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type (
	TelegramRequest struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}

	TelegramClient struct {
		token  string
		client *http.Client
	}
)

func NewTelegramClient(token string, client *http.Client) *TelegramClient {
	return &TelegramClient{
		token:  token,
		client: client,
	}
}

func (c *TelegramClient) SendMessage(ctx context.Context, chatID, text string) error {
	body, err := json.Marshal(TelegramRequest{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("unexpected status code", "status_code", resp.StatusCode)
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("read response body", "error", err)
		} else {
			slog.Warn("response body", "body", string(respBytes))
		}
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
