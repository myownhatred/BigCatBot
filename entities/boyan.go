package entities

type Boyan struct {
	ID   int      `json:"id"`
	Link string   `json:"link"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type TagCount struct {
	TagID int    `json:"tag_id"`
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}
