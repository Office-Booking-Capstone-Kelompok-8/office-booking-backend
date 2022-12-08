package dto

import "office-booking-backend/pkg/entity"

type PaymentResponse struct {
	ID            uint   `json:"id"`
	BankIcon      string `json:"bankIcon"`
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
	Description   string `json:"description"`
}

func NewPaymentResponse(payment *entity.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:            payment.Model.ID,
		BankIcon:      payment.Bank.Icon,
		BankName:      payment.Bank.Name,
		AccountNumber: payment.AccountNumber,
		AccountName:   payment.AccountName,
		Description:   payment.Description,
	}
}

type PaymentsResponse []PaymentResponse

func NewPaymentsResponse(payments *entity.Payments) *PaymentsResponse {
	var paymentsResponse PaymentsResponse
	for _, payment := range *payments {
		paymentsResponse = append(paymentsResponse, *NewPaymentResponse(&payment))
	}
	return &paymentsResponse
}

type BankResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

func NewBankResponse(bank *entity.Bank) *BankResponse {
	return &BankResponse{
		ID:   bank.ID,
		Name: bank.Name,
		Icon: bank.Icon,
	}
}

type BanksResponse []BankResponse

func NewBanksResponse(banks *entity.Banks) *BanksResponse {
	var banksResponse BanksResponse
	for _, bank := range *banks {
		banksResponse = append(banksResponse, *NewBankResponse(&bank))
	}
	return &banksResponse
}
