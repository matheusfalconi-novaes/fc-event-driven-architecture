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

type BalanceUpdatedOutputDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIdTo   float64 `json:"balance_account_id_to"`
}

type CreateTransactionUseCase struct {
	UnitOfWork         unityofwork.UnityOfWorkInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated     events.EventInterface
}

func NewCreateTransactionUseCase(
	unitOfWork unityofwork.UnityOfWorkInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	balanceUpdated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		UnitOfWork:         unitOfWork,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
		BalanceUpdated:     balanceUpdated,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	transactionOutput := &CreateTransactionOutputDTO{}
	balanceOutput := &BalanceUpdatedOutputDTO{}
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
		transactionOutput.ID = transaction.ID
		transactionOutput.AccountIDFrom = input.AccountIDFrom
		transactionOutput.AccountIDTo = input.AccountIDTo
		transactionOutput.Amount = input.Amount

		balanceOutput.AccountIDFrom = input.AccountIDFrom
		balanceOutput.AccountIDTo = input.AccountIDTo
		balanceOutput.BalanceAccountIDFrom = accountFrom.Balance
		balanceOutput.BalanceAccountIdTo = accountTo.Balance
		return nil
	})

	if err != nil {
		return nil, err
	}

	uc.TransactionCreated.SetPayload(transactionOutput)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)

	uc.BalanceUpdated.SetPayload(balanceOutput)
	uc.EventDispatcher.Dispatch(uc.BalanceUpdated)

	return transactionOutput, nil
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
