package teststats

import (
	"math/rand"
	tests "tests/game-data"
	"time"
)

var outcomes []int = make([]int, 100)

func rand_eng(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

type rldata struct {
	reels []int
}





func RunReelSlotMachine(run int, pd *tests.GameData) {

	for i := 0; i < run; i++ {
		r1 := rand_eng(0, len(pd.Data[0].Data)-1)
		r2 := rand_eng(0, len(pd.Data[1].Data)-1)
		r3 := rand_eng(0, len(pd.Data[2].Data)-1)
		r4 := rand_eng(0, len(pd.Data[3].Data)-1)
		r5 := rand_eng(0, len(pd.Data[4].Data)-1)
		//		fmt.Printf("[loop %d] r1 %d r2 %d r3 %d r4 %d r5 %d \r\n", i,
		//			pd.Data[0].Data[r1],
		//			pd.Data[1].Data[r2],
		//			pd.Data[2].Data[r3],
		//			pd.Data[3].Data[r4],
		//			pd.Data[4].Data[r5])
		//		for j := 0; j < 5; j++ {
		//		}

	}

}
