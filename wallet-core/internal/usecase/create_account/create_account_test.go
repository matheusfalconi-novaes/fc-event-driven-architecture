package create_account

import (
	"testing"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	clientMock := &mocks.ClientGatewayMock{}
	clientMock.On("Get", client.ID).Return(client, nil)
	accountMock := &mocks.AccountGatewayMock{}

	accountMock.On("Save", mock.Anything).Return(nil)

	inputDto := CreateAccountInputDTO{
		ClientID: client.ID,
	}
	uc := NewCreateAccountUseCase(accountMock, clientMock)
	output, err := uc.Execute(inputDto)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.True(t, IsValidUUID(output.ID))
	clientMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t, "Get", 1)
	accountMock.AssertExpectations(t)
	accountMock.AssertNumberOfCalls(t, "Save", 1)
}
