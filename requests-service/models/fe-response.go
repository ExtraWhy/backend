package feresponse

type Fe_resp struct {
	Id    uint64  `json:"id"`
	Won   uint64  `json:"won"`
	Name  string  `json:"name"` //todo to be removed
	Lines []uint8 `json:"lines"`
	Reels []uint8 `json:"reels"`
}

type Fe_resp_cleo struct {
	Pay   float32  `json:"Pay,omitempty"`
	Mult  float32  `json:"Mult,omitempty"`
	Sym   uint32   `json:"Sym,omitempty"`
	Num   []uint32 `json:"Num,omitempty"`
	Line  uint32   `json:"Line,omitempty"`
	Linex []uint32 `json:"Linex,omitempty"`
	Free  uint32   `json:"Free,omitempty"`
	BID   uint32   `json:"BID,omitempty"`
	Bon   any      `json:"Bon,omitempty"`
	JID   uint32   `json:"JID,omitempty"`
	Jack  float32  `json:"Jack,omitempty"`
}

type Fe_resp_slots struct {
	Cleo []Fe_resp_cleo `json:"cleo"`
}
