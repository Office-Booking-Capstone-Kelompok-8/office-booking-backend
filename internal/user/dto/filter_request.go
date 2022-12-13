package dto

type UserFilterRequest struct {
	Query  string `query:"q" validate:"omitempty,min=3,max=50"`
	Role   int    `query:"role" validate:"omitempty,oneof=1 2"`
	Page   int    `query:"page" validate:"omitempty,gte=1"`
	Limit  int    `query:"limit" validate:"omitempty,gte=1"`
	Offset int    `query:"-"`
}
