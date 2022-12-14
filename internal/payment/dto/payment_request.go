package dto

import (
	"office-booking-backend/pkg/entity"

	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

type CreatePaymentRequest struct {
	BankID        int    `json:"bankId" validate:"required,gte=1"`
	AccountNumber string `json:"accountNumber" validate:"required"`
	AccountName   string `json:"accountName" validate:"required"`
	Description   string `json:"description" validate:"omitempty,max=255"`
}

func (c *CreatePaymentRequest) ToEntity() *entity.Payment {
	return &entity.Payment{
		BankID:        c.BankID,
		AccountName:   c.AccountName,
		AccountNumber: c.AccountNumber,
		Description:   c.Description,
	}
}

type UpdatePaymentRequest struct {
	BankID        int    `json:"bankId" validate:"omitempty,gte=1"`
	AccountName   string `json:"accountName" validate:"omitempty"`
	AccountNumber string `json:"accountNumber" validate:"omitempty"`
	Description   string `json:"description" validate:"omitempty,max=255"`
}

func (u *UpdatePaymentRequest) ToEntity() *entity.Payment {
	return &entity.Payment{
		BankID:        u.BankID,
		AccountName:   u.AccountName,
		AccountNumber: u.AccountNumber,
		Description:   u.Description,
	}
}

type CreateReservationPaymentRequest struct {
	MethodID int `form:"methodId" validate:"required"`
}

func (c *CreateReservationPaymentRequest) ToEntity(proof *uploader.UploadResult) *entity.Transaction {
	return &entity.Transaction{
		PaymentID: uint(c.MethodID),
		ProofID:   proof.FileId,
		Proof: entity.PaymentProof{
			ID:  proof.FileId,
			URL: proof.Url,
		},
	}
}
