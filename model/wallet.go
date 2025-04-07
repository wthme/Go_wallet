package model

import "github.com/google/uuid"

type OperationType string

const (
	DEPOSIT  OperationType = "DEPOSIT"
	WITHDRAW OperationType = "WITHDRAW"
)

type WalletOperation struct {
	WalletID    uuid.UUID     `json:"walletId"`    
	OperationType OperationType `json:"operationType"`
	Amount      float64         `json:"amount"`
}

type Wallet struct {
	ID      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
}




