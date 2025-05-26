package slots

type CombinationGenerator struct {
	a, b, c, d, e int
	RX, RY        int
	done          bool
}

func NewCombinationGenerator(rangex, rangey int) *CombinationGenerator {
	return &CombinationGenerator{rangex, rangex, rangex, rangex, 0, rangex, rangey, false}
}

func (g *CombinationGenerator) Next() ([5]int, bool) {
	if g.done {
		return [5]int{}, false
	}

	// Advance e
	g.e++
	if g.e > g.RY {
		g.e = g.RX
		g.d++
		if g.d > g.RY {
			g.d = g.RX
			g.c++
			if g.c > g.RY {
				g.c = g.RX
				g.b++
				if g.b > g.RY {
					g.b = g.RX
					g.a++
					if g.a > g.RY {
						g.done = true
						return [5]int{}, false
					}
				}
			}
		}
	}

	return [5]int{g.a, g.b, g.c, g.d, g.e}, true
}
