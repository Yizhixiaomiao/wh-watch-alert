package types

type RequestTenantSearch struct {
	Keyword string `json:"keyword" form:"keyword"`
	Page    int    `json:"page" form:"page" default:"1"`
	Size    int    `json:"size" form:"size" default:"20"`
}
