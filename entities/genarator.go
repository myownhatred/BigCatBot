package entities

type GenModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Signa struct {
	Callsign string     `json:"callsign"`
	Models   []GenModel `json:"models"`
	Armed    bool       `json:"armed"`
}
