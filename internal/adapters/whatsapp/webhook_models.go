package whatsapp

import (
	"encoding/json"
	"fmt"
	"time"
)

type InfobipWebhookRequest struct {
	Results []InfobipInboundMessage `json:"results"`
}

type InfobipInboundMessage struct {
	MessageID  string                `json:"messageId"`
	From       string                `json:"from"`
	To         string                `json:"to"`
	ReceivedAt time.Time             `json:"receivedAt"`
	Message    InfobipMessageContent `json:"message"`
	Contact    *InfobipContact       `json:"contact,omitempty"`
	Price      *InfobipPrice         `json:"price,omitempty"`
}

type InfobipMessageContent struct {
	Type     string                  `json:"type"`
	Text     *string                 `json:"text,omitempty"`
	Image    *InfobipMediaContent    `json:"image,omitempty"`
	Document *InfobipMediaContent    `json:"document,omitempty"`
	Audio    *InfobipMediaContent    `json:"audio,omitempty"`
	Video    *InfobipMediaContent    `json:"video,omitempty"`
	Location *InfobipLocationContent `json:"location,omitempty"`
}

type InfobipMediaContent struct {
	Caption *string `json:"caption,omitempty"`
	URL     string  `json:"url"`
}

type InfobipLocationContent struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      *string `json:"name,omitempty"`
	Address   *string `json:"address,omitempty"`
}

type InfobipContact struct {
	Name string `json:"name"`
}

type InfobipPrice struct {
	PricePerMessage float64 `json:"pricePerMessage"`
	Currency        string  `json:"currency"`
}

func (r *InfobipWebhookRequest) ExtractMessages() []ParsedMessage {
	var messages []ParsedMessage

	for _, result := range r.Results {
		message := ParsedMessage{
			ID:        result.MessageID,
			From:      result.From,
			To:        result.To,
			Timestamp: result.ReceivedAt,
			Type:      result.Message.Type,
		}

		switch result.Message.Type {
		case "TEXT":
			if result.Message.Text != nil {
				message.Text = *result.Message.Text
			}
		case "IMAGE":
			if result.Message.Image != nil {
				message.MediaURL = result.Message.Image.URL
				if result.Message.Image.Caption != nil {
					message.Text = *result.Message.Image.Caption
				}
			}
		case "LOCATION":
			if result.Message.Location != nil {
				locationData := map[string]interface{}{
					"latitude":  result.Message.Location.Latitude,
					"longitude": result.Message.Location.Longitude,
				}
				if result.Message.Location.Name != nil {
					locationData["name"] = *result.Message.Location.Name
				}
				if result.Message.Location.Address != nil {
					locationData["address"] = *result.Message.Location.Address
				}

				locationJSON, _ := json.Marshal(locationData)
				message.Text = string(locationJSON)
			}
		default:
			message.Text = fmt.Sprintf("Unsupported message type: %s", result.Message.Type)
		}

		if result.Contact != nil {
			message.ContactName = result.Contact.Name
		}

		messages = append(messages, message)
	}

	return messages
}

type ParsedMessage struct {
	ID          string    `json:"id"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`
	Text        string    `json:"text"`
	MediaURL    string    `json:"media_url,omitempty"`
	ContactName string    `json:"contact_name,omitempty"`
}
