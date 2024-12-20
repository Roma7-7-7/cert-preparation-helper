package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	CompletionsRequest struct {
		Model               string           `json:"model"`
		Messages            []RequestMessage `json:"messages"`
		ResponseFormat      ResponseFormat   `json:"response_format"`
		Temperature         float64          `json:"temperature"`
		MaxCompletionTokens int              `json:"max_completion_tokens"`
		TopP                float64          `json:"top_p"`
		FrequencyPenalty    float64          `json:"frequency_penalty"`
		PresencePenalty     float64          `json:"presence_penalty"`
	}

	RequestMessage struct {
		Role    string    `json:"role"`
		Content []Content `json:"content"`
	}

	Content struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}

	ResponseFormat struct {
		Type string `json:"type"`
	}

	CompletionsResponse struct {
		ID                string   `json:"id"`
		Object            string   `json:"object"`
		Created           int64    `json:"created"`
		Model             string   `json:"model"`
		Choices           []Choice `json:"choices"`
		Usage             Usage    `json:"usage"`
		SystemFingerprint string   `json:"system_fingerprint"`
	}

	Choice struct {
		Index        int             `json:"index"`
		Message      ResponseMessage `json:"message"`
		FinishReason string          `json:"finish_reason"`
	}

	Usage struct {
		PromptTokens            int                    `json:"prompt_tokens"`
		CompletionTokens        int                    `json:"completion_tokens"`
		TotalTokens             int                    `json:"total_tokens"`
		PromptTokensDetails     TokenDetails           `json:"prompt_tokens_details"`
		CompletionTokensDetails CompletionTokenDetails `json:"completion_tokens_details"`
	}

	ResponseMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	TokenDetails struct {
		CachedTokens int `json:"cached_tokens"`
		AudioTokens  int `json:"audio_tokens"`
	}

	CompletionTokenDetails struct {
		ReasoningTokens          int `json:"reasoning_tokens"`
		AudioTokens              int `json:"audio_tokens"`
		AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
		RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
	}

	OpenAIClient struct {
		openAIAPIToken string
		client         *http.Client
	}
)

func NewOpenAIClient(openAIAPIToken string, client *http.Client) *OpenAIClient {
	return &OpenAIClient{
		openAIAPIToken: openAIAPIToken,
		client:         client,
	}
}

func (c *OpenAIClient) GetNextMessage(ctx context.Context) (*CompletionsResponse, error) {
	requestPayload := CompletionsRequest{
		Model: "gpt-4o",
		Messages: []RequestMessage{
			{
				Role: "system",
				Content: []Content{
					{
						Type: "text",
						Text: "You are chat bot that should help user with \"AWS Certified Solutions Architect - Associate\" certification preparation. The user will ask you for the next random fact about AWS and you should generate a summary of some AWS Service, deep detail/fact about AWS Service, nice to know feature or anything else that is required to pass certification exam. Response should be structured and formatted so it could be send to WhatsApp/Telegram chat. At the end of the message it would be nice to have links to documentation/blog post/any other document to see details. DO NOT format answer with Markdown since it is not well supported in target chat apps.",
					},
				},
			},
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: "give me next message",
					},
				},
			},
		},
		ResponseFormat: ResponseFormat{
			Type: "text",
		},
		Temperature:         1,
		MaxCompletionTokens: 2048,
		TopP:                1,
		FrequencyPenalty:    0,
		PresencePenalty:     0,
	}

	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("marshal request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.openAIAPIToken))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	var response CompletionsResponse
	if err = json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response body: %w", err)
	}

	return &response, nil
}
