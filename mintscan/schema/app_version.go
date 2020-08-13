package schema

// AppVersion defines the schema for App version information
type AppVersion struct {
	ID          int32  `json:"id" sql:",pk"`
	Passwd      string `json:"passwd" sql:",notnull"`
	Version     string `json:"version" sql:",notnull"`
	DownloadURL string `json:"downloadurl" sql:",notnull"`
}
