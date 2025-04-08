package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ExtraWhy/internal-libs/models/games"
)

// check win and position
func check_help2(a [5][3]uint16, g *games.Game, w *games.Win) *games.Win {
	var i, l = 0, 0
	for l = 0; l < 5; l++ {
		for i = 0; i < 5 && a[i][g.Lines[l][i]] == a[0][g.Lines[l][0]]; i++ {
		}

		if l == 0 && i == 5 {
			w.Mid = 1
		}
		if l == 1 && i == 5 {
			w.Top = 1
		}
		if l == 2 && i == 5 {
			w.Bottom = 1
		}
		if l == 3 && i == 5 {
			w.DLow = 1
		}
		if l == 4 && i == 5 {
			w.DHigh = 1
		}
	}
	return w
}

// slide window by 1 back 1 forth [-1 , rand element , +1]
func check_help1(a []uint16, g *games.Game, w *games.Win) *games.Win {

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
	return check_help2(m1, g, w)
}

func CheckWin(a []uint16, g *games.Game) (*games.Win, error) {
	w := &games.Win{Top: 0, Bottom: 0, Mid: 0, DLow: 0, DHigh: 0}

	if len(a) != 5 {
		return nil, errors.New("logical error, array must be 5 elements")
	}

	return check_help1(a, g, w), nil
}

func rand_eng(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func rand_eng_dummy(min, max int) int {
	return min + max
}
func main() {

	fmt.Println("test")
	var data = make([]uint16, 5)
	for i := 0; i < 1000000; i++ {
		for j := 0; j < 5; j++ {
			ridx := rand_eng(0, len(games.Game1.Reels[j])-1)
			data[j] = uint16(ridx)
		}
		if res, err := CheckWin(data, &games.Game1); err != nil {
			log.Fatal("Exception in logic", err)
		} else {
			fmt.Println(res)
		}
	}

}
