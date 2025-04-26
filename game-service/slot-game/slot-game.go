package slot

import (
	"math/rand"
	"proto/player/server/bitvector"
	"proto/player/server/cleopatra"
	"proto/player/server/slots"
	"sync"
	"time"

	"github.com/ExtraWhy/internal-libs/models/games"
)

var gameMode *games.Game
var (
	bvec *bitvector.Bitvector
	once sync.Once
)

func bvecInst() *bitvector.Bitvector {
	once.Do(func() {
		bvec = &bitvector.Bitvector{}
		bvec.NewBitvector(1)
	})
	return bvec
}

// check wilds/scatter/lines
func check_help3(a, b, wildscat uint8) bool {
	return a == b || b == wildscat
}

// check win and position
func check_help2(a [5][3]uint8, g *games.Game) ([]uint8, [5][3]uint8) {
	var scat, i, l = 0, 0, 0
	var lenl = len(g.Lines)
	for l = 0; l < lenl; l++ {
		//		for i = 0; i < 5 && a[i][g.Lines[l][i]-1] == a[0][g.Lines[l][0]-1]; i++ {
		//		}
		for i = 0; i < 5 && check_help3(a[i][g.Lines[l][i]-1], a[0][g.Lines[l][0]-1], 13); i++ {
		}
		for scat = 0; scat < 5 && check_help3(a[scat][g.Lines[l][scat]-1], a[0][g.Lines[l][0]-1], 1); scat++ {
		}
		if i == 5 {
			bvecInst().Add(l)
		}

	}

	return bvecInst().Indices(), a
}

// slide window by 1 back 1 forth [-1 , rand element , +1]
func check_help1(a []uint8, g *games.Game) ([]uint8, [5][3]uint8) {

	var m1 = [5][3]uint8{}
	for i := 0; i < 5; i++ {
		if a[i] >= uint8(len(g.Reels[0]))-1 {
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

func CheckWin(a []uint8, g *games.Game) ([]uint8, [5][3]uint8) {

	return check_help1(a, g)
}

func rand_eng(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func CleopatraSpin(bet uint64) *slots.Wins {
	var wins slots.Wins
	cl := cleopatra.NewGame()
	cl.Spin(99)
	cl.Scanner(&wins)
	return &wins
}

func RollLines() (uint64, []uint8, [5][3]uint8) {
	var data = make([]uint8, 5)

	var multiplyer uint64
	for j := 0; j < 5; j++ {
		ridx := rand_eng(0, len(gameMode.Reels[j])-1)
		data[j] = uint8(ridx)
	}
	res, symb := CheckWin(data, gameMode)
	bvecInst().Reset()
	if len(res) > 0 {
		for i := 0; i < len(res); i++ {
			multiplyer += uint64(gameMode.LinePay[0][3]) // take fixed paytable for now
		}
		return multiplyer, res, symb

	}
	return multiplyer, res, symb

}
