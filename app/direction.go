package app

import "github.com/go-gl/mathgl/mgl32"

type Direction uint

const (
	north Direction = iota // -z
	south                  // +z
	down                   // -y
	up                     // +y
	west                   // -x
	east                   // +x
	none                   // not calculated
)

var directions []mgl32.Vec3 = []mgl32.Vec3{
	// north
	{0, 0, -1},

	// south
	{0, 0, 1},

	// down
	{0, -1, 0},

	// up
	{0, 1, 0},

	// west
	{-1, 0, 0},

	// east
	{1, 0, 0},
}

// Returns a new boxface with the maching direction.
func newDirection(p mgl32.Vec3) Direction {
	for i, v := range directions {
		if v == p {
			return Direction(i)
		}
	}

	return none
}

func (d Direction) Direction() mgl32.Vec3 {
	new := directions[d]
	return new
}