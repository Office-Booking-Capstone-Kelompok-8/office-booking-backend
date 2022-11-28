package impl

import (
	"context"
	"office-booking-backend/internal/building/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"time"

	"gorm.io/gorm"
)

type BuildingRepositoryImpl struct {
	db *gorm.DB
}

func NewBuildingRepositoryImpl(db *gorm.DB) repository.BuildingRepository {
	return &BuildingRepositoryImpl{
		db: db,
	}
}

func (b *BuildingRepositoryImpl) GetAllBuildings(ctx context.Context, q string, cityID int, districtID int, startDate time.Time, endDate time.Time, limit int, offset int) (*entity.Buildings, int64, error) {
	buildings := &entity.Buildings{}
	var count int64

	query := b.db.WithContext(ctx).
		Preload("Pictures", func(db *gorm.DB) *gorm.DB {
			return db.Select("building_id, thumbnail_url")
		}).
		Joins("District").
		Joins("City").
		Model(&entity.Building{})
	if q != "" {
		query = query.Where("`buildings`.`name` LIKE ?", "%"+q+"%")
	}

	if cityID != 0 {
		query = query.Where("`City`.`id` = ?", cityID)
	}

	if districtID != 0 {
		query = query.Where("`District`.`id` = ?", districtID)
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("NOT EXISTS (SELECT * FROM `reservations` WHERE `reservations`.`building_id` = `buildings`.`id` AND `reservations`.`start_date` <= ? AND `reservations`.`end_date` >= ?)", endDate, startDate)
	}

	err := query.Limit(limit).Offset(offset).Find(buildings).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return buildings, count, nil
}

func (b *BuildingRepositoryImpl) GetBuildingDetailByID(ctx context.Context, id string) (*entity.Building, error) {
	building := &entity.Building{}
	err := b.db.WithContext(ctx).
		Preload("Pictures").
		Preload("Facilities", func(db *gorm.DB) *gorm.DB {
			return db.Joins("Category")
		}).
		Joins("District").
		Joins("City").
		Model(&entity.Building{}).
		Where("`buildings`.`id` = ?", id).
		First(building).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err2.ErrBuildingNotFound
		}

		return nil, err
	}

	return building, nil
}
