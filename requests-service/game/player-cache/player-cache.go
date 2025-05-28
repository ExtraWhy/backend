package playercache

import (
	"sort"

	"github.com/ExtraWhy/internal-libs/models/player"
)

type cachedPlayer struct {
	Hits uint64
	Pl   player.Player
}
type skv struct {
	k uint64
	v cachedPlayer
}

// bumpy cache to display most recent players that won
var players_cache = make(map[uint64]cachedPlayer)

func DropThem() {
	tb := []skv{}
	if len(players_cache) < 5 {
		return
	}
	for k, v := range players_cache {
		tb = append(tb, skv{k, v})
	}
	sort.Slice(tb, func(i, j int) bool {
		return tb[i].v.Hits > tb[j].v.Hits
	})
	for i := 0; i < len(tb)-5; i++ {
		delete(players_cache, tb[i].k)
	}
}

func CacheSize() int {
	return len(players_cache)
}

func GetThem() []player.Player {
	pl := []player.Player{}
	for _, v := range players_cache {
		pl = append(pl, v.Pl)
	}
	return pl
}

func PutToCache(pl *player.Player) {

	if found, ok := players_cache[pl.Id]; ok {
		players_cache[pl.Id] = cachedPlayer{Hits: found.Hits + 1, Pl: *pl}
	} else {
		players_cache[pl.Id] = cachedPlayer{Hits: 1, Pl: *pl}
	}
}
