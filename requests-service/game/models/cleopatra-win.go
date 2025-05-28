package models

type CleopatraWin struct {
	Pay   float32
	Mult  float32
	Sym   uint32
	Num   uint32
	Line  uint32
	Linex []uint32
	Free  uint32
	BID   uint32
	Bon   interface{}
	JID   uint32
	Jack  float32
}

type CleopatraWins struct {
	Wins []CleopatraWin
	Syms []int32
}
