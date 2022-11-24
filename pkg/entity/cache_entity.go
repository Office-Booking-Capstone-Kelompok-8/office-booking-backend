package entity

type CachedToken struct {
	AccessID  string `json:"accessID"`
	RefreshID string `json:"refreshID"`
}
