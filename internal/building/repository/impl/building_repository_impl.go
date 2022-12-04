package impl

import (
	"context"
	"fmt"
	"office-booking-backend/internal/building/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/utils/ptr"
	"strings"
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

func (b *BuildingRepositoryImpl) GetAllBuildings(ctx context.Context, q string, cityID int, districtID int, startDate time.Time, endDate time.Time, limit int, offset int, isPublishedOnly bool) (*entity.Buildings, int64, error) {
	buildings := &entity.Buildings{}
	var count int64

	query := b.db.WithContext(ctx).
		Preload("Pictures", "`pictures`.`index` = 0").
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

	if isPublishedOnly {
		query = query.Where("`buildings`.`is_published` = ?", true)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Find(buildings).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return buildings, count, nil
}

func (b *BuildingRepositoryImpl) GetBuildingDetailByID(ctx context.Context, id string, isPublishedOnly bool) (*entity.Building, error) {
	building := &entity.Building{}

	query := b.db.WithContext(ctx).
		Preload("Pictures", func(db *gorm.DB) *gorm.DB {
			return db.Order("`pictures`.`index` ASC")
		}).
		Preload("Facilities", func(db *gorm.DB) *gorm.DB {
			return db.Joins("Category")
		}).
		Joins("District").
		Joins("City").
		Model(&entity.Building{}).
		Where("`buildings`.`id` = ?", id)

	if isPublishedOnly {
		query = query.Where("`buildings`.`is_published` = ?", true)
	}

	err := query.First(building).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err2.ErrBuildingNotFound
		}

		return nil, err
	}

	return building, nil
}

func (b *BuildingRepositoryImpl) GetFacilityCategories(ctx context.Context) (*entity.Categories, error) {
	categories := &entity.Categories{}
	err := b.db.WithContext(ctx).
		Model(&entity.Category{}).
		Find(categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (b *BuildingRepositoryImpl) CreateBuilding(ctx context.Context, building *entity.Building) error {
	err := b.db.WithContext(ctx).
		Model(&entity.Building{}).
		Create(building).Error
	if err != nil {
		return err
	}

	return nil
}

func (b *BuildingRepositoryImpl) UpdateBuildingByID(ctx context.Context, building *entity.Building) error {
	// Standard version
	// res := b.db.WithContext(ctx).
	// 	Model(&entity.Building{}).
	// 	Where("id = ?", building.ID).
	// 	Updates(building)
	// if res.Error != nil {
	// 	return res.Error
	// }
	//
	// if res.RowsAffected == 0 {
	// 	return err2.ErrBuildingNotFound
	// }
	// return nil

	return b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).
			Model(&entity.Building{}).
			Where("id = ?", building.ID).
			Updates(building).Error
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "CONSTRAINT `fk_buildings_city`"):
				return err2.ErrInavalidCityID
			case strings.Contains(err.Error(), "CONSTRAINT `fk_buildings_district`"):
				return err2.ErrInvalidDistrictID
			default:
				return err
			}
		}

		for _, picture := range building.Pictures {
			fmt.Println(picture)
			err := tx.WithContext(ctx).
				Model(&entity.Picture{}).
				Where("id = ?", picture.ID).
				Where("building_id = ?", building.ID).
				Updates(picture).Error
			if err != nil {
				return err
			}
		}

		for _, facility := range building.Facilities {
			err := tx.WithContext(ctx).
				Model(&entity.Facility{}).
				Where("id = ?", facility.ID).
				Where("building_id = ?", building.ID).
				Updates(facility).Error
			if err != nil {
				if strings.Contains(err.Error(), "CONSTRAINT `fk_facilities_category` FOREIGN KEY (`category_id`)") {
					return err2.ErrInvalidCategoryID
				}
				return err
			}
		}

		return nil
	})
}

func (b *BuildingRepositoryImpl) IsBuildingExist(ctx context.Context, buildingId string) (bool, error) {
	var count int64
	err := b.db.WithContext(ctx).
		Model(&entity.Building{}).
		Where("id = ?", buildingId).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (b *BuildingRepositoryImpl) IsBuildingPublished(ctx context.Context, buildingID string) (bool, error) {
	var building entity.Building
	err := b.db.WithContext(ctx).
		Model(&entity.Building{}).
		Where("id = ?", buildingID).
		First(&building).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, err2.ErrBuildingNotFound
		}

		return false, err
	}

	return *building.IsPublished, nil
}

func (b *BuildingRepositoryImpl) CountBuildingPicturesByID(ctx context.Context, buildingId string) (int64, error) {
	var count int64
	err := b.db.WithContext(ctx).
		Model(&entity.Picture{}).
		Where("building_id = ?", buildingId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (b *BuildingRepositoryImpl) AddPicture(ctx context.Context, picture *entity.Picture) error {
	err := b.db.WithContext(ctx).
		Model(&entity.Picture{}).
		Create(picture).Error
	if err != nil {
		return err
	}

	return nil
}

func (b *BuildingRepositoryImpl) AddFacility(ctx context.Context, facility *entity.Facilities) error {
	err := b.db.WithContext(ctx).
		Model(&entity.Facility{}).
		Create(facility).Error
	if err != nil {
		if strings.Contains(err.Error(), "CONSTRAINT `fk_facilities_category` FOREIGN KEY (`category_id`)") {
			return err2.ErrInvalidCategoryID
		}
		return err
	}

	return nil
}

func (b *BuildingRepositoryImpl) DeleteBuildingPicturesByID(ctx context.Context, buildingID string, pictureID string) error {
	// res := b.db.WithContext(ctx).
	// 	Unscoped().
	// 	Model(&entity.Picture{}).
	// 	Where("id = ?", pictureID).
	// 	Where("building_id = ?", buildingID).
	// 	Delete(&entity.Picture{})
	// if res.Error != nil {
	// 	return res.Error
	// }

	// if res.RowsAffected == 0 {
	// 	return err2.ErrPictureNotFound
	// }

	err := b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get picture that will be deleted
		var picture entity.Picture
		err := tx.WithContext(ctx).
			Model(&entity.Picture{}).
			Where("id = ?", pictureID).
			Where("building_id = ?", buildingID).
			First(&picture).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return err2.ErrPictureNotFound
			}
			return err
		}

		err = tx.WithContext(ctx).
			Unscoped().
			Model(&entity.Picture{}).
			Where("id = ?", pictureID).
			Where("building_id = ?", buildingID).
			Delete(&entity.Picture{}).Error
		if err != nil {
			return err
		}

		if *picture.Index == 0 {
			return changeNextPictureToIndexZero(ctx, tx, buildingID)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func changeNextPictureToIndexZero(ctx context.Context, tx *gorm.DB, buildingID string) error {
	nextPicture := entity.Picture{}
	err := tx.WithContext(ctx).
		Model(&entity.Picture{}).
		Where("building_id = ?", buildingID).
		Order("`pictures`.`index` ASC").
		First(&nextPicture).Error
	fmt.Println("nextPicture", nextPicture)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		// change building isPublished to false
		return tx.WithContext(ctx).
			Model(&entity.Building{}).
			Where("id = ?", buildingID).
			Updates(entity.Building{IsPublished: ptr.Bool(true)}).Error
	} else {
		fmt.Println("nextPicture", nextPicture.ID, "changed to index 0")
		// change next picture index to 0
		return tx.WithContext(ctx).
			Model(&entity.Picture{}).
			Select("index").
			Where("id = ?", nextPicture.ID).
			Update("index", 0).Error
	}
}

func (b *BuildingRepositoryImpl) DeleteBuildingFacilityByID(ctx context.Context, buildingID string, facilityID int) error {
	res := b.db.WithContext(ctx).
		Model(&entity.Facility{}).
		Where("id = ?", facilityID).
		Where("building_id = ?", buildingID).
		Delete(&entity.Facility{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrFacilityNotFound
	}

	return nil
}

func (b *BuildingRepositoryImpl) DeleteBuildingByID(ctx context.Context, buildingID string) error {
	res := b.db.WithContext(ctx).
		Where("id = ?", buildingID).
		Delete(&entity.Building{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrBuildingNotFound
	}

	return nil
}
