package impl

import (
	"context"
	"fmt"
	"office-booking-backend/internal/building/dto"
	"office-booking-backend/internal/building/repository"
	"office-booking-backend/pkg/custom"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BuildingRepositoryImpl struct {
	db *gorm.DB
}

func NewBuildingRepositoryImpl(db *gorm.DB) repository.BuildingRepository {
	return &BuildingRepositoryImpl{
		db: db,
	}
}

func (b *BuildingRepositoryImpl) GetAllBuildings(ctx context.Context, filter *dto.SearchBuildingQueryParam, isPublishedOnly bool) (*entity.Buildings, int64, error) {
	buildings := &entity.Buildings{}
	var count int64

	query := b.db.WithContext(ctx).
		Preload("Pictures", "`pictures`.`index` = 0").
		Joins("District").
		Joins("City").
		Model(&entity.Building{})
	if filter.BuildingName != "" {
		query = query.Where("`buildings`.`name` LIKE ?", "%"+filter.BuildingName+"%")
	}

	if filter.CityID != 0 {
		query = query.Where("`City`.`id` = ?", filter.CityID)
	}

	if filter.DistrictID != 0 {
		query = query.Where("`District`.`id` = ?", filter.DistrictID)
	}

	if !filter.StartDate.ToTime().IsZero() && !filter.EndDate.IsZero() {
		query = query.Where("NOT EXISTS (SELECT * FROM `reservations` WHERE `reservations`.`building_id` = `buildings`.`id` AND `reservations`.`start_date` <= ? AND `reservations`.`end_date` >= ?)", filter.EndDate, filter.StartDate.ToTime())
	}

	if filter.AnnualPriceMin != 0 {
		query = query.Where("`buildings`.`annual_price` >= ?", filter.AnnualPriceMin)
	}

	if filter.AnnualPriceMax != 0 {
		query = query.Where("`buildings`.`annual_price` <= ?", filter.AnnualPriceMax)
	}

	if filter.MonthlyPriceMin != 0 {
		query = query.Where("`buildings`.`monthly_price` >= ?", filter.MonthlyPriceMin)
	}

	if filter.MonthlyPriceMax != 0 {
		query = query.Where("`buildings`.`monthly_price` <= ?", filter.MonthlyPriceMax)
	}

	if filter.CapacityMin != 0 {
		query = query.Where("`buildings`.`capacity` >= ?", filter.CapacityMin)
	}

	if filter.CapacityMax != 0 {
		query = query.Where("`buildings`.`capacity` <= ?", filter.CapacityMax)
	}

	if filter.SortBy != "" {
		// POW(69.1 * (latitude - [startlat]), 2) +
		// POW(69.1 * ([startlng] - longitude) * COS(latitude / 57.3), 2))
		if filter.SortBy == "pinpoint" {
			query = query.Clauses(clause.OrderBy{
				Expression: clause.Expr{
					SQL:                fmt.Sprintf("POW(69.1 * (`buildings`.`latitude` - ?), 2) + POW(69.1 * (? - `buildings`.`longitude`) * COS(`buildings`.`latitude` / 57.3), 2) %s", filter.Order),
					Vars:               []interface{}{filter.Latitude, filter.Longitude},
					WithoutParentheses: true,
				},
			})
		} else {
			query = query.Order(fmt.Sprintf("`buildings`.`%s` %s", filter.SortBy, filter.Order))
		}
	}

	if isPublishedOnly {
		query = query.Where("`buildings`.`is_published` = ?", true)
	}

	err := query.
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(buildings).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return buildings, count, nil

	// db, err := b.db.DB()
	// if err != nil {
	// 	return nil, err
	// }

	// query := sq.Select("b.id, b.name, b.annual_price, b.monthly_price, b.owner, b.is_published, c.name, d.name").
	// 	From("buildings b").
	// 	LeftJoin("cities c ON b.city_id = c.id").
	// 	LeftJoin("districts d ON b.district_id = d.id").
	// 	LeftJoin("pictures p ON p.id = (SELECT p1.id FROM pictures p1 WHERE b.id = p1.building_id AND p1.index = 0 LIMIT 1)").
	// 	Where("b.deleted_at IS NULL").
	// 	Where("b.is_published = ?", isPublishedOnly)

	// if filter.BuildingName != "" {
	// 	query = query.Where("b.name LIKE ?", "%"+filter.BuildingName+"%")
	// }

	// if filter.CityID != 0 {
	// 	query = query.Where("b.city_id = ?", filter.CityID)
	// }

	// if filter.DistrictID != 0 {
	// 	query = query.Where("b.district_id = ?", filter.DistrictID)
	// }

	// if !filter.StartDate.IsZero() && !filter.EndDate.IsZero() {
	// 	query = query.Where("NOT EXISTS (SELECT * FROM reservations r WHERE r.building_id = b.id AND r.start_date <= ? AND r.end_date >= ?)", filter.EndDate, filter.StartDate)
	// }

	// if filter.AnnualPriceMin != 0 {
	// 	query = query.Where("b.annual_price >= ?", filter.AnnualPriceMin)
	// }

	// if filter.AnnualPriceMax != 0 {
	// 	query = query.Where("b.annual_price <= ?", filter.AnnualPriceMax)
	// }

	// if filter.MonthlyPriceMin != 0 {
	// 	query = query.Where("b.monthly_price >= ?", filter.MonthlyPriceMin)
	// }

	// if filter.MonthlyPriceMax != 0 {
	// 	query = query.Where("b.monthly_price <= ?", filter.MonthlyPriceMax)
	// }

	// if filter.CapacityMin != 0 {
	// 	query = query.Where("b.capacity >= ?", filter.CapacityMin)
	// }

	// if filter.CapacityMax != 0 {
	// 	query = query.Where("b.capacity <= ?", filter.CapacityMax)
	// }

	// if filter.SortBy != "" {
	// 	// POW(69.1 * (latitude - [startlat]), 2) +
	// 	// POW(69.1 * ([startlng] - longitude) * COS(latitude / 57.3), 2))
	// 	if filter.SortBy == "pinpoint" {
	// 		query = query.OrderByClause("POW(69.1 * (b.latitude - ?), 2) + POW(69.1 * (? - b.longitude) * COS(b.latitude / 57.3), 2) "+filter.Order, filter.Latitude, filter.Longitude, filter.Order)
	// 	} else {
	// 		query = query.OrderBy(filter.SortBy + " " + filter.Order)
	// 	}
	// }

	// rows, err := query.
	// 	Offset(uint64(filter.Offset)).
	// 	Limit(uint64(filter.Limit)).
	// 	RunWith(db).QueryContext(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// defer func() {
	// 	err = rows.Close()
	// }()

	// var buildings entity.Buildings
	// for rows.Next() {
	// 	var building entity.Building
	// 	err = rows.Scan(&building.ID, &building.Name, &building.AnnualPrice, &building.MonthlyPrice, &building.Owner, &building.IsPublished, &building.City, &building.District)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	buildings = append(buildings, building)
	// }

	// return &buildings, nil
}

func (b *BuildingRepositoryImpl) GetBuildingDetailByID(ctx context.Context, id string, isPublishedOnly bool) (*entity.Building, error) {
	building := &entity.Building{}

	// TODO: Optimize this query (maybe use raw query instead of gorm)
	query := b.db.WithContext(ctx).
		Preload("Pictures", func(db *gorm.DB) *gorm.DB {
			return db.Order("`pictures`.`index` ASC")
		}).
		Preload("Facilities", func(db *gorm.DB) *gorm.DB {
			return db.Joins("Category")
		}).
		Preload("CreatedBy.Detail.Picture").
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

func (b *BuildingRepositoryImpl) GetCities(ctx context.Context) (*entity.Cities, error) {
	cities := new(entity.Cities)
	err := b.db.WithContext(ctx).
		Model(entity.City{}).
		Find(cities).Error
	if err != nil {
		return nil, err
	}

	return cities, nil
}
func (b *BuildingRepositoryImpl) GetDistrictsByCityID(ctx context.Context, cityID int) (*entity.Districts, error) {
	districts := new(entity.Districts)
	err := b.db.WithContext(ctx).
		Model(entity.District{}).
		Where("city_id = ?", cityID).
		Find(districts).Error
	if err != nil {
		return nil, err
	}

	return districts, nil
}

func (b *BuildingRepositoryImpl) GetDistrictByID(ctx context.Context, districtID int) (*entity.District, error) {
	district := new(entity.District)
	err := b.db.WithContext(ctx).
		Model(entity.District{}).
		Where("id = ?", districtID).
		First(district).Error
	if err != nil {
		return nil, err
	}

	return district, nil
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
		switch {
		case strings.Contains(err.Error(), "CONSTRAINT `fk_buildings_facilities` FOREIGN KEY (`building_id`)"):
			return err2.ErrBuildingNotFound
		case strings.Contains(err.Error(), "CONSTRAINT `fk_facilities_category` FOREIGN KEY (`category_id`)"):
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
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		// change building isPublished to false
		return tx.WithContext(ctx).
			Model(&entity.Building{}).
			Where("id = ?", buildingID).
			Updates(entity.Building{IsPublished: custom.Bool(true)}).Error
	} else {
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
