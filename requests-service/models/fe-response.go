package feresponse

type Fe_resp struct {
	Id    uint64  `json:"id"`
	Won   uint64  `json:"won"`
	Name  string  `json:"name"` //todo to be removed
	Lines []uint8 `json:"lines"`
	Reels []uint8 `json:"reels"`
}
