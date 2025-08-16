# Alarm Agent - WhatsApp Calendar Bot

Um serviço em Go que integra WhatsApp via Infobip com LLMs (Anthropic/OpenAI) para gerenciar compromissos e enviar lembretes automatizados.

## Características

- **Integração WhatsApp**: Recebe e responde mensagens via webhook Infobip
- **IA Conversacional**: Usa LLMs (Anthropic Claude ou OpenAI GPT) para interpretar comandos em português
- **Sistema de Agenda**: Cria, atualiza, cancela e lista compromissos por usuário
- **Lembretes Inteligentes**: Sistema de notificações configuráveis com confirmação opcional
- **Segurança**: Whitelist de números, validação de webhooks, rate limiting
- **Observabilidade**: Métricas Prometheus, logs estruturados, health checks
- **Arquitetura Limpa**: Clean/Hexagonal Architecture com separação clara de responsabilidades

## Funcionalidades

### Comandos Suportados (Português-BR)

- **Criar**: "Marcar dentista dia 22/08 às 14h, lembrar 1h antes"
- **Atualizar**: "Adia a reunião para amanhã 9:30"
- **Cancelar**: "Cancelar o café com Ana sexta"
- **Listar**: "O que tenho semana que vem?"
- **Confirmar**: "OK", "Confirmo", "Sim"

### Sistema de Lembretes

- Configurável por usuário (tempo antes, frequência, max notificações)
- Opção de requerer confirmação do usuário
- Status do evento: scheduled → confirmed → completed
- Retry automático com backoff exponencial

## Arquitetura

```
cmd/
├── server/           # Ponto de entrada da aplicação
internal/
├── domain/          # Entidades e regras de negócio
├── usecase/         # Casos de uso (business logic)
├── ports/           # Interfaces (contratos)
├── adapters/        # Implementações (infra)
│   ├── repo/        # PostgreSQL repositories
│   ├── llm/         # OpenAI/Anthropic clients
│   ├── whatsapp/    # Infobip integration
│   └── http/        # HTTP handlers
├── config/          # Configuração da aplicação
├── infra/           # Database, logging, metrics
└── workers/         # Background jobs (scheduler)
db/
├── migrations/      # SQL migration files
pkg/                 # Packages compartilháveis
scripts/             # Scripts de desenvolvimento
```

## Tecnologias

- **Go 1.22** com modules
- **PostgreSQL** para persistência
- **Infobip** para integração WhatsApp
- **OpenAI/Anthropic** para processamento de linguagem natural
- **Prometheus** para métricas
- **Zap** para logging estruturado
- **Docker** para containerização

## Configuração

### Variáveis de Ambiente

```bash
# Aplicação
PORT=8080
ENV=development
TIMEZONE_DEFAULT=America/Sao_Paulo

# Database
POSTGRES_DSN=postgres://user:pass@localhost/alarm_agent?sslmode=disable

# WhatsApp/Infobip
INFOBIP_BASE_URL=https://api.infobip.com
INFOBIP_API_KEY=your_infobip_api_key
INFOBIP_WHATSAPP_SENDER=your_whatsapp_number
INFOBIP_WEBHOOK_SECRET=your_webhook_secret

# LLM Configuration
LLM_PROVIDER=anthropic  # ou openai
ANTHROPIC_API_KEY=your_anthropic_key
OPENAI_API_KEY=your_openai_key
LLM_MODEL=claude-3-haiku-20240307  # ou gpt-3.5-turbo

# Segurança
WHITELIST_NUMBERS=+5511999999999,+5511888888888
RATE_LIMIT_PER_MINUTE=30

# Workers
REMINDER_TICK_SECONDS=30
```

### Banco de Dados

O sistema usa PostgreSQL com migrações versionadas:

```sql
-- Principais tabelas:
users                # Usuários com preferências
whitelist_numbers    # Números autorizados
events              # Compromissos/lembretes
inbound_messages    # Cache para idempotência
```

## Desenvolvimento Local

### Pré-requisitos

- Go 1.22+
- Docker e Docker Compose
- Make

### Setup Rápido

```bash
# 1. Clone e configure
git clone <repo>
cd alarm-agent
cp .env.example .env
# Edite .env com suas credenciais

# 2. Suba a stack completa
make up

# 3. Execute migrações
make migrate-up

# 4. Teste a aplicação
make test

# 5. Simule um webhook
curl -X POST http://localhost:8080/webhook/whatsapp \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"from":"5511999999999","text":"Marcar dentista amanhã 14h"}]}'
```

### Comandos Disponíveis

```bash
make help          # Lista todos os comandos
make up             # Sobe stack (app + postgres)
make down           # Para stack
make build          # Build da aplicação
make test           # Executa testes
make test-cover     # Testes com cobertura
make lint           # Linting (golangci-lint)
make migrate-up     # Aplica migrações
make migrate-down   # Reverte migrações
make logs           # Logs da aplicação
```

## API Endpoints

### Webhooks
- `POST /webhook/whatsapp` - Recebe mensagens do Infobip

### Health & Metrics
- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /metrics` - Métricas Prometheus

### Desenvolvimento (opcional)
- `POST /api/v1/events` - CRUD de eventos (para debug)
- `GET /api/v1/events` - Lista eventos

## Testes

### Estrutura

- **Unit Tests**: Casos de uso isolados com mocks
- **Integration Tests**: Handlers HTTP + banco real
- **E2E Tests**: Fluxo completo com webhooks simulados

```bash
# Testes unitários
go test ./internal/usecase/... -v

# Testes de integração
go test ./internal/adapters/... -v -tags=integration

# Cobertura
make test-cover
```

## Deployment

### Docker

```bash
# Build da imagem
docker build -t alarm-agent .

# Executar
docker run -p 8080:8080 --env-file .env alarm-agent
```

### Docker Compose

```bash
# Produção
docker-compose -f docker-compose.prod.yml up -d

# Desenvolvimento
make up
```

## Monitoramento

### Métricas Disponíveis

- `whatsapp_messages_received_total`
- `whatsapp_messages_sent_total`
- `llm_requests_total`
- `events_created_total`
- `reminders_sent_total`
- `http_requests_duration_seconds`

### Health Checks

- `/health` - Status geral da aplicação
- `/ready` - Dependências (DB, APIs externas)

### Logs

Logs estruturados em JSON com níveis:
- `INFO` - Operações normais
- `WARN` - Situações de atenção
- `ERROR` - Erros que requerem investigação

## Segurança

- **Whitelist**: Apenas números autorizados podem usar o bot
- **Webhook Validation**: Validação de assinatura Infobip (se disponível)
- **Rate Limiting**: Limitação de requisições por minuto
- **PII Protection**: Dados pessoais não aparecem em logs
- **Environment Secrets**: Chaves via variáveis de ambiente

## FAQ

### Como adicionar um número na whitelist?
Adicione o número na variável `WHITELIST_NUMBERS` e reinicie a aplicação, ou insira diretamente na tabela `whitelist_numbers`.

### Como trocar o provedor de LLM?
Altere `LLM_PROVIDER=openai` (ou `anthropic`) e configure as respectivas API keys.

### Como customizar lembretes padrão?
Ajuste as preferências na tabela `users` ou permita que o usuário configure via mensagem.

### Webhook não está funcionando?
1. Verifique se o endpoint está acessível publicamente
2. Confirme as credenciais Infobip
3. Verifique os logs para erros de validação

## Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature
3. Commit suas mudanças
4. Execute testes e linting
5. Abra um Pull Request

## Licença

MIT License
