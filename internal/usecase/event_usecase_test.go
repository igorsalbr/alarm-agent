package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alarm-agent/internal/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByWANumber(ctx context.Context, waNumber string) (*domain.User, error) {
	args := m.Called(ctx, waNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(ctx context.Context, event *domain.Event) error {
	args := m.Called(ctx, event)
	if args.Error(0) == nil {
		event.ID = 1
	}
	return args.Error(0)
}

func (m *MockEventRepository) Update(ctx context.Context, event *domain.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventRepository) GetByID(ctx context.Context, id int) (*domain.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Event), args.Error(1)
}

func (m *MockEventRepository) GetByUserID(ctx context.Context, userID int) ([]domain.Event, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func (m *MockEventRepository) GetByUserIDAndDateRange(ctx context.Context, userID int, start, end time.Time) ([]domain.Event, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func (m *MockEventRepository) GetPendingReminders(ctx context.Context, reminderWindow time.Duration) ([]domain.EventWithUser, error) {
	args := m.Called(ctx, reminderWindow)
	return args.Get(0).([]domain.EventWithUser), args.Error(1)
}

func (m *MockEventRepository) FindByUserAndIdentifier(ctx context.Context, userID int, identifier *domain.EventIdentifier) ([]domain.Event, error) {
	args := m.Called(ctx, userID, identifier)
	return args.Get(0).([]domain.Event), args.Error(1)
}

type MockRepositories struct {
	userRepo  *MockUserRepository
	eventRepo *MockEventRepository
}

func (m *MockRepositories) User() *MockUserRepository {
	return m.userRepo
}

func (m *MockRepositories) Event() *MockEventRepository {
	return m.eventRepo
}

func (m *MockRepositories) Whitelist() interface{} {
	return nil
}

func (m *MockRepositories) InboundMessage() interface{} {
	return nil
}

func TestEventUseCase_CreateEvent(t *testing.T) {
	ctx := context.Background()
	
	mockRepos := &MockRepositories{
		userRepo:  &MockUserRepository{},
		eventRepo: &MockEventRepository{},
	}
	
	useCase := NewEventUseCase(mockRepos)
	
	title := "Test Event"
	startsAt := time.Now().Add(time.Hour)
	
	entities := &domain.EventEntities{
		Title:    &title,
		StartsAt: &startsAt,
	}
	
	user := &domain.User{
		ID:                              1,
		WANumber:                        "+5511999999999",
		DefaultRemindBeforeMinutes:      30,
		DefaultRemindFrequencyMinutes:   15,
		DefaultRequireConfirmation:      true,
	}
	
	mockRepos.userRepo.On("GetByWANumber", ctx, mock.AnythingOfType("string")).Return(user, nil)
	mockRepos.eventRepo.On("Create", ctx, mock.AnythingOfType("*domain.Event")).Return(nil)
	
	event, err := useCase.CreateEvent(ctx, 1, entities)
	
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, title, event.Title)
	assert.Equal(t, startsAt, event.StartsAt)
	assert.Equal(t, domain.EventStatusScheduled, event.Status)
	assert.Equal(t, user.DefaultRemindBeforeMinutes, event.RemindBeforeMinutes)
	
	mockRepos.eventRepo.AssertExpectations(t)
}

func TestEventUseCase_CreateEvent_MissingTitle(t *testing.T) {
	ctx := context.Background()
	
	mockRepos := &MockRepositories{
		userRepo:  &MockUserRepository{},
		eventRepo: &MockEventRepository{},
	}
	
	useCase := NewEventUseCase(mockRepos)
	
	startsAt := time.Now().Add(time.Hour)
	
	entities := &domain.EventEntities{
		StartsAt: &startsAt,
	}
	
	event, err := useCase.CreateEvent(ctx, 1, entities)
	
	assert.Error(t, err)
	assert.Nil(t, event)
	assert.Contains(t, err.Error(), "title is required")
}

func TestEventUseCase_CreateEvent_MissingStartTime(t *testing.T) {
	ctx := context.Background()
	
	mockRepos := &MockRepositories{
		userRepo:  &MockUserRepository{},
		eventRepo: &MockEventRepository{},
	}
	
	useCase := NewEventUseCase(mockRepos)
	
	title := "Test Event"
	
	entities := &domain.EventEntities{
		Title: &title,
	}
	
	event, err := useCase.CreateEvent(ctx, 1, entities)
	
	assert.Error(t, err)
	assert.Nil(t, event)
	assert.Contains(t, err.Error(), "start time is required")
}