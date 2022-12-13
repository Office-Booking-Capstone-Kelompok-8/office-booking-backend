package impl

import (
	"fmt"
	"office-booking-backend/internal/reservation/dto"
	"office-booking-backend/internal/reservation/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type ReservationRepositoryImpl struct {
	db *gorm.DB
}

func NewReservationRepositoryImpl(db *gorm.DB) repository.ReservationRepository {
	return &ReservationRepositoryImpl{
		db: db,
	}
}

func (r *ReservationRepositoryImpl) CountUserActiveReservations(ctx context.Context, userID string) (int64, error) {
	var count int64
	// Count user active reservations with status id not 2 (rejected), 3 (canceled) or 5 (completed) and not expired
	err := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("user_id = ? AND status_id NOT IN (2, 3, 5) AND end_date > ?", userID, time.Now()).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReservationRepositoryImpl) CountBuildingActiveReservations(ctx context.Context, buildingID string) (int64, error) {
	var count int64
	// Count building active reservations with status id not 2 (rejected), 3 (canceled) or 5 (completed) and not expired
	err := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("building_id = ? AND status_id NOT IN (2, 3, 5) AND end_date > ?", buildingID, time.Now()).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReservationRepositoryImpl) CountUserReservation(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReservationRepositoryImpl) CountReservation(ctx context.Context, filter *dto.ReservationQueryParam) (int64, error) {
	db, err := r.db.DB()
	if err != nil {
		return 0, err
	}

	query := sq.Select("COUNT(*)").
		From("reservations r").
		Join("users u ON u.id = r.user_id").
		Join("user_details ud ON u.id = ud.user_id").
		Where("r.deleted_at IS NULL")

	if filter.UserID != "" {
		query = query.Where(sq.Eq{"r.user_id": filter.UserID})
	}

	if filter.UserName != "" {
		query = query.Where("ud.name LIKE ?", "%"+filter.UserName+"%")
	}

	if filter.BuildingID != "" {
		query = query.Where(sq.Eq{"r.building_id": filter.BuildingID})
	}

	if filter.StatusID != 0 {
		query = query.Where(sq.Eq{"r.status_id": filter.StatusID})
	}

	if !filter.StartDate.ToTime().IsZero() {
		query = query.Where(sq.GtOrEq{"r.start_date": filter.StartDate.ToTime()})
	}

	if !filter.EndDate.ToTime().IsZero() {
		query = query.Where(sq.LtOrEq{"r.end_date": filter.EndDate.ToTime()})
	}

	rows, err := query.RunWith(db).QueryContext(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		err = rows.Close()
	}()

	var count int64
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (r *ReservationRepositoryImpl) IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time, excludedReservationID ...string) (bool, error) {
	var count int64
	// Count building active reservations with status id not 2 (rejected), 3 (canceled) or 5 (completed)
	// and not in the same time range as the new reservation
	query := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("building_id = ? AND status_id NOT IN (2, 3, 5)", buildingID).
		Where("start_date <= ? AND end_date >= ?", end, start)

	for _, id := range excludedReservationID {
		query = query.Where("id != ?", id)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *ReservationRepositoryImpl) GetReservations(ctx context.Context, filter *dto.ReservationQueryParam) (*entity.Reservations, error) {
	db, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := sq.Select("r.id, r.company_name, r.building_id, r.start_date, r.end_date,  r.status_id, r.created_at, r.updated_at, s.id, s.message, b.id, b.name, p.thumbnail_url, c.name, u.id, u.email, ud.name, pp.url").
		From("reservations r").
		Join("statuses s ON r.status_id = s.id").
		Join("buildings b ON r.building_id = b.id").
		Join("pictures p ON p.id = (SELECT p1.id FROM pictures p1 WHERE b.id = p1.building_id AND p1.index = 0 LIMIT 1)").
		Join("cities c ON c.id = b.city_id").
		Join("users u ON u.id = r.user_id").
		Join("user_details ud ON u.id = ud.user_id").
		LeftJoin("profile_pictures pp ON ud.picture_id = pp.id").
		Where("r.deleted_at IS NULL")

	if filter.UserID != "" {
		query = query.Where(sq.Eq{"r.user_id": filter.UserID})
	}

	if filter.UserName != "" {
		query = query.Where("ud.name LIKE ?", "%"+filter.UserName+"%")
	}

	if filter.BuildingID != "" {
		query = query.Where(sq.Eq{"r.building_id": filter.BuildingID})
	}

	if filter.StatusID != 0 {
		query = query.Where(sq.Eq{"r.status_id": filter.StatusID})
	}

	if !filter.StartDate.ToTime().IsZero() {
		query = query.Where(sq.GtOrEq{"r.start_date": filter.StartDate.ToTime()})
	}

	if !filter.EndDate.ToTime().IsZero() {
		query = query.Where(sq.LtOrEq{"r.end_date": filter.EndDate.ToTime()})
	}

	if !filter.CreatedStart.ToTime().IsZero() && !filter.CreatedEnd.ToTime().IsZero() {
		query = query.Where("DATE(r.created_at) BETWEEN DATE(?) AND DATE(?)", filter.CreatedStart.ToTime(), filter.CreatedEnd.ToTime())
	}

	switch filter.SortBy {
	case "":
	case "created_at":
		fmt.Println("Not this one")
		filter.SortBy = "r.created_at"
	case "start_date":
		filter.SortBy = "r.start_date"
	case "end_date":
		filter.SortBy = "r.end_date"
	case "building_name":
		filter.SortBy = "b.name"
	case "user_name":
		filter.SortBy = "ud.name"
	default:
		filter.SortBy = ""
	}

	rows, err := query.
		Offset(uint64(filter.Offset)).
		Limit(uint64(filter.Limit)).
		OrderBy(fmt.Sprintf("%s %s", filter.SortBy, filter.SortOrder)).
		RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	var reservations entity.Reservations
	for rows.Next() {
		var reservation entity.Reservation
		NullAbleProfilePicture := &entity.NullAbleProfilePicture{}
		reservation.Building.Pictures = append(reservation.Building.Pictures, entity.Picture{})
		err = rows.Scan(&reservation.ID, &reservation.CompanyName, &reservation.BuildingID, &reservation.StartDate, &reservation.EndDate, &reservation.StatusID, &reservation.CreatedAt, &reservation.UpdatedAt,
			&reservation.Status.ID, &reservation.Status.Message,
			&reservation.Building.ID, &reservation.Building.Name, &reservation.Building.Pictures[0].ThumbnailUrl, &reservation.Building.City.Name,
			&reservation.User.ID, &reservation.User.Email, &reservation.User.Detail.Name, &NullAbleProfilePicture.Url)
		if err != nil {
			return nil, err
		}

		reservation.User.Detail.Picture = NullAbleProfilePicture.ConvertToProfilePicture()

		reservations = append(reservations, reservation)
	}

	return &reservations, nil
}

func (r *ReservationRepositoryImpl) GetUserReservations(ctx context.Context, userID string, offset int, limit int) (*entity.Reservations, error) {
	db, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	rows, err := sq.Select("r.id, r.company_name, r.building_id, r.start_date, r.end_date,  r.status_id, r.created_at, r.updated_at, s.id, s.message, b.id, b.name, p.thumbnail_url, c.name").
		From("reservations r").
		Join("statuses s ON r.status_id = s.id").
		Join("buildings b ON r.building_id = b.id").
		Join("pictures p ON p.id = (SELECT p1.id FROM pictures p1 WHERE b.id = p1.building_id AND p1.index = 0 LIMIT 1)").
		Join("cities c ON c.id = b.city_id").
		Where("r.deleted_at IS NULL").
		Where(sq.Eq{"r.user_id": userID}).
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	var reservations entity.Reservations
	for rows.Next() {
		var reservation entity.Reservation
		reservation.Building.Pictures = append(reservation.Building.Pictures, entity.Picture{})
		err = rows.Scan(&reservation.ID, &reservation.CompanyName, &reservation.BuildingID, &reservation.StartDate, &reservation.EndDate, &reservation.StatusID, &reservation.CreatedAt, &reservation.UpdatedAt,
			&reservation.Status.ID, &reservation.Status.Message,
			&reservation.Building.ID, &reservation.Building.Name, &reservation.Building.Pictures[0].ThumbnailUrl, &reservation.Building.City.Name)
		if err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	return &reservations, nil
}

func (r *ReservationRepositoryImpl) GetReservationByID(ctx context.Context, reservationID string) (*entity.Reservation, error) {
	db, err := r.db.DB()
	if err != nil {
		return nil, err
	}
	rows, err := sq.Select("r.id, r.company_name, r.building_id, r.start_date, r.end_date, r.user_id, r.status_id, r.message, r.created_at, r.updated_at, s.id, s.message, b.id, b.name, b.address, p.thumbnail_url, c.name, d.name, u.id, u.email, ud.name, pp.url").
		From("reservations r").
		Join("statuses s ON s.id = r.status_id").
		Join("buildings b ON b.id = r.building_id").
		Join("pictures p ON p.id = (SELECT p1.id FROM pictures p1 WHERE b.id = p1.building_id AND p1.index = 0 LIMIT 1)").
		Join("cities c ON c.id = b.city_id").
		Join("districts d ON d.id = b.district_id").
		Join("users u ON u.id = r.user_id").
		Join("user_details ud ON u.id = ud.user_id").
		LeftJoin("profile_pictures pp ON ud.picture_id = pp.id").
		Where("r.deleted_at IS NULL AND r.id = ?", reservationID).RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	if !rows.Next() {
		return nil, err2.ErrReservationNotFound
	}

	var reservation entity.Reservation
	NullAbleProfilePicture := &entity.NullAbleProfilePicture{}
	reservation.Building.Pictures = append(reservation.Building.Pictures, entity.Picture{})
	err = rows.Scan(&reservation.ID, &reservation.CompanyName, &reservation.BuildingID, &reservation.StartDate, &reservation.EndDate,
		&reservation.UserID, &reservation.StatusID, &reservation.Message, &reservation.CreatedAt, &reservation.UpdatedAt, &reservation.Status.ID, &reservation.Status.Message,
		&reservation.Building.ID, &reservation.Building.Name, &reservation.Building.Address, &reservation.Building.Pictures[0].ThumbnailUrl,
		&reservation.Building.City.Name, &reservation.Building.District.Name, &reservation.User.ID, &reservation.User.Email, &reservation.User.Detail.Name,
		&NullAbleProfilePicture.Url)
	if err != nil {
		return nil, err
	}

	reservation.User.Detail.Picture = NullAbleProfilePicture.ConvertToProfilePicture()

	return &reservation, nil
}

func (r *ReservationRepositoryImpl) GetUserReservationByID(ctx context.Context, reservationID string, userID string) (*entity.Reservation, error) {
	db, err := r.db.DB()
	if err != nil {
		return nil, err
	}
	rows, err := sq.Select("r.id, r.company_name, r.building_id, r.start_date, r.end_date, r.user_id, r.status_id, r.message, r.created_at, r.updated_at, s.id, s.message, b.id, b.name, b.address, p.thumbnail_url, c.name, d.name").
		From("reservations r").
		Join("statuses s ON s.id = r.status_id").
		Join("buildings b ON b.id = r.building_id").
		Join("pictures p ON p.id = (SELECT p1.id FROM pictures p1 WHERE b.id = p1.building_id AND p1.index = 0 LIMIT 1)").
		Join("cities c ON c.id = b.city_id").
		Join("districts d ON d.id = b.district_id").
		Where("r.deleted_at IS NULL").
		Where("r.id = ?", reservationID).
		Where("r.user_id = ?", userID).
		RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	if !rows.Next() {
		return nil, err2.ErrReservationNotFound
	}

	var reservation entity.Reservation
	reservation.Building.Pictures = append(reservation.Building.Pictures, entity.Picture{})
	err = rows.Scan(&reservation.ID, &reservation.CompanyName, &reservation.BuildingID, &reservation.StartDate, &reservation.EndDate,
		&reservation.UserID, &reservation.StatusID, &reservation.Message, &reservation.CreatedAt, &reservation.UpdatedAt, &reservation.Status.ID, &reservation.Status.Message,
		&reservation.Building.ID, &reservation.Building.Name, &reservation.Building.Address, &reservation.Building.Pictures[0].ThumbnailUrl,
		&reservation.Building.City.Name, &reservation.Building.District.Name)
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *ReservationRepositoryImpl) GetReservationTotal(ctx context.Context) (*entity.StatusesStat, error) {
	rows, err := r.db.WithContext(ctx).
		Table("statuses").
		Select("statuses.id, statuses.message, COUNT(reservations.id) AS total").
		Joins("left join reservations on reservations.status_id = statuses.id").
		Group("statuses.id").
		Rows()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	stat := entity.StatusesStat{}
	for rows.Next() {
		var status entity.StatusStat
		err = rows.Scan(&status.StatusID, &status.StatusName, &status.Total)
		if err != nil {
			return nil, err
		}

		stat = append(stat, status)
	}

	return &stat, nil
}

func (r *ReservationRepositoryImpl) GetReservationCount(ctx context.Context) (*entity.TimeframeStat, error) {
	rows, err := r.db.WithContext(ctx).
		Table(
			"(?) AS today, (?) AS thisWeek, (?) AS thisMonth, (?) AS thisYear",
			r.db.Table("reservations").Select("count(*)").Where("DATE(created_at) = DATE(?)", time.Now().Format("2006-01-02")),
			r.db.Table("reservations").Select("count(*)").Where("YEARWEEK(created_at) = YEARWEEK(?)", time.Now().Format("2006-01-02")),
			r.db.Table("reservations").Select("count(*)").Where("MONTH(created_at) = MONTH(?)", time.Now().Format("2006-01-02")),
			r.db.Table("reservations").Select("count(*)").Where("YEAR(created_at) = YEAR(?)", time.Now().Format("2006-01-02")),
		).Rows()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	stat := new(entity.TimeframeStat)
	for rows.Next() {
		err = rows.Scan(&stat.Day, &stat.Week, &stat.Month, &stat.Year)
		if err != nil {
			return nil, err
		}
	}

	return stat, nil
}

func (r *ReservationRepositoryImpl) AddBuildingReservation(ctx context.Context, reservation *entity.Reservation) error {
	err := r.db.WithContext(ctx).Create(reservation).Error
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "CONSTRAINT `fk_reservations_building`"):
			return err2.ErrBuildingNotFound
		case strings.Contains(err.Error(), "CONSTRAINT `fk_reservations_user`"):
			return err2.ErrInvalidUserID
		default:
			return err
		}
	}

	return nil
}

func (r *ReservationRepositoryImpl) UpdateReservation(ctx context.Context, reservation *entity.Reservation) error {
	res := r.db.WithContext(ctx).
		Model(entity.Reservation{}).
		Where("id = ?", reservation.ID).
		Updates(reservation)
	if res.Error != nil {
		switch {
		case strings.Contains(res.Error.Error(), "CONSTRAINT `fk_reservations_building`"):
			return err2.ErrBuildingNotFound
		case strings.Contains(res.Error.Error(), "CONSTRAINT `fk_reservations_user`"):
			return err2.ErrUserNotFound
		case strings.Contains(res.Error.Error(), "CONSTRAINT `fk_reservations_status`"):
			return err2.ErrInvalidStatus
		default:
			return res.Error
		}
	}

	if res.RowsAffected == 0 {
		return err2.ErrReservationNotFound
	}

	return nil
}

func (r *ReservationRepositoryImpl) DeleteReservationByID(ctx context.Context, reservationID string) error {
	res := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("id = ?", reservationID).
		Delete(&entity.Reservation{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrReservationNotFound
	}

	return nil
}

// get review for reservation
func (r *ReservationRepositoryImpl) GetReservationReviews(ctx context.Context) (*entity.Reviews, error) {
	review := new(entity.Reviews)
	err := r.db.WithContext(ctx).
		Model(&entity.Review{}).
		Joins("Reservation").
		Find(review).Error
	if err != nil {
		return nil, err
	}

	return review, nil
}

// add review after reservation status is completed
func (r *ReservationRepositoryImpl) AddReservationReviews(ctx context.Context, review *entity.Review) error {
	err := r.db.WithContext(ctx).Create(review).Error
	if err != nil {
		return err
	}

	return nil
}

// delete review for admin
func (r *ReservationRepositoryImpl) DeleteAdminReservationReviews(ctx context.Context, reservationID string) error {
	res := r.db.WithContext(ctx).
		Model(&entity.Review{}).
		Where("id = ?", reservationID).
		Delete(&entity.Review{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrReservationNotFound
	}

	return nil
}
