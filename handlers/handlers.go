package handler

import (
	"net/http"

	"Go_wallet/model"
	"Go_wallet/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	service service.WalletService
}

func NewWalletHandler(service service.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) HandleWalletOperation(c *gin.Context) {

	var op model.WalletOperation
	
	if err := c.ShouldBindJSON(&op); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if op.OperationType != model.DEPOSIT && op.OperationType != model.WITHDRAW {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation type"})
		return
	}

	if op.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return
	}

	if err := h.service.ProcessOperation(c.Request.Context(), op); err != nil {
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


	balance, err := h.service.GetWalletBalance(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}