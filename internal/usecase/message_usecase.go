package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alarm-agent/internal/adapters/llm"
	"github.com/alarm-agent/internal/adapters/whatsapp"
	"github.com/alarm-agent/internal/config"
	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type MessageUseCase struct {
	repos           ports.Repositories
	whatsappSender  ports.WhatsAppSender
	eventUseCase    *EventUseCase
	defaultTimezone string
	config          *config.Config
}

func NewMessageUseCase(
	repos ports.Repositories,
	whatsappSender ports.WhatsAppSender,
	eventUseCase *EventUseCase,
	defaultTimezone string,
	config *config.Config,
) *MessageUseCase {
	return &MessageUseCase{
		repos:           repos,
		whatsappSender:  whatsappSender,
		eventUseCase:    eventUseCase,
		defaultTimezone: defaultTimezone,
		config:          config,
	}
}

func (uc *MessageUseCase) ProcessInboundMessage(ctx context.Context, parsedMessage whatsapp.ParsedMessage) error {
	exists, err := uc.repos.InboundMessage().Exists(ctx, parsedMessage.ID)
	if err != nil {
		return fmt.Errorf("failed to check message existence: %w", err)
	}
	if exists {
		return nil
	}

	rawPayload, _ := json.Marshal(parsedMessage)
	inboundMessage := &domain.InboundMessage{
		ProviderMessageID: parsedMessage.ID,
		FromNumber:        parsedMessage.From,
		RawPayload:        rawPayload,
	}

	if err := uc.repos.InboundMessage().Create(ctx, inboundMessage); err != nil {
		return fmt.Errorf("failed to create inbound message: %w", err)
	}

	isWhitelisted, err := uc.repos.Whitelist().IsWhitelisted(ctx, parsedMessage.From)
	if err != nil {
		return fmt.Errorf("failed to check whitelist: %w", err)
	}
	if !isWhitelisted {
		return nil
	}

	return uc.processWhitelistedMessage(ctx, parsedMessage)
}

func (uc *MessageUseCase) processWhitelistedMessage(ctx context.Context, parsedMessage whatsapp.ParsedMessage) error {
	user, err := uc.getOrCreateUser(ctx, parsedMessage.From, parsedMessage.ContactName)
	if err != nil {
		return fmt.Errorf("failed to get or create user: %w", err)
	}

	userPreferences := map[string]interface{}{
		"timezone":                         user.Timezone,
		"default_remind_before_minutes":    user.DefaultRemindBeforeMinutes,
		"default_remind_frequency_minutes": user.DefaultRemindFrequencyMinutes,
		"default_require_confirmation":     user.DefaultRequireConfirmation,
	}

	// Get LLM client from user's database configuration
	llmClient, err := llm.NewLLMClientFromDB(ctx, uc.repos.LLMConfig(), uc.config, user.ID)
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}

	systemPrompt := llm.BuildSystemPrompt(user.Timezone)
	userMessage := llm.BuildUserMessage(parsedMessage.From, parsedMessage.Text, userPreferences)

	llmResponse, err := llmClient.Chat(ctx, systemPrompt, userMessage)
	if err != nil {
		return fmt.Errorf("failed to get LLM response: %w", err)
	}

	if llmResponse.FollowUpQuestion != nil {
		return uc.sendWhatsAppMessage(ctx, parsedMessage.From, *llmResponse.FollowUpQuestion)
	}

	return uc.handleLLMIntent(ctx, user, llmResponse)
}

func (uc *MessageUseCase) handleLLMIntent(ctx context.Context, user *domain.User, llmResponse *domain.LLMResponse) error {
	switch llmResponse.Intent {
	case domain.IntentCreateEvent:
		return uc.handleCreateEvent(ctx, user, llmResponse)
	case domain.IntentUpdateEvent:
		return uc.handleUpdateEvent(ctx, user, llmResponse)
	case domain.IntentCancelEvent:
		return uc.handleCancelEvent(ctx, user, llmResponse)
	case domain.IntentListEvents:
		return uc.handleListEvents(ctx, user, llmResponse)
	case domain.IntentConfirmEvent:
		return uc.handleConfirmEvent(ctx, user, llmResponse)
	case domain.IntentDeclineEvent:
		return uc.handleDeclineEvent(ctx, user, llmResponse)
	case domain.IntentSmallTalk:
		return uc.sendWhatsAppMessage(ctx, user.WANumber, "Olá! Como posso ajudar com seus compromissos hoje?")
	default:
		return uc.sendWhatsAppMessage(ctx, user.WANumber, "Desculpe, não consegui entender sua mensagem. Pode tentar novamente?")
	}
}

func (uc *MessageUseCase) handleCreateEvent(ctx context.Context, user *domain.User, llmResponse *domain.LLMResponse) error {
	entities, err := uc.parseEventEntities(llmResponse.Entities)
	if err != nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, "Erro ao processar os dados do evento. Pode tentar novamente?")
	}

	event, err := uc.eventUseCase.CreateEvent(ctx, user.ID, entities)
	if err != nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, fmt.Sprintf("Erro ao criar evento: %s", err.Error()))
	}

	location := ""
	if event.Location != nil {
		location = fmt.Sprintf(" em %s", *event.Location)
	}

	message := fmt.Sprintf("✅ Evento criado: %s em %s%s. Lembrete: %d minutos antes.",
		event.Title,
		event.StartsAt.Format("02/01/2006 15:04"),
		location,
		event.RemindBeforeMinutes,
	)

	return uc.sendWhatsAppMessage(ctx, user.WANumber, message)
}

func (uc *MessageUseCase) handleUpdateEvent(ctx context.Context, user *domain.User, llmResponse *domain.LLMResponse) error {
	entities, err := uc.parseEventEntities(llmResponse.Entities)
	if err != nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, "Erro ao processar os dados do evento.")
	}

	event, err := uc.eventUseCase.UpdateEvent(ctx, user.ID, entities)
	if err != nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, fmt.Sprintf("Erro ao atualizar evento: %s", err.Error()))
	}

	message := fmt.Sprintf("✏️ Evento atualizado: %s em %s", event.Title, event.StartsAt.Format("02/01/2006 15:04"))
	return uc.sendWhatsAppMessage(ctx, user.WANumber, message)
}

func (uc *MessageUseCase) handleCancelEvent(ctx context.Context, user *domain.User, llmResponse *domain.LLMResponse) error {
	entities, err := uc.parseEventEntities(llmResponse.Entities)
	if err != nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, "Erro ao identificar o evento.")
	}

	event, err := uc.eventUseCase.CancelEvent(ctx, user.ID, entities.Identifier)
	if err != nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, fmt.Sprintf("Erro ao cancelar evento: %s", err.Error()))
	}

	message := fmt.Sprintf("❌ Evento cancelado: %s", event.Title)
	return uc.sendWhatsAppMessage(ctx, user.WANumber, message)
}

func (uc *MessageUseCase) handleListEvents(ctx context.Context, user *domain.User, llmResponse *domain.LLMResponse) error {
	events, err := uc.eventUseCase.ListEvents(ctx, user.ID, nil, nil)
	if err != nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, "Erro ao listar eventos.")
	}

	if len(events) == 0 {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, "Você não tem nenhum evento agendado.")
	}

	var message strings.Builder
	message.WriteString("📅 *Seus próximos eventos:*\n\n")

	for i, event := range events {
		if i >= 10 { // Limit to 10 events
			break
		}
		location := ""
		if event.Location != nil {
			location = fmt.Sprintf(" - %s", *event.Location)
		}
		message.WriteString(fmt.Sprintf("%d. %s\n📅 %s%s\n\n",
			i+1,
			event.Title,
			event.StartsAt.Format("02/01/2006 15:04"),
			location,
		))
	}

	return uc.sendWhatsAppMessage(ctx, user.WANumber, message.String())
}

func (uc *MessageUseCase) handleConfirmEvent(ctx context.Context, user *domain.User, llmResponse *domain.LLMResponse) error {
	entities, err := uc.parseEventEntities(llmResponse.Entities)
	if err != nil || entities.Identifier == nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, "Não consegui identificar qual evento confirmar.")
	}

	event, err := uc.eventUseCase.ConfirmEvent(ctx, user.ID, entities.Identifier)
	if err != nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, fmt.Sprintf("Erro ao confirmar evento: %s", err.Error()))
	}

	message := fmt.Sprintf("✅ Evento confirmado: %s", event.Title)
	return uc.sendWhatsAppMessage(ctx, user.WANumber, message)
}

func (uc *MessageUseCase) handleDeclineEvent(ctx context.Context, user *domain.User, llmResponse *domain.LLMResponse) error {
	entities, err := uc.parseEventEntities(llmResponse.Entities)
	if err != nil || entities.Identifier == nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, "Não consegui identificar qual evento cancelar.")
	}

	event, err := uc.eventUseCase.CancelEvent(ctx, user.ID, entities.Identifier)
	if err != nil {
		return uc.sendWhatsAppMessage(ctx, user.WANumber, fmt.Sprintf("Erro ao cancelar evento: %s", err.Error()))
	}

	message := fmt.Sprintf("❌ Evento cancelado: %s", event.Title)
	return uc.sendWhatsAppMessage(ctx, user.WANumber, message)
}

func (uc *MessageUseCase) getOrCreateUser(ctx context.Context, waNumber, contactName string) (*domain.User, error) {
	user, err := uc.repos.User().GetByWANumber(ctx, waNumber)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	user = &domain.User{
		WANumber:                      waNumber,
		Timezone:                      uc.defaultTimezone,
		DefaultRemindBeforeMinutes:    30,
		DefaultRemindFrequencyMinutes: 15,
		DefaultRequireConfirmation:    true,
	}

	if contactName != "" {
		user.Name = &contactName
	}

	if err := uc.repos.User().Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *MessageUseCase) parseEventEntities(entities map[string]interface{}) (*domain.EventEntities, error) {
	entitiesJSON, err := json.Marshal(entities)
	if err != nil {
		return nil, err
	}

	var eventEntities domain.EventEntities
	if err := json.Unmarshal(entitiesJSON, &eventEntities); err != nil {
		return nil, err
	}

	return &eventEntities, nil
}

func (uc *MessageUseCase) sendWhatsAppMessage(ctx context.Context, to, text string) error {
	return uc.whatsappSender.SendText(ctx, to, text)
}
