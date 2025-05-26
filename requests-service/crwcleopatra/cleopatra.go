package crwcleopatra

import (
	"casino/bitvector"
	models "casino/models"
	"casino/slots"
	"errors"
	"fmt"
	"sync"

	"github.com/ExtraWhy/internal-libs/models/games"
)

var gameMode *games.Game
var (
	bvec  *bitvector.Bitvector
	once  sync.Once
	once1 sync.Once
	cleo  *Game = NewGame()
)

func set20Lines() {
	once1.Do(func() {
		cleo.SetSel(20)
	})
}

func bvecInst() *bitvector.Bitvector {
	once.Do(func() {
		bvec = &bitvector.Bitvector{}
		bvec.NewBitvector(1)
	})
	return bvec
}

func SetupGame(b bool) {
	if b {
		gameMode = &games.GameTest
	} else {
		gameMode = &games.Game1
	}
}

func CleopatraSpinV2(bet uint64) (*slots.Wins, *Game) {
	set20Lines()
	var wins slots.Wins
	var n = 0

	cleo.Prepare()
	cleo.SetBet(float64(bet))
	for { // repeat until get valid screen
		cleo.Spin(99.517383)
		if cleo.Scanner(&wins) == nil {
			break
		}
		n++
		if n > 300 {
			break
		}
		wins.Reset()
	}
	cleo.Spawn(wins, float64(bet), 99.517383)
	//debit = cost*(1-jprate/100) - wins.Gain()
	//jack = wins.Jackpot()
	//wins.Reset()
	return &wins, cleo
}

func GetWinForCleopatra(msg *models.MessageBet) (*models.CleopatraWins, error) {

	retwins := models.CleopatraWins{}
	retwins.Wins = make([]models.CleopatraWin, 1)
	wins, cl := CleopatraSpinV2(msg.Money)

	for j := 0; j < 5; j++ {
		for i := 0; i < 3; i++ {
			retwins.Syms = append(retwins.Syms, int32(cl.Scr[j][i]))
		}
	}

	for _, j := range *wins {
		bid := uint32(j.BID)
		free := uint32(j.Free)
		jid := uint32(j.JID)
		jack := float32(j.Jack)
		line := uint32(j.Line)
		mult := float32(j.Mult)
		pay := float32(j.Pay)
		num := uint32(j.Num)
		sym := uint32(j.Sym)
		retwins.Wins = append(retwins.Wins, models.CleopatraWin{})

		retwins.Wins[len(retwins.Wins)-1].BID = bid
		retwins.Wins[len(retwins.Wins)-1].Free = free
		retwins.Wins[len(retwins.Wins)-1].JID = jid
		retwins.Wins[len(retwins.Wins)-1].Jack = jack
		retwins.Wins[len(retwins.Wins)-1].Line = line
		retwins.Wins[len(retwins.Wins)-1].Mult = mult
		retwins.Wins[len(retwins.Wins)-1].Pay = pay
		retwins.Wins[len(retwins.Wins)-1].Num = num
		retwins.Wins[len(retwins.Wins)-1].Sym = sym

		for h := 0; h < len(j.XY); h++ {
			retwins.Wins[len(retwins.Wins)-1].Linex = append(retwins.Wins[len(retwins.Wins)-1].Linex, uint32(j.XY[h]))
		}

	}
	return &retwins, nil
}

func SendWin4Cleo(msg *models.MessageBet) (*models.CleopatraWins, error) {
	cleoresp, err := GetWinForCleopatra(msg) //IVZ
	if err != nil {
		return nil, errors.New(fmt.Sprint("could not greet: %v", err))
	}

	return cleoresp, nil
}
