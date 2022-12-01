package dto

import "office-booking-backend/pkg/entity"

type BriefPublishedBuildingResponse struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Pictures string    `json:"pictures"`
	Prices   *Price    `json:"price"`
	Owner    string    `json:"owner"`
	Location *Location `json:"location"`
}

func NewBriefPublishedBuildingResponse(building *entity.Building) *BriefPublishedBuildingResponse {
	pictureUrl := ""
	if len(building.Pictures) > 0 {
		pictureUrl = building.Pictures[0].ThumbnailUrl
	}

	return &BriefPublishedBuildingResponse{
		ID:       building.ID,
		Name:     building.Name,
		Pictures: pictureUrl,
		Prices: &Price{
			AnnualPrice:  building.AnnualPrice,
			MonthlyPrice: building.MonthlyPrice,
		},
		Owner: building.Owner,
		Location: &Location{
			City:     building.City.Name,
			District: building.District.Name,
		},
	}
}

type BriefPublishedBuildingsResponse []BriefPublishedBuildingResponse

func NewBriefPublishedBuildingsResponse(buildings *entity.Buildings) *BriefPublishedBuildingsResponse {
	var briefBuildings BriefPublishedBuildingsResponse
	for _, building := range *buildings {
		briefBuildings = append(briefBuildings, *NewBriefPublishedBuildingResponse(&building))
	}
	return &briefBuildings
}

type BriefBuildingResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Pictures    string    `json:"pictures"`
	Prices      *Price    `json:"price"`
	Owner       string    `json:"owner"`
	Location    *Location `json:"location"`
	IsPublished bool      `json:"isPublished"`
}

func NewBriefBuildingResponse(building *entity.Building) *BriefBuildingResponse {
	pictureUrl := ""
	if len(building.Pictures) > 0 {
		pictureUrl = building.Pictures[0].ThumbnailUrl
	}

	return &BriefBuildingResponse{
		ID:       building.ID,
		Name:     building.Name,
		Pictures: pictureUrl,
		Prices: &Price{
			AnnualPrice:  building.AnnualPrice,
			MonthlyPrice: building.MonthlyPrice,
		},
		Owner: building.Owner,
		Location: &Location{
			City:     building.City.Name,
			District: building.District.Name,
		},
		IsPublished: building.IsPublished,
	}
}

type BriefBuildingsResponse []BriefBuildingResponse

func NewBriefBuildingsResponse(buildings *entity.Buildings) *BriefBuildingsResponse {
	var briefBuildings BriefBuildingsResponse
	for _, building := range *buildings {
		briefBuildings = append(briefBuildings, *NewBriefBuildingResponse(&building))
	}
	return &briefBuildings
}

type Facility struct {
	ID          int    `json:"id"`
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Icon        string `json:"icon" `
	IconName    string `json:"iconName"`
	Description string `json:"description" validate:"omitempty,min=3,max=100"`
}

func NewFacility(facility *entity.Facility) *Facility {
	return &Facility{
		ID:          facility.ID,
		Name:        facility.Name,
		Icon:        facility.Category.Url,
		IconName:    facility.Category.Name,
		Description: facility.Description,
	}
}

type Facilities []Facility

func NewFacilities(facilities *entity.Facilities) *Facilities {
	var facs Facilities
	for _, facility := range *facilities {
		facs = append(facs, *NewFacility(&facility))
	}
	return &facs
}

type Price struct {
	AnnualPrice  int `json:"annual" validate:"required"`
	MonthlyPrice int `json:"monthly" validate:"required"`
}

type Picture struct {
	ID    string `json:"id" validate:"required"`
	Index int    `json:"index" validate:"required,gte=0,lte=9"`
	Url   string `json:"url" validate:"required,url"`
	Alt   string `json:"alt" validate:"omitempty,min=3,max=100"`
}

func NewPicture(picture *entity.Picture) *Picture {
	return &Picture{
		ID:    picture.ID,
		Index: picture.Index,
		Url:   picture.Url,
		Alt:   picture.Alt,
	}
}

type Pictures []Picture

func NewPictures(pictures *entity.Pictures) *Pictures {
	var pics Pictures
	for _, picture := range *pictures {
		pics = append(pics, *NewPicture(&picture))
	}
	return &pics
}

type Location struct {
	Address  string `json:"address,omitempty"`
	City     string `json:"city"`
	District string `json:"district"`
	Geo      *Geo   `json:"geo,omitempty"`
}

type Geo struct {
	Longitude float64 `json:"long"`
	Latitude  float64 `json:"lat"`
}

type FullPublishedBuildingResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Pictures    *Pictures   `json:"pictures"`
	Description string      `json:"description"`
	Facilities  *Facilities `json:"facilities"`
	Capacity    int         `json:"capacity"`
	Prices      *Price      `json:"price"`
	Owner       string      `json:"owner"`
	Locations   *Location   `json:"location"`
}

func NewFullPublishedBuildingResponse(building *entity.Building) *FullPublishedBuildingResponse {
	return &FullPublishedBuildingResponse{
		ID:          building.ID,
		Name:        building.Name,
		Pictures:    NewPictures(&building.Pictures),
		Description: building.Description,
		Facilities:  NewFacilities(&building.Facilities),
		Capacity:    building.Capacity,
		Prices: &Price{
			AnnualPrice:  building.AnnualPrice,
			MonthlyPrice: building.MonthlyPrice,
		},
		Owner: building.Owner,
		Locations: &Location{
			Address:  building.Address,
			City:     building.City.Name,
			District: building.District.Name,
			Geo: &Geo{
				Longitude: building.Longitude,
				Latitude:  building.Latitude,
			},
		},
	}
}

type FullBuildingResponse struct {
	ID          string      `json:"id" validate:"required,uuid"`
	Name        string      `json:"name" validate:"required,min=3,max=100"`
	Pictures    *Pictures   `json:"pictures" validate:"required,min=1,dive"`
	Description string      `json:"description" validate:"required,min=3,max=10000"`
	Facilities  *Facilities `json:"facilities" validate:"required,min=1,dive"`
	Capacity    int         `json:"capacity" validate:"required,min=1,max=100"`
	Prices      *Price      `json:"price" validate:"required,dive"`
	Owner       string      `json:"owner" validate:"required"`
	Locations   *Location   `json:"location" validate:"required,dive"`
	IsPublished bool        `json:"isPublished" `
}

func NewFullBuildingResponse(building *entity.Building) *FullBuildingResponse {
	return &FullBuildingResponse{
		ID:          building.ID,
		Name:        building.Name,
		Pictures:    NewPictures(&building.Pictures),
		Description: building.Description,
		Facilities:  NewFacilities(&building.Facilities),
		Capacity:    building.Capacity,
		Prices: &Price{
			AnnualPrice:  building.AnnualPrice,
			MonthlyPrice: building.MonthlyPrice,
		},
		Owner: building.Owner,
		Locations: &Location{
			Address:  building.Address,
			City:     building.City.Name,
			District: building.District.Name,
			Geo: &Geo{
				Longitude: building.Longitude,
				Latitude:  building.Latitude,
			},
		},
		IsPublished: building.IsPublished,
	}
}

type FacilityCategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func NewFacilityCategoryResponse(category *entity.Category) *FacilityCategoryResponse {
	return &FacilityCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Url:  category.Url,
	}
}

type FacilityCategoriesResponse []FacilityCategoryResponse

func NewFacilityCategoriesResponse(categories *entity.Categories) *FacilityCategoriesResponse {
	var categoriesResponse FacilityCategoriesResponse
	for _, category := range *categories {
		categoriesResponse = append(categoriesResponse, *NewFacilityCategoryResponse(&category))
	}
	return &categoriesResponse
}

type AddPictureResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
	Alt string `json:"alt"`
}

func NewAddPictureResponse(picture *entity.Picture) *AddPictureResponse {
	return &AddPictureResponse{
		ID:  picture.ID,
		URL: picture.Url,
		Alt: picture.Alt,
	}
}
