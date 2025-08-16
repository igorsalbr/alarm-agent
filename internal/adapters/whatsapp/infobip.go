package whatsapp

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alarm-agent/internal/ports"
)

type InfobipClient struct {
	baseURL    string
	apiKey     string
	sender     string
	httpClient *http.Client
}

func NewInfobipClient(baseURL, apiKey, sender string) ports.WhatsAppSender {
	return &InfobipClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		sender:  sender,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type InfobipTextMessage struct {
	From    string             `json:"from"`
	To      string             `json:"to"`
	Content InfobipTextContent `json:"content"`
}

type InfobipTextContent struct {
	Text string `json:"text"`
}

type InfobipSendRequest struct {
	Messages []InfobipTextMessage `json:"messages"`
}

func (c *InfobipClient) SendText(ctx context.Context, to, text string) error {
	request := InfobipSendRequest{
		Messages: []InfobipTextMessage{
			{
				From: c.sender,
				To:   to,
				Content: InfobipTextContent{
					Text: text,
				},
			},
		},
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/whatsapp/1/message/text", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("App %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("infobip API error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

type InfobipWebhookVerifier struct {
	secret string
}

func NewInfobipWebhookVerifier(secret string) ports.WhatsAppWebhookVerifier {
	return &InfobipWebhookVerifier{secret: secret}
}

func (v *InfobipWebhookVerifier) VerifySignature(payload []byte, signature string) bool {
	if v.secret == "" || signature == "" {
		return true // Skip verification if no secret configured
	}

	mac := hmac.New(sha256.New, []byte(v.secret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
