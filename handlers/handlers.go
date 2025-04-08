package handlers

import (
	"context"
	"fmt"
	"net/http"

	"Go_wallet/model"
	"Go_wallet/pglogic"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	wallet pglogic.Walletdb
}

func NewWalletHandler(wallet pglogic.Walletdb) *WalletHandler {
	return &WalletHandler{wallet: wallet}
}

func (h *WalletHandler) HandleWalletOperation(c *gin.Context) {

	var op model.WalletOperation
	
	if err := c.BindJSON(&op); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error uncorrect request": err.Error()})
		return
	}
	
	fmt.Printf("%v %v %v", op.Amount , op.OperationType , op.WalletID)


	if op.OperationType != model.DEPOSIT && op.OperationType != model.WITHDRAW {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation type"})
		return
	}

	if op.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return
	}

	if err := h.wallet.ProcessOperation(context.Background(), op); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}


func (h *WalletHandler) GetWalletBalance(c *gin.Context) {
	walletIDStr := c.Param("walletId")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet ID"})
		return
	}

	balance, err := h.wallet.GetWalletBalance(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}