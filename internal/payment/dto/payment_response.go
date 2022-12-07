package dto

import "office-booking-backend/pkg/entity"

type PaymentResponse struct {
	ID            uint   `json:"id"`
	BankIcon      string `json:"bankIcon"`
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	Description   string `json:"description"`
}

func NewPaymentResponse(payment *entity.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:            payment.Model.ID,
		BankIcon:      payment.Bank.Icon,
		BankName:      payment.Bank.Name,
		AccountNumber: payment.AccountNumber,
		Description:   payment.Description,
	}
}
