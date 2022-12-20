package impl

import (
	"context"
	"io"
	"log"
	"office-booking-backend/internal/payment/dto"
	"office-booking-backend/internal/payment/repository"
	"office-booking-backend/internal/payment/service"
	reservationRepo "office-booking-backend/internal/reservation/repository"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/utils/imagekit"
	"time"

	"github.com/google/uuid"
)

type PaymentServiceImpl struct {
	repo            repository.PaymentRepository
	reservationRepo reservationRepo.ReservationRepository
	imgKitService   imagekit.ImgKitService
}

func NewPaymentServiceImpl(repo repository.PaymentRepository, reservationRepo reservationRepo.ReservationRepository, imgKitService imagekit.ImgKitService) service.PaymentService {
	return &PaymentServiceImpl{
		repo:            repo,
		reservationRepo: reservationRepo,
		imgKitService:   imgKitService,
	}
}

func (p *PaymentServiceImpl) GetBanks(ctx context.Context) (*dto.BanksResponse, error) {
	banks, err := p.repo.GetAllBank(ctx)
	if err != nil {
		log.Println("error when getting all banks: ", err)
		return nil, err
	}

	return dto.NewBanksResponse(banks), nil
}

func (p *PaymentServiceImpl) GetAllPaymentMethod(ctx context.Context) (*dto.PaymentMethodsResponse, error) {
	payments, err := p.repo.GetAllPaymentMethod(ctx)
	if err != nil {
		log.Println("error when getting all payments: ", err)
		return nil, err
	}

	return dto.NewPaymentMethodsResponse(payments), nil
}

func (p *PaymentServiceImpl) GetPaymentMethodByID(ctx context.Context, paymentID int) (*dto.PaymentMethodResponse, error) {
	payment, err := p.repo.GetPaymentMethodByID(ctx, paymentID)
	if err != nil {
		log.Println("error when get payment by id: ", err)
		return nil, err
	}

	return dto.NewPaymentMethodResponse(payment), nil
}

func (p *PaymentServiceImpl) GetReservationPaymentByID(ctx context.Context, reservationID string) (*dto.PaymentDetailResponse, error) {
	payment, err := p.repo.GetReservationPaymentByID(ctx, reservationID, "")
	if err != nil {
		log.Println("error when get reservation payment by id: ", err)
		return nil, err
	}

	return dto.NewPaymentDetailResponse(payment), nil
}

func (p *PaymentServiceImpl) GetUserReservationPaymentByID(ctx context.Context, reservationID string, userID string) (*dto.PaymentDetailResponse, error) {
	payment, err := p.repo.GetReservationPaymentByID(ctx, reservationID, userID)
	if err != nil {
		log.Println("error when get reservation payment by id: ", err)
		return nil, err
	}

	return dto.NewPaymentDetailResponse(payment), nil
}

func (p *PaymentServiceImpl) CreatePaymentMethod(ctx context.Context, payment *dto.CreatePaymentRequest) (uint, error) {
	paymentEntity := payment.ToEntity()
	err := p.repo.CreatePaymentMethod(ctx, paymentEntity)
	if err != nil {
		log.Println("error when create payment: ", err)
		return 0, err
	}

	return paymentEntity.ID, nil
}

func (p *PaymentServiceImpl) AddPaymentProof(ctx context.Context, reservationID string, payment *dto.CreateReservationPaymentRequest, file io.Reader) error {
	reservation, err := p.reservationRepo.GetReservationByID(ctx, reservationID)
	if err != nil {
		log.Println("error when get reservation by id: ", err)
		return err
	}

	if reservation.StatusID == 5 {
		return err2.ErrReservationAlreadyPaid
	}

	if reservation.StatusID == 3 || reservation.ExpiredAt.Before(time.Now()) {
		return err2.ErrPaymentAlreadyExpired
	}

	if reservation.StatusID != 4 {
		return err2.ErrReservationNotAwaitingPayment
	}

	filename := uuid.New().String()
	uploadResponse, err := p.imgKitService.UploadFile(ctx, file, filename, "PaymentProof")
	if err != nil {
		log.Println("error when upload image to imagekit: ", err)
		return err2.ErrPictureServiceFailed
	}

	paymentEntity := payment.ToEntity(uploadResponse)
	paymentEntity.ReservationID = reservationID
	err = p.repo.CreateNewReservationPayment(ctx, paymentEntity)
	if err != nil {
		log.Println("error when create reservation payment: ", err)
		return err
	}

	return nil
}

func (p *PaymentServiceImpl) UpdatePaymentMethod(ctx context.Context, paymentID int, payment *dto.UpdatePaymentRequest) error {
	paymentEntity := payment.ToEntity()
	paymentEntity.ID = uint(paymentID)
	err := p.repo.UpdatePaymentMethod(ctx, paymentEntity)
	if err != nil {
		log.Println("error when update payment: ", err)
		return err
	}

	return nil
}

func (p *PaymentServiceImpl) DeletePaymentMethod(ctx context.Context, paymentID int) error {
	err := p.repo.DeletePaymentMethod(ctx, paymentID)
	if err != nil {
		log.Println("error when delete payment: ", err)
		return err
	}

	return nil
}
