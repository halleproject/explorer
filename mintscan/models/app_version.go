package models

type (
	// AppVersion defines the structure for App version API
	AppVersion struct {
		ID          int32  `json:"id" sql:",pk"`
		Passwd      string `json:"passwd" sql:",notnull"`
		Version     string `json:"version" sql:",notnull"`
		DownloadURL string `json:"downloadurl" sql:",notnull"`
	}

	ResultVersion struct {
		Version     string `json:"version" sql:",notnull"`
		DownloadURL string `json:"downloadurl" sql:",notnull"`
	}
)
