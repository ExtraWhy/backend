package slots

type CombinationGenerator struct {
	a, b, c, d, e int
	done          bool
}

func NewCombinationGenerator() *CombinationGenerator {
	return &CombinationGenerator{1, 1, 1, 1, 0, false}
}

func (g *CombinationGenerator) Next() ([5]int, bool) {
	if g.done {
		return [5]int{}, false
	}

	// Advance e
	g.e++
	if g.e > 34 {
		g.e = 1
		g.d++
		if g.d > 34 {
			g.d = 1
			g.c++
			if g.c > 34 {
				g.c = 1
				g.b++
				if g.b > 34 {
					g.b = 1
					g.a++
					if g.a > 34 {
						g.done = true
						return [5]int{}, false
					}
				}
			}
		}
	}

	return [5]int{g.a, g.b, g.c, g.d, g.e}, true
}
