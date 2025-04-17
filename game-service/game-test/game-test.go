package gametest

import (
	"math/rand"
	"proto/player/server/bitvector"
	"sync"
	"time"

	"github.com/ExtraWhy/internal-libs/models/games"
)

var gameMode *games.Game
var (
	bvec     *bitvector.Bitvector
	once     sync.Once
	paytable = []uint64{50, 100, 50,
		25, 25, 25,
		20, 20, 20,
		10, 10, 10, 5}
)

func bvecInst() *bitvector.Bitvector {
	once.Do(func() {
		bvec = &bitvector.Bitvector{}
		bvec.NewBitvector(1)
	})
	return bvec
}

// check win and position
func check_help2(a [5][3]uint16, g *games.Game) ([]uint8, [5][3]uint16) {
	var i, l = 0, 0
	for l = 0; l < 10; l++ {
		for i = 0; i < 5 && a[i][g.Lines[l][i]] == a[0][g.Lines[l][0]]; i++ {
		}
		if i == 5 {
			bvecInst().Add(l)
		}
	}

	return bvecInst().Indices(), a
}

// slide window by 1 back 1 forth [-1 , rand element , +1]
func check_help1(a []uint16, g *games.Game) ([]uint8, [5][3]uint16) {

	var m1 = [5][3]uint16{}
	for i := 0; i < 5; i++ {
		if a[i] >= uint16(len(g.Reels[0]))-1 {
			m1[i][0] = g.Reels[i][a[i]-1]
			m1[i][1] = g.Reels[i][a[i]]
			m1[i][2] = g.Reels[i][0]
		} else if a[i] == 0 {
			m1[i][0] = g.Reels[i][len(g.Reels[0])-1]
			m1[i][1] = g.Reels[i][a[i]]
			m1[i][2] = g.Reels[i][a[i]+1]
		} else {
			m1[i][0] = g.Reels[i][a[i]-1]
			m1[i][1] = g.Reels[i][a[i]]
			m1[i][2] = g.Reels[i][a[i]+1]

		}
	}
	return check_help2(m1, g)
}

func SetupGame(b bool) {
	if b {
		gameMode = &games.GameTest
	} else {
		gameMode = &games.Game1
	}
}

func CheckWin(a []uint16, g *games.Game) ([]uint8, [5][3]uint16) {

	return check_help1(a, g)
}

func rand_eng(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func RollLines() (uint64, []uint8, [5][3]uint16) {
	var data = make([]uint16, 5)
	var multiplyer uint64
	for j := 0; j < 5; j++ {
		ridx := rand_eng(0, len(gameMode.Reels[j])-1)
		data[j] = uint16(ridx)
	}
	res, symb := CheckWin(data, gameMode)
	bvecInst().Reset()
	if len(res) > 0 {
		for i := 0; i < len(res); i++ {
			multiplyer += paytable[res[i]]
		}
		return multiplyer, res, symb

	}
	return multiplyer, res, symb

}
