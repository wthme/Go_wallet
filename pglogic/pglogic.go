package pglogic

import (
	"Go_wallet/model"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Wallet_actions interface {
	GetWallet(ctx context.Context, walletID uuid.UUID) (*model.Wallet, error)
	UpdateWalletBalance(ctx context.Context, walletID uuid.UUID, amount float64) error
	CreateWallet(ctx context.Context, walletID uuid.UUID) error
}

type wallet struct {
	db *pgxpool.Pool
}

func NewWallet(db *pgxpool.Pool) *wallet {
	return &wallet{db: db}
}

func (r *wallet) GetWallet(ctx context.Context, walletID uuid.UUID) (*model.Wallet, error) {

	var wallet model.Wallet
	err := r.db.QueryRow(ctx,
		"SELECT id, balance FROM wallets WHERE id = $1", walletID).
		Scan(&wallet.ID, &wallet.Balance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

func (r *wallet) UpdateWalletBalance(ctx context.Context, walletID uuid.UUID, amount float64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Use SELECT FOR UPDATE to lock the row
	var balance float64
	err = tx.QueryRow(ctx,
		"SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", walletID).
		Scan(&balance)
	if err != nil {
		return fmt.Errorf("failed to lock wallet: %w", err)
	}

	// Update the balance
	newBalance := balance + amount
	if newBalance < 0 {
		return fmt.Errorf("insufficient funds")
	}

	_, err = tx.Exec(ctx,
		"UPDATE wallets SET balance = $1 WHERE id = $2", newBalance, walletID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *wallet) CreateWallet(ctx context.Context, walletID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO wallets (id, balance) VALUES ('$1', 0)", walletID)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}
	return nil
}


func InitDB(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS wallets (
			id UUID PRIMARY KEY,
			balance BIGINT NOT NULL DEFAULT 0,
			CHECK (balance >= 0)
		);
	`)
	return err
}