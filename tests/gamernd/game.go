package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Game struct {
	Reels [][]uint16
	Lines [][]uint16
}

type WinPos struct {
	Top     bool
	Mid     bool
	Bottom  bool
	UpDiag  bool
	LowDiag bool
}

func (wp *WinPos) reset() {
	wp.Bottom = false
	wp.LowDiag = false
	wp.UpDiag = false
	wp.Top = false
	wp.Mid = false
}

var (
	//99% rtp cleopatra test
	g1 = Game{Reels: [][]uint16{
		{11, 5, 10, 2, 7, 5, 9, 1, 7, 3, 8, 4, 11, 12, 6, 10, 3, 8, 2, 11, 3, 8, 5, 12, 6, 9, 12, 4, 10, 13, 2, 7, 6, 9},
		{11, 12, 3, 9, 8, 4, 9, 5, 7, 2, 9, 6, 7, 2, 8, 5, 11, 6, 8, 10, 3, 12, 2, 10, 11, 1, 10, 4, 13, 6, 7, 5, 12, 3},
		{8, 5, 11, 3, 12, 6, 9, 7, 3, 10, 4, 8, 13, 2, 7, 4, 11, 2, 9, 5, 12, 3, 10, 9, 2, 12, 1, 7, 5, 8, 6, 11, 10, 6},
		{3, 10, 5, 11, 6, 7, 11, 6, 7, 2, 13, 3, 9, 2, 10, 1, 12, 11, 3, 12, 4, 9, 8, 5, 10, 2, 8, 5, 7, 12, 4, 8, 6, 9},
		{6, 7, 10, 9, 1, 12, 6, 11, 2, 9, 3, 8, 2, 7, 5, 10, 7, 4, 8, 3, 12, 13, 6, 11, 5, 8, 2, 10, 5, 12, 4, 9, 3, 11},
	},
		Lines: [][]uint16{
			{1, 1, 1, 1, 1}, // 1
			{0, 0, 0, 0, 0}, // 2
			{2, 2, 2, 2, 2}, // 3
			{0, 1, 2, 1, 0}, // 4
			{2, 1, 0, 1, 2}, // 5
		},
	}
)

var gWin = WinPos{Top: false, Mid: false, Bottom: false, UpDiag: false, LowDiag: false}

func check_help2(a [5][3]uint16) *WinPos {

	for l := 0; l < 5; l++ {
		wnr := a[l][0]
		for i := 0; i < 5; i++ {
			if a[i][g1.Lines[l][i]] == wnr && l == 0 {
				gWin.Top = true
			}
			if a[i][g1.Lines[l][i]] == wnr && l == 1 {
				gWin.Mid = true
			}
			if a[i][g1.Lines[l][i]] == wnr && l == 2 {
				gWin.Bottom = true
			}
			if a[i][g1.Lines[l][i]] == wnr && l == 3 {
				gWin.LowDiag = true
			}
			if a[i][g1.Lines[l][i]] == wnr && l == 4 {
				gWin.UpDiag = true
			}
		}
	}
	return &gWin
}

func check_help1(a []uint16) *WinPos {

	var m1 = [5][3]uint16{}
	for i := 0; i < 5; i++ {
		if a[i] >= uint16(len(g1.Reels[0]))-1 {
			m1[i][0] = g1.Reels[i][a[i]-1]
			m1[i][1] = g1.Reels[i][a[i]]
			m1[i][2] = g1.Reels[i][0]
		} else if a[i] == 0 {
			m1[i][0] = g1.Reels[i][len(g1.Reels[0])-1]
			m1[i][1] = g1.Reels[i][a[i]]
			m1[i][2] = g1.Reels[i][a[i]+1]
		} else {
			m1[i][0] = g1.Reels[i][a[i]-1]
			m1[i][1] = g1.Reels[i][a[i]]
			m1[i][2] = g1.Reels[i][a[i]+1]

		}
	}
	return check_help2(m1)
}

func CheckWin(a []uint16) (*WinPos, error) {
	if len(a) != 5 {
		return nil, errors.New("logical error, array must be 5 elements")
	}

	return check_help1(a), nil
}

func rand_eng(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func main() {

	fmt.Println("test")
	var data = make([]uint16, 5)
	for i := 0; i < 2; i++ {
		for j := 0; j < 5; j++ {
			ridx := rand_eng(0, len(g1.Reels[j])-1)
			data[j] = uint16(ridx)
		}
		if res, err := CheckWin(data); err != nil {
			log.Fatal("Exception in logic", err)
		} else {
			fmt.Println(res)
			gWin.reset()
		}
	}

}
