package app

import (
	"math/rand"
)

type NoiseMapGenerator struct {
	seed int64
	perm []int
}

func newNoiseMapGenerator() *NoiseMapGenerator {
	n := &NoiseMapGenerator{}
	return n
}

func (n *NoiseMapGenerator) Generate3D(width, height, depth int, scale float32) [][][]float32 {
	o := make([][][]float32, 0)
	for i := 0; i < width; i++ {
		o = append(o, make([][]float32, 0))
		for j := 0; j < height; j++ {
			o[i] = append(o[i], make([]float32, 0))
			for k := 0; k < depth; k++ {
				noise := n.altitudePerlinNoise3D(float32(i)*scale, float32(j)*scale, float32(k), height, n.perm)
				noise += 1
				noise /= 2
				o[i][j] = append(o[i][j], noise)
			}
		}
	}
	return o
}

func (n *NoiseMapGenerator) Generate2D(width, depth, height int, scale float32) [][]float32 {
	o := make([][]float32, 0)
	for i := 0; i < depth; i++ {
		o = append(o, make([]float32, 0))
		for j := 0; j < width; j++ {
			noise := n.perlinNoise2D(float32(j)*scale, float32(i)*scale, n.perm)
			noise += 1
			noise /= 2
			noise *= float32(height)
			o[i] = append(o[i], noise)
		}
	}
	return o
}

func (n *NoiseMapGenerator) Seed(seed int64) {
	n.seed = seed
	n.perm = n.generatePermutation(seed)
}

// Fade function for smoothing
func (n *NoiseMapGenerator) fade(t float32) float32 {
	return t * t * t * (t*(t*6-15) + 10)
}

// Linear interpolation
func (n *NoiseMapGenerator) lerp(t, a, b float32) float32 {
	return a + t*(b-a)
}

// Gradient function
func (n *NoiseMapGenerator) grad3D(hash int, x, y, z float32) float32 {
	h := hash & 3
	var u float32
	var v float32
	if h < 8 || h == 12 || h == 13 {
		u = x
	} else {
		u = y
	}
	if h < 4 || h == 12 || h == 13 {
		v = y
	} else {
		v = z
	}
	if h&1 == 0 {
		u = -u
	}
	if h&2 == 0 {
		v = -v
	}
	return u + v
}

func (n *NoiseMapGenerator) grad2D(hash int, x, y float32) float32 {
	h := hash & 3
	u := x
	if h&1 == 0 {
		u = -x
	}
	v := y
	if h&2 == 0 {
		v = -y
	}
	return u + v
}

func (n *NoiseMapGenerator) altitudePerlinNoise3D(x, y, z float32, height int, perm []int) float32 {
	base := n.perlinNoise3D(x, y, z, perm)
	// TODO:
	altFactor := 1.0 - max(0, y-10)/float32(height)
	if altFactor < 0 {
		altFactor = 0
	}
	return base * altFactor
}

func (n *NoiseMapGenerator) perlinNoise3D(x, y, z float32, perm []int) float32 {
	// Find unit grid cell containing point
	x0 := floor(x)
	y0 := floor(y)
	z0 := floor(z)

	// Relative coordinates within grid cell
	relX := x - x0
	relY := y - y0
	relZ := z - z0

	// Wrap the integer grid points (for permutation table lookup)
	x0i := int(x0) & 255
	y0i := int(y0) & 255
	z0i := int(z0) & 255

	// Gradient hashes for the eight cube corners
	h000 := perm[perm[perm[x0i]+y0i]+z0i]
	h001 := perm[perm[perm[x0i]+y0i]+z0i+1]
	h010 := perm[perm[perm[x0i]+y0i+1]+z0i]
	h011 := perm[perm[perm[x0i]+y0i+1]+z0i+1]
	h100 := perm[perm[perm[x0i+1]+y0i]+z0i]
	h101 := perm[perm[perm[x0i+1]+y0i]+z0i+1]
	h110 := perm[perm[perm[x0i+1]+y0i+1]+z0i]
	h111 := perm[perm[perm[x0i+1]+y0i+1]+z0i+1]

	// Gradient contributions
	g000 := n.grad3D(h000, relX, relY, relZ)
	g001 := n.grad3D(h001, relX, relY, relZ-1)
	g010 := n.grad3D(h010, relX, relY-1, relZ)
	g011 := n.grad3D(h011, relX, relY-1, relZ-1)
	g100 := n.grad3D(h100, relX-1, relY, relZ)
	g101 := n.grad3D(h101, relX-1, relY, relZ-1)
	g110 := n.grad3D(h110, relX-1, relY-1, relZ)
	g111 := n.grad3D(h111, relX-1, relY-1, relZ-1)

	// Fade curves for each coordinate
	u := n.fade(relX)
	v := n.fade(relY)
	w := n.fade(relZ)

	// Interpolate along x, then y, then z
	lx00 := n.lerp(u, g000, g100)
	lx01 := n.lerp(u, g001, g101)
	lx10 := n.lerp(u, g010, g110)
	lx11 := n.lerp(u, g011, g111)

	ly0 := n.lerp(v, lx00, lx10)
	ly1 := n.lerp(v, lx01, lx11)

	return n.lerp(w, ly0, ly1)
}

func (n *NoiseMapGenerator) perlinNoise2D(x, y float32, perm []int) float32 {
	x0 := floor(x)
	y0 := floor(y)
	x1 := x0 + 1
	y1 := y0 + 1

	// Compute relative coordinates
	relX := x - x0
	relY := y - y0

	// Wrap coordinates to permutation table
	x0i := int(x0) & 255
	y0i := int(y0) & 255
	x1i := int(x1) & 255
	y1i := int(y1) & 255

	// Calculate hash values
	h00 := perm[perm[x0i]+y0i]
	h10 := perm[perm[x1i]+y0i]
	h01 := perm[perm[x0i]+y1i]
	h11 := perm[perm[x1i]+y1i]

	// Compute gradients
	g00 := n.grad2D(h00, relX, relY)
	g10 := n.grad2D(h10, relX-1, relY)
	g01 := n.grad2D(h01, relX, relY-1)
	g11 := n.grad2D(h11, relX-1, relY-1)

	// Compute fade values
	u := n.fade(relX)
	v := n.fade(relY)

	// Interpolate
	lx0 := n.lerp(u, g00, g10)
	lx1 := n.lerp(u, g01, g11)
	return n.lerp(v, lx0, lx1)
}

// Generate a permutation table
func (n *NoiseMapGenerator) generatePermutation(seed int64) []int {
	perm := make([]int, 256)
	for i := range perm {
		perm[i] = i
	}

	gen := rand.New(rand.NewSource(seed))
	gen.Shuffle(len(perm), func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	// Duplicate for easier wrapping
	return append(perm, perm...)
}