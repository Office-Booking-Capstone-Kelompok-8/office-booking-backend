package impl

import (
	"context"
	"office-booking-backend/internal/buildings/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"

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

func (b *BuildingRepositoryImpl) GetAllBuildings(ctx context.Context, q string, limit int, offset int) (*entity.Building, int64, error) {
	buildings := &entity.Building{}
	var count int64

	query := b.db.WithContext(ctx).Joins("Pictures").Model(&entity.Building{})
	if q != "" {
		query = query.Where("`Pictures`.`name` LIKE ?", "%"+q+"%")
	}

	err := query.Limit(limit).Offset(offset).Find(buildings).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return buildings, count, nil
}

func (b *BuildingRepositoryImpl) GetBuildingDetailByID(ctx context.Context, id string) (*entity.Building, error) {
	building := &entity.Building{}
	err := b.db.WithContext(ctx).Model(&entity.Building{}).Joins("Pictures").Where("id = ?", id).First(building).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err2.ErrUserNotFound
		}

		return nil, err
	}

	return building, nil
}
