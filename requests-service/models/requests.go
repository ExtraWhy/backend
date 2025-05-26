package models

import (
	pb "github.com/ExtraWhy/internal-libs/proto-models/player"
)

type WinRequest struct {
	PlayerRequest     *pb.PlayerRequest
	PlayerResponse    *pb.PlayerResponse
	CleopatraResponse *pb.CleopatraWins
}
