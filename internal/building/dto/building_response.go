package dto

import (
	"office-booking-backend/pkg/constant"
	"office-booking-backend/pkg/entity"
)

type BriefPublishedBuildingResponse struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Pictures string    `json:"pictures"`
	Review   *Review   `json:"review"`
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
		Review: &Review{
			Rating: building.Rating,
			Count:  building.ReviewCount,
		},
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
	Review      *Review   `json:"review"`
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
		Review: &Review{
			Rating: building.Rating,
			Count:  building.ReviewCount,
		},
		Prices: &Price{
			AnnualPrice:  building.AnnualPrice,
			MonthlyPrice: building.MonthlyPrice,
		},
		Owner: building.Owner,
		Location: &Location{
			City:     building.City.Name,
			District: building.District.Name,
		},
		IsPublished: *building.IsPublished,
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
	Icon        string `json:"icon"`
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
	Index int    `json:"index" validate:"gte=0,lte=9"`
	Url   string `json:"url" validate:"required,url"`
	Alt   string `json:"alt" validate:"omitempty,min=3,max=100"`
}

func NewPicture(picture *entity.Picture) *Picture {
	return &Picture{
		ID:    picture.ID,
		Index: *picture.Index,
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

type Review struct {
	Rating float64 `json:"rating"`
	Count  int     `json:"count"`
}

type Location struct {
	Address  string `json:"address,omitempty"`
	City     string `json:"city"`
	District string `json:"district"`
	Geo      *Geo   `json:"geo,omitempty"`
}

type FullLocation struct {
	Address  string            `json:"address,omitempty"`
	City     *CityResponse     `json:"city"`
	District *DistrictResponse `json:"district"`
	Geo      *Geo              `json:"geo,omitempty"`
}

type Geo struct {
	Longitude float64 `json:"long"`
	Latitude  float64 `json:"lat"`
}

type Status struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

type BriefReservationResponse struct {
	Status    *Status `json:"status"`
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
}

func NewBriefReservationResponse(reservation *entity.Reservation) *BriefReservationResponse {
	return &BriefReservationResponse{
		Status: &Status{
			ID:      reservation.Status.ID,
			Message: reservation.Status.Message,
		},
		StartDate: reservation.StartDate.Format("2006-01-02"),
		EndDate:   reservation.EndDate.Format("2006-01-02"),
	}
}

type BriefReservationsResponse []BriefReservationResponse

func NewBriefReservationsResponse(reservations *entity.Reservations) *BriefReservationsResponse {
	var briefReservations BriefReservationsResponse
	for _, reservation := range *reservations {
		briefReservations = append(briefReservations, *NewBriefReservationResponse(&reservation))
	}
	return &briefReservations
}

type FullPublishedBuildingResponse struct {
	ID           string                     `json:"id"`
	Name         string                     `json:"name"`
	Pictures     *Pictures                  `json:"pictures"`
	Description  string                     `json:"description"`
	Facilities   *Facilities                `json:"facilities"`
	Reservations *BriefReservationsResponse `json:"reservations"`
	Capacity     int                        `json:"capacity"`
	Size         int                        `json:"size"`
	Review       *Review                    `json:"review"`
	Prices       *Price                     `json:"price"`
	Owner        string                     `json:"owner"`
	Locations    *FullLocation              `json:"location"`
	Agent        *Agent                     `json:"agent"`
}

func NewFullPublishedBuildingResponse(building *entity.Building) *FullPublishedBuildingResponse {
	return &FullPublishedBuildingResponse{
		ID:           building.ID,
		Name:         building.Name,
		Pictures:     NewPictures(&building.Pictures),
		Description:  building.Description,
		Facilities:   NewFacilities(&building.Facilities),
		Capacity:     building.Capacity,
		Size:         building.Size,
		Reservations: NewBriefReservationsResponse(&building.Reservations),
		Review: &Review{
			Rating: building.Rating,
			Count:  building.ReviewCount,
		},
		Prices: &Price{
			AnnualPrice:  building.AnnualPrice,
			MonthlyPrice: building.MonthlyPrice,
		},
		Owner: building.Owner,
		Locations: &FullLocation{
			Address:  building.Address,
			City:     NewCityResponse(&building.City),
			District: NewDistrictResponse(&building.District),
			Geo: &Geo{
				Longitude: building.Longitude,
				Latitude:  building.Latitude,
			},
		},
		Agent: NewAgent(&building.CreatedBy),
	}
}

type FullBuildingResponse struct {
	ID           string                     `json:"id" validate:"required,uuid"`
	Name         string                     `json:"name" validate:"required,min=3,max=100"`
	Pictures     *Pictures                  `json:"pictures" validate:"required,min=1,dive"`
	Description  string                     `json:"description" validate:"required,min=3,max=10000"`
	Facilities   *Facilities                `json:"facilities" validate:"required,min=1,dive"`
	Capacity     int                        `json:"capacity" validate:"required,min=1"`
	Size         int                        `json:"size" validate:"required,gte=1"`
	Reservations *BriefReservationsResponse `json:"reservations,omitempty"`
	Review       *Review                    `json:"review"`
	Prices       *Price                     `json:"price" validate:"required,dive"`
	Owner        string                     `json:"owner" validate:"required"`
	Locations    *FullLocation              `json:"location" validate:"required,dive"`
	Agent        *Agent                     `json:"agent,omitempty"`
	IsPublished  bool                       `json:"isPublished" `
}

func NewFullBuildingResponse(building *entity.Building) *FullBuildingResponse {
	return &FullBuildingResponse{
		ID:           building.ID,
		Name:         building.Name,
		Pictures:     NewPictures(&building.Pictures),
		Description:  building.Description,
		Facilities:   NewFacilities(&building.Facilities),
		Capacity:     building.Capacity,
		Size:         building.Size,
		Reservations: NewBriefReservationsResponse(&building.Reservations),
		Review: &Review{
			Rating: building.Rating,
			Count:  building.ReviewCount,
		},
		Prices: &Price{
			AnnualPrice:  building.AnnualPrice,
			MonthlyPrice: building.MonthlyPrice,
		},
		Owner: building.Owner,
		Locations: &FullLocation{
			Address:  building.Address,
			City:     NewCityResponse(&building.City),
			District: NewDistrictResponse(&building.District),
			Geo: &Geo{
				Longitude: building.Longitude,
				Latitude:  building.Latitude,
			},
		},
		Agent:       NewAgent(&building.CreatedBy),
		IsPublished: *building.IsPublished,
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

type CityResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewCityResponse(cities *entity.City) *CityResponse {
	return &CityResponse{
		ID:   cities.ID,
		Name: cities.Name,
	}
}

type CitiesResponse []CityResponse

func NewCitiesResponse(cities *entity.Cities) *CitiesResponse {
	var citiesResponse CitiesResponse
	for _, city := range *cities {
		citiesResponse = append(citiesResponse, *NewCityResponse(&city))
	}
	return &citiesResponse
}

type DistrictResponse struct {
	ID     int    `json:"id"`
	CityID int    `json:"cityId"`
	Name   string `json:"name"`
}

func NewDistrictResponse(district *entity.District) *DistrictResponse {
	return &DistrictResponse{
		ID:     district.ID,
		CityID: district.CityID,
		Name:   district.Name,
	}
}

type DistrictsResponse []DistrictResponse

func NewDistrictsResponse(districts *entity.Districts) *DistrictsResponse {
	var districtsResponse DistrictsResponse
	for _, district := range *districts {
		districtsResponse = append(districtsResponse, *NewDistrictResponse(&district))
	}
	return &districtsResponse
}

type Agent struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Picture string `json:"picture"`
}

func NewAgent(agent *entity.User) *Agent {
	return &Agent{
		ID:      agent.ID,
		Name:    agent.Detail.Name,
		Email:   agent.Email,
		Phone:   agent.Detail.Phone,
		Picture: agent.Detail.Picture.Url,
	}
}

type BuildingStatResponse struct {
	ByCity      *TotalByCities    `json:"byCity"`
	ByTimeFrame *TotalByTimeFrame `json:"byTimeFrame"`
}

type TotalByCity struct {
	CityID   int64  `json:"cityId"`
	CityName string `json:"cityName"`
	Total    int64  `json:"total"`
}

func NewTotalByCity(stat *entity.CityStat) *TotalByCity {
	return &TotalByCity{
		CityID:   stat.CityID,
		CityName: stat.CityName,
		Total:    stat.Total,
	}
}

type TotalByCities []TotalByCity

func NewTotalByCities(stats *entity.CitiesStat) *TotalByCities {
	var totalByCities TotalByCities
	for _, stat := range *stats {
		totalByCities = append(totalByCities, *NewTotalByCity(&stat))
	}
	return &totalByCities
}

type TotalByTimeFrame struct {
	Today     int64 `json:"today"`
	ThisWeek  int64 `json:"thisWeek"`
	ThisMonth int64 `json:"thisMonth"`
	ThisYear  int64 `json:"thisYear"`
	AllTime   int64 `json:"allTime"`
}

func NewTimeframeStat(stats *entity.TimeframeStat) *TotalByTimeFrame {
	return &TotalByTimeFrame{
		Today:     stats.Day.Int64,
		ThisWeek:  stats.Week.Int64,
		ThisMonth: stats.Month.Int64,
		ThisYear:  stats.Year.Int64,
		AllTime:   stats.All.Int64,
	}
}

type BriefBuildingReviewResponse struct {
	ID        string            `json:"id"`
	User      BriefUserResponse `json:"user"`
	Rating    int               `json:"rating"`
	Message   string            `json:"message"`
	CreatedAt string            `json:"createdAt"`
}

func NewBriefBuildingReviewResponse(review *entity.Review) *BriefBuildingReviewResponse {
	return &BriefBuildingReviewResponse{
		ID:        review.ID,
		Rating:    review.Rating,
		Message:   review.Message,
		User:      *NewBriefUserResponse(&review.User),
		CreatedAt: review.CreatedAt.Format(constant.DATE_RESPONSE_FORMAT),
	}
}

type BriefBuildingReviewsResponse []BriefBuildingReviewResponse

func NewBriefBuildingReviewsResponse(reviews *entity.Reviews) *BriefBuildingReviewsResponse {
	response := new(BriefBuildingReviewsResponse)
	for _, review := range *reviews {
		*response = append(*response, *NewBriefBuildingReviewResponse(&review))
	}
	return response
}

type BriefUserResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewBriefUserResponse(user *entity.User) *BriefUserResponse {
	url := user.Detail.Picture.Url
	if url == "" {
		url = constant.DEFAULT_USER_AVATAR
	}
	return &BriefUserResponse{
		ID:      user.ID,
		Name:    user.Detail.Name,
		Picture: url,
	}
}
