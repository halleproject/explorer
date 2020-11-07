package schema

// TwoAuth defines the schema for two factor auth information
type TwoAuthForDCI struct {
	//ID      int32  `json:"id" sql:",pk"`
	Key     string `json:"key" sql:",notnull"`
	Address string `json:"address" sql:",notnull, unique,pk"`
	Bind    bool   `json:"bind"`
}

type TwoAuth struct {
	ID  int32  `json:"id" sql:",pk"`
	Key string `json:"key" sql:",notnull"`
}

// type TwoAuth struct {
// 	//ID      int32  `json:"id" sql:",pk"`
// 	Key     string `json:"key" sql:",notnull"`
// 	Address string `json:"address" sql:",notnull, unique,pk"`
// 	Bind    bool   `json:"bind"`
// }
