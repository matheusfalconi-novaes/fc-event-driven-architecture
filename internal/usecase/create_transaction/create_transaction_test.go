package create_transaction

import (
	"context"
	"testing"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/event"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/mocks"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/events"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("Client1", "j@j.com")
	account1, _ := entity.NewAccount(client1)
	account1.Credit(1000)
	client2, _ := entity.NewClient("Client2", "j@j2.com")
	account2, _ := entity.NewAccount(client2)
	account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	dispatcher := events.NewEventDispatcher()
	transactionEvent := event.NewTransactionCreated()
	balanceEvent := event.NewBalanceUpdated()

	ctx := context.Background()

	uc := NewCreateTransactionUseCase(mockUow, dispatcher, transactionEvent, balanceEvent)
	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}
	output, err := uc.Execute(ctx, inputDto)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
