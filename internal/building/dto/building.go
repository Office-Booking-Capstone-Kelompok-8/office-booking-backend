package dto

import "office-booking-backend/pkg/entity"

type BriefBuildingResponse struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Pictures string    `json:"pictures"`
	Prices   *Price    `json:"price"`
	Owner    string    `json:"owner"`
	Location *Location `json:"location"`
}

func NewBriefBuildingResponse(building *entity.Building) *BriefBuildingResponse {
	return &BriefBuildingResponse{
		ID:       building.ID,
		Name:     building.Name,
		Pictures: building.Pictures[0].ThumbnailUrl,
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

type BriefBuildingsResponse []BriefBuildingResponse

func NewBriefBuildingsResponse(buildings *entity.Buildings) *BriefBuildingsResponse {
	var briefBuildings BriefBuildingsResponse
	for _, building := range *buildings {
		briefBuildings = append(briefBuildings, *NewBriefBuildingResponse(&building))
	}
	return &briefBuildings
}

type Facility struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	IconName    string `json:"iconName"`
	Description string `json:"description"`
}

func NewFacility(facility *entity.Facility) *Facility {
	return &Facility{
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
	AnnualPrice  int `json:"annual"`
	MonthlyPrice int `json:"monthly"`
}

type Picture struct {
	ID    string `json:"id"`
	Index int    `json:"index"`
	Url   string `json:"url"`
	Alt   string `json:"alt"`
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

type FullBuildingResponse struct {
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
	}
}
