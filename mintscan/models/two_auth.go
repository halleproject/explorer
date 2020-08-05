package models

type (
	// TwoAuth defines the structure for TwoAuth API
	TwoAuth struct {
		ID  int32  `json:"id" sql:",pk"`
		Key string `json:"key" sql:",notnull"`
	}
)
