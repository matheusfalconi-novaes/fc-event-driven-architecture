package create_transaction

import (
	"errors"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("Client1", "j@j.com")
	account1, _ := entity.NewAccount(client1)
	account1.Credit(1000)
	client2, _ := entity.NewClient("Client2", "j@j2.com")
	account2, _ := entity.NewAccount(client2)
	account2.Credit(1000)

	transactionGatewayMock := &TransactionGatewayMock{}
	accountGatewayMock := &AccountGatewayMock{}
	accountGatewayMock.On("FindByID", account1.ID).Return(account1, nil)
	accountGatewayMock.On("FindByID", account2.ID).Return(account2, nil)
	transactionGatewayMock.On("Create", mock.Anything).Return(nil)

	uc := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)
	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}
	output, err := uc.Execute(inputDto)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.True(t, IsValidUUID(output.ID))
	transactionGatewayMock.AssertExpectations(t)
	transactionGatewayMock.AssertNumberOfCalls(t, "Create", 1)
	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "FindByID", 2)
}

func TestCreateTransactionUseCase_ExecuteWithInvalidAccountIds(t *testing.T) {
	client1, _ := entity.NewClient("Client1", "j@j.com")
	account1, _ := entity.NewAccount(client1)
	account1.Credit(1000)
	client2, _ := entity.NewClient("Client2", "j@j2.com")
	account2, _ := entity.NewAccount(client2)
	account2.Credit(1000)

	transactionGatewayMock := &TransactionGatewayMock{}
	accountGatewayMock := &AccountGatewayMock{}
	accountGatewayMock.On("FindByID", account1.ID).Return(account1, errors.New(mock.Anything))

	uc := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)
	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}
	output, err := uc.Execute(inputDto)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	transactionGatewayMock.AssertExpectations(t)
	transactionGatewayMock.AssertNotCalled(t, "Create")
	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "FindByID", 1)

	accountGatewayMock.On("FindByID", account2.ID).Return(account2, errors.New(mock.Anything))

	uc = NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)
	output, err = uc.Execute(inputDto)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	transactionGatewayMock.AssertNotCalled(t, "Create")
	accountGatewayMock.AssertNumberOfCalls(t, "FindByID", 2)
}

func TestCreateTransactionUseCase_ExecuteWithInvalidTransaction(t *testing.T) {
	client1, _ := entity.NewClient("Client1", "j@j.com")
	account1, _ := entity.NewAccount(client1)
	client2, _ := entity.NewClient("Client2", "j@j2.com")
	account2, _ := entity.NewAccount(client2)

	transactionGatewayMock := &TransactionGatewayMock{}
	accountGatewayMock := &AccountGatewayMock{}
	accountGatewayMock.On("FindByID", account1.ID).Return(account1, nil)
	accountGatewayMock.On("FindByID", account2.ID).Return(account2, nil)

	uc := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)
	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}
	output, err := uc.Execute(inputDto)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	transactionGatewayMock.AssertExpectations(t)
	transactionGatewayMock.AssertNotCalled(t, "Create")
	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "FindByID", 2)
}

func TestCreateTransactionUseCase_ExecuteWithErrorOnGateway(t *testing.T) {
	client1, _ := entity.NewClient("Client1", "j@j.com")
	account1, _ := entity.NewAccount(client1)
	account1.Credit(1000)
	client2, _ := entity.NewClient("Client2", "j@j2.com")
	account2, _ := entity.NewAccount(client2)
	account2.Credit(1000)

	transactionGatewayMock := &TransactionGatewayMock{}
	accountGatewayMock := &AccountGatewayMock{}
	accountGatewayMock.On("FindByID", account1.ID).Return(account1, nil)
	accountGatewayMock.On("FindByID", account2.ID).Return(account2, nil)
	transactionGatewayMock.On("Create", mock.Anything).Return(errors.New(mock.Anything))

	uc := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)
	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}
	output, err := uc.Execute(inputDto)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	transactionGatewayMock.AssertExpectations(t)
	transactionGatewayMock.AssertNumberOfCalls(t, "Create", 1)
	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "FindByID", 2)
}
