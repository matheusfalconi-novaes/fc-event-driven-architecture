package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	client1, _ := NewClient("John Doe 1", "j@j.com")
	account1, _ := NewAccount(client1)
	client2, _ := NewClient("John Doe 2", "j@j2.com")
	account2, _ := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 100)
	assert.Nil(t, err)
	assert.NotNil(t, transaction)

	assert.Equal(t, 900.0, account1.Balance)
	assert.Equal(t, 1100.0, account2.Balance)
}

func TestCreateTransactionWithInvalidAccounts(t *testing.T) {
	client, _ := NewClient("John Doe 1", "j@j.com")
	account, _ := NewAccount(client)

	transaction, err := NewTransaction(account, nil, 100)
	assert.Error(t, err, "none of the accounts can be nil")
	assert.Nil(t, transaction)

	transaction, err = NewTransaction(nil, account, 100)
	assert.Error(t, err, "none of the accounts can be nil")
	assert.Nil(t, transaction)
}

func TestCreateTransactionWithInvalidAmount(t *testing.T) {
	client1, _ := NewClient("John Doe 1", "j@j.com")
	account1, _ := NewAccount(client1)
	client2, _ := NewClient("John Doe 2", "j@j2.com")
	account2, _ := NewAccount(client2)

	transaction, err := NewTransaction(account1, account2, 0)

	assert.Error(t, err, "amount must be greater than zero")
	assert.Nil(t, transaction)

	transaction, err = NewTransaction(account1, account2, -10)

	assert.Error(t, err, "amount must be greater than zero")
	assert.Nil(t, transaction)
}

func TestCreateTransactionWithInsufficientBalance(t *testing.T) {
	client1, _ := NewClient("John Doe 1", "j@j.com")
	account1, _ := NewAccount(client1)
	client2, _ := NewClient("John Doe 2", "j@j2.com")
	account2, _ := NewAccount(client2)
	account1.Credit(1000)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 2000)

	assert.Error(t, err, "insufficient funds")
	assert.Nil(t, transaction)
	assert.Equal(t, 1000.0, account1.Balance)
	assert.Equal(t, 1000.0, account2.Balance)
}
