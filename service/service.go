package service

import (
	"Go_wallet/model"
	"Go_wallet/pglogic"
	"context"

	"github.com/google/uuid"
)

type WalletService interface {
	ProcessOperation(ctx context.Context, op model.WalletOperation) error
	GetWalletBalance(ctx context.Context, walletID uuid.UUID) (float64, error)
}

type walletService struct {
	repo pglogic.Wallet_actions
}

func NewWalletService(rep pglogic.Wallet_actions) *walletService {
	return &walletService{repo: rep}
}

func (s *walletService) ProcessOperation(ctx context.Context, op model.WalletOperation) error {
	// Check if wallet exists
	wallet, err := s.repo.GetWallet(ctx, op.WalletID)
	if err != nil {
		return err
	}

	// Create wallet if it doesn't exist
	if wallet == nil {
		if err := s.repo.CreateWallet(ctx, op.WalletID); err != nil {
			return err
		}
	}

	// Adjust amount based on operation type
	amount := op.Amount
	if op.OperationType == model.WITHDRAW {
		amount = -amount
	}

	// Update balance
	return s.repo.UpdateWalletBalance(ctx, op.WalletID, amount)
}

func (s *walletService) GetWalletBalance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	wallet, err := s.repo.GetWallet(ctx, walletID)
	if err != nil {
		return 0, err
	}
	if wallet == nil {
		return 0, nil
	}
	return wallet.Balance, nil
}