package dto

import (
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/entity"
)

type PaymentMethodResponse struct {
	ID            uint   `json:"id"`
	BankIcon      string `json:"bankIcon"`
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
	Description   string `json:"description"`
}

func NewPaymentMethodResponse(payment *entity.Payment) *PaymentMethodResponse {
	return &PaymentMethodResponse{
		ID:            payment.Model.ID,
		BankIcon:      payment.Bank.Icon,
		BankName:      payment.Bank.Name,
		AccountNumber: payment.AccountNumber,
		AccountName:   payment.AccountName,
		Description:   payment.Description,
	}
}

type PaymentMethodsResponse []PaymentMethodResponse

func NewPaymentMethodsResponse(payments *entity.Payments) *PaymentMethodsResponse {
	var paymentsResponse PaymentMethodsResponse
	for _, payment := range *payments {
		paymentsResponse = append(paymentsResponse, *NewPaymentMethodResponse(&payment))
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

type PaymentDetailResponse struct {
	ID        string                     `json:"id"`
	Ammount   int                        `json:"ammount"`
	StartDate string                     `json:"startDate"`
	EndDate   string                     `json:"endDate"`
	Method    BriefPaymentMethodResponse `json:"method"`
	Proof     string                     `json:"proof"`
	CreatedAt string                     `json:"createdAt"`
	UpdatedAt string                     `json:"updatedAt"`
}

func NewPaymentDetailResponse(payment *entity.Transaction) *PaymentDetailResponse {
	return &PaymentDetailResponse{
		ID:        payment.ID,
		Ammount:   payment.Reservation.Amount,
		StartDate: payment.Reservation.StartDate.Format(config.DATE_RESPONSE_FORMAT),
		EndDate:   payment.Reservation.EndDate.Format(config.DATE_RESPONSE_FORMAT),
		Method: BriefPaymentMethodResponse{
			ID:            payment.Payment.ID,
			BankIcon:      payment.Payment.Bank.Icon,
			BankName:      payment.Payment.Bank.Name,
			AccountNumber: payment.Payment.AccountNumber,
			AccountName:   payment.Payment.AccountName,
		},
		Proof:     payment.Proof.URL,
		CreatedAt: payment.CreatedAt.Format(config.DATE_RESPONSE_FORMAT),
		UpdatedAt: payment.UpdatedAt.Format(config.DATE_RESPONSE_FORMAT),
	}
}

type BriefPaymentMethodResponse struct {
	ID            uint   `json:"id"`
	BankIcon      string `json:"bankIcon"`
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
}
