// Package mapgen provides procedural map generation utilities.
// Implements OpenSimplex 2D noise for terrain variety.
package mapgen

import "math"

// Noise2D is a 2D OpenSimplex-style noise generator.
type Noise2D struct {
	perm [512]int
}

// NewNoise2D creates a seeded 2D noise generator using a simple LCG shuffle.
func NewNoise2D(seed int64) *Noise2D {
	n := &Noise2D{}
	// Initialize identity permutation
	var base [256]int
	for i := range base {
		base[i] = i
	}
	// Fisher-Yates shuffle using a simple LCG seeded RNG
	s := seed
	for i := 255; i > 0; i-- {
		s = (s*6364136223846793005 + 1442695040888963407) // LCG step
		j := int(uint64(s>>16) % uint64(i+1))
		if j < 0 {
			j = -j
		}
		base[i], base[j] = base[j], base[i]
	}
	for i := 0; i < 512; i++ {
		n.perm[i] = base[i&255]
	}
	return n
}

// Eval returns noise at (x, y) in the range approximately [-1, 1].
// Uses a 2D simplex grid with gradient lookup.
func (n *Noise2D) Eval(x, y float64) float64 {
	const (
		F2 = 0.3660254037844386 // (sqrt(3) - 1) / 2
		G2 = 0.2113248654051871 // (3 - sqrt(3)) / 6
	)

	// Skew input space to determine which simplex cell we're in
	s := (x + y) * F2
	i := fastFloor(x + s)
	j := fastFloor(y + s)

	t := float64(i+j) * G2
	X0 := float64(i) - t
	Y0 := float64(j) - t
	x0 := x - X0
	y0 := y - Y0

	// Determine which simplex triangle we're in
	var i1, j1 int
	if x0 > y0 {
		i1, j1 = 1, 0
	} else {
		i1, j1 = 0, 1
	}

	x1 := x0 - float64(i1) + G2
	y1 := y0 - float64(j1) + G2
	x2 := x0 - 1.0 + 2.0*G2
	y2 := y0 - 1.0 + 2.0*G2

	ii := i & 255
	jj := j & 255

	// Calculate contribution from three corners
	var val float64

	t0 := 0.5 - x0*x0 - y0*y0
	if t0 >= 0 {
		t0 *= t0
		gi := n.perm[ii+n.perm[jj]] & 7
		val += t0 * t0 * grad2(gi, x0, y0)
	}

	t1 := 0.5 - x1*x1 - y1*y1
	if t1 >= 0 {
		t1 *= t1
		gi := n.perm[ii+i1+n.perm[jj+j1]] & 7
		val += t1 * t1 * grad2(gi, x1, y1)
	}

	t2 := 0.5 - x2*x2 - y2*y2
	if t2 >= 0 {
		t2 *= t2
		gi := n.perm[ii+1+n.perm[jj+1]] & 7
		val += t2 * t2 * grad2(gi, x2, y2)
	}

	// Scale to [-1, 1]
	return 70.0 * val
}

// FBM returns fractal Brownian motion (layered noise) for richer terrain.
// octaves controls detail layers, lacunarity is frequency multiplier,
// persistence is amplitude decay per octave.
func (n *Noise2D) FBM(x, y float64, octaves int, lacunarity, persistence float64) float64 {
	var total, amplitude, frequency float64
	amplitude = 1.0
	frequency = 1.0
	var maxAmp float64

	for i := 0; i < octaves; i++ {
		total += n.Eval(x*frequency, y*frequency) * amplitude
		maxAmp += amplitude
		amplitude *= persistence
		frequency *= lacunarity
	}

	return total / maxAmp
}

func fastFloor(x float64) int {
	xi := int(x)
	if x < float64(xi) {
		return xi - 1
	}
	return xi
}

// 2D gradient vectors (8 directions)
var grad2Table = [8][2]float64{
	{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	{1, 1}, {-1, 1}, {1, -1}, {-1, -1},
}

func grad2(hash int, x, y float64) float64 {
	g := grad2Table[hash]
	return g[0]*x + g[1]*y
}

// DistFromCenter returns normalized distance from (0,0) for the given map half-size.
// Returns 0.0 at center, 1.0 at the edge.
func DistFromCenter(x, y, halfSize int) float64 {
	dx := float64(x) / float64(halfSize)
	dy := float64(y) / float64(halfSize)
	return math.Sqrt(dx*dx+dy*dy) / math.Sqrt(2)
}
