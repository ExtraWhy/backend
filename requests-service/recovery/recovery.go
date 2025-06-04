package recovery

import (
	"github.com/ExtraWhy/internal-libs/models/player"
	"github.com/ExtraWhy/internal-libs/models/user"
)

const (
	CRW_RecordOK = 0
	CRW_Failed   = 1
	CRW_DB_Error = 2
	CRW_Unknown  = 0xff
)

// ensure to specialize types for the record example is to allow player or user or anytihg that is json formatted
type Record interface {
	player.Player | user.User
}

func AddRecord[T Record](rec T) (uint8, error) {
	return CRW_RecordOK, nil
}

// test only for now !!!
