//go:build !prod || full || igt

package cleopatra

import (
	"context"
	"fmt"
	"io"
	slot "proto/player/server/slots"
)

func CalcStatBon(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	g.FSR = 15 // set free spins mode
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = float64(s.Count(1))
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var q = float64(s.FreeCount()) / reshuf
		var sq = 1 / (1 - q)
		var rtp = sq * rtpsym
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits()))
		fmt.Fprintf(w, "RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%\n", sq, rtpsym, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx, mrtp)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = float64(s.Count(1))
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var q = float64(s.FreeCount()) / reshuf
		var sq = 1 / (1 - q)
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits()))
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}

var Info = slot.GameInfo{
	Aliases: []slot.GameAlias{
		{Prov: "IGT", Name: "Cleopatra"}, // see: https://www.slotsmate.com/software/igt/igt-cleopatra
	},
	GP: slot.GPlpay |
		slot.GPsel |
		slot.GPretrig |
		slot.GPfgmult |
		slot.GPscat |
		slot.GPwild |
		slot.GPwmult,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: slot.MakeRtpList(ReelsMap),
}

func init() {
	Info.SetupFactory(func() any { return NewGame() }, CalcStatReg)
}
