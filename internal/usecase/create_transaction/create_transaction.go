package create_transaction

import (
	"context"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/gateway"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/events"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/unityofwork"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string  `json:"account_id_from"`
	AccountIDTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID            string  `json:"id"`
	AccountIDFrom string  `json:"account_id_from"`
	AccountIDTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionUseCase struct {
	UnitOfWork         unityofwork.UnityOfWorkInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
}

func NewCreateTransactionUseCase(
	unitOfWork unityofwork.UnityOfWorkInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		UnitOfWork:         unitOfWork,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	output := &CreateTransactionOutputDTO{}
	err := uc.UnitOfWork.Do(ctx, func(_ *unityofwork.UnityOfWork) error {
		accountGateway := uc.getAccountGateway(ctx)
		transactionGateway := uc.getTransactionGateway(ctx)

		accountFrom, err := accountGateway.FindByID(input.AccountIDFrom)
		if err != nil {
			return err
		}
		accountTo, err := accountGateway.FindByID(input.AccountIDTo)
		if err != nil {
			return err
		}
		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			return err
		}
		err = accountGateway.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}
		err = accountGateway.UpdateBalance(accountTo)
		if err != nil {
			return err
		}
		err = transactionGateway.Create(transaction)
		if err != nil {
			return err
		}
		output.ID = transaction.ID
		output.AccountIDFrom = input.AccountIDFrom
		output.AccountIDTo = input.AccountIDTo
		output.Amount = input.Amount
		return nil
	})

	if err != nil {
		return nil, err
	}

	uc.TransactionCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)

	return output, nil
}

func (uc *CreateTransactionUseCase) getAccountGateway(ctx context.Context) gateway.AccountGateway {
	accountGateway, err := uc.UnitOfWork.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return accountGateway.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionGateway(ctx context.Context) gateway.TransactionGateway {
	transactionGateway, err := uc.UnitOfWork.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}
	return transactionGateway.(gateway.TransactionGateway)
}
