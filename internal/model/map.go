package model

import (
	"math/rand"

	"github.com/capucinoxx/forlorn/pkg/codec"
	"github.com/capucinoxx/forlorn/pkg/model"
)

const (
	mapWidth  = 10
	mapHeight = 10
	cellSize  = 10.0
)

// Grid represents a 2D grid structure in a game environment.
type Grid struct {
	height, width int
	cells         map[model.Point]map[model.Point]bool
}

// isInBounds returns true if the model.Point is within the grid's boundaries, otherwise false.
func (g *Grid) isInBounds(pos *model.Point) bool {
	return pos.X >= 0 && pos.X < float32(g.width) && pos.Y >= 0 && pos.Y < float32(g.height)
}

// wallCount returns the number of walls surrounding a given model.Point.
func (g *Grid) wallCount(p model.Point) int {
	if tile, ok := g.cells[p]; ok {
		return len(tile)
	}
	return 0
}

// generateGrid creates and returns a new Grid of the specified dimensions.
func generateGrid(width, height int) *Grid {
	grid := &Grid{
		height: height,
		width:  width,
		cells:  make(map[model.Point]map[model.Point]bool),
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			grid.cells[model.Point{X: float32(x), Y: float32(y)}] = make(map[model.Point]bool)
		}
	}

	visited := make(map[model.Point]bool)

	var dfs func(model.Point)
	dfs = func(pos model.Point) {
		visited[pos] = true

		dirs := make([]model.Point, len(model.Directions))
		copy(dirs, model.Directions)
		rand.Shuffle(len(dirs), func(i, j int) {
			dirs[i], dirs[j] = dirs[j], dirs[i]
		})

		for _, dir := range dirs {
			v := dir.Add(&pos)
			if grid.isInBounds(v) && !visited[*v] {
				dfs(*v)
				grid.cells[pos][dir] = true
				grid.cells[*v][*dir.Reflect(model.NullPoint())] = true
			}
		}
	}

	dfs(model.Point{X: 0, Y: 0})
	randomPos := model.Point{X: float32(rand.Intn(width)), Y: float32(rand.Intn(height))}
	dfs(randomPos)

	return grid
}

// Map represents the game map containing colliders and spawn model.Points.
type Map struct {
	colliders  []*model.Collider
	spawns     []*model.Point
	cellSize   float32
	tilesWalls [][]uint8

	r *rand.Rand
}

// Colliders returns the list of colliders in the map.
func (m *Map) Colliders() []*model.Collider {
	return m.colliders
}

// Spawns returns the list of spawn model.Points in the map.
func (m *Map) Spawns() []*model.Point {
	return m.spawns
}

func (m *Map) Size() int {
	return mapWidth
}

func (m *Map) DiscreteMap() [][]uint8 {
	return m.tilesWalls
}

// Setup initializes the map by setting up the grid and generating spawns.
func (m *Map) Setup() {
	m.cellSize = cellSize
	m.r = rand.New(rand.NewSource(int64(rand.Uint64())))

	grid := generateGrid(mapWidth, mapHeight)
	m.populate(grid)
	m.generateSpawns(grid)
}

// populate fills the map with colliders using the specified grid.
func (m *Map) populate(grid *Grid) {
	m.colliders = []*model.Collider{
		{Points: []*model.Point{{X: 0, Y: 0}, {X: 0, Y: float32(grid.height) * m.cellSize}}},
		{Points: []*model.Point{{X: 0, Y: 0}, {X: float32(grid.width) * m.cellSize, Y: 0}}},
		{Points: []*model.Point{{X: float32(grid.width) * m.cellSize, Y: 0}, {X: float32(grid.width) * m.cellSize, Y: float32(grid.height) * m.cellSize}}},
		{Points: []*model.Point{{X: 0, Y: float32(grid.height) * m.cellSize}, {X: float32(grid.width) * m.cellSize, Y: float32(grid.height) * m.cellSize}}},
	}

	for y := 0; y < grid.height; y++ {
		for x := 0; x < grid.width; x++ {
			cell := grid.cells[model.Point{X: float32(x), Y: float32(y)}]
			x1, y1 := float32(x)*m.cellSize, float32(y)*m.cellSize
			x2, y2 := x1+m.cellSize, y1+m.cellSize

			if _, exist := cell[model.RIGHT]; exist {
				m.colliders = append(m.colliders, &model.Collider{
					Points: []*model.Point{{X: x2, Y: y1}, {X: x2, Y: y2}},
				})
			}

			if _, exist := cell[model.DOWN]; exist {
				m.colliders = append(m.colliders, &model.Collider{
					Points: []*model.Point{{X: x1, Y: y2}, {X: x2, Y: y2}},
				})
			}
		}
	}
}

// generateSpawns generates spawn model.Points for the map based on the grid.
func (m *Map) generateSpawns(grid *Grid) {
	m.tilesWalls = make([][]uint8, mapHeight)

	const (
		minValue  = 0.1
		maxValue  = 0.9
		diffValue = maxValue - minValue
	)

	for i := 0; i < mapHeight; i++ {
		m.tilesWalls[i] = make([]uint8, mapWidth)
		for j := 0; j < mapWidth; j++ {
			wallCount := grid.wallCount(model.Point{X: float32(i), Y: float32(j)})
			m.tilesWalls[i][j] = uint8(wallCount)

			if wallCount != 4 {
				for k := 0; k < 3; k++ {
					m.spawns = append(m.spawns, &model.Point{
						X: float32(j) + (minValue+m.r.Float32()*diffValue)*m.cellSize,
						Y: float32(i) + (minValue+m.r.Float32()*diffValue)*m.cellSize,
					})
				}
			}
		}
	}
}

func (m *Map) Encode(w *codec.ByteWriter) error {
	w.WriteInt8(mapWidth)

	for i := 0; i < mapWidth; i++ {
		for j := 0; j < mapWidth; j++ {
			w.WriteInt8(int8(m.tilesWalls[i][j]))
		}
	}

	w.WriteInt32(int32(len(m.colliders)))
	for _, c := range m.colliders {
		c.Encode(w)
	}

	return nil
}

func (m *Map) Decode(r *codec.ByteReader) error {
	width, err := r.ReadInt8()
	if err != nil {
		return err
	}

	m.tilesWalls = make([][]uint8, width)

	for i := 0; i < int(width); i++ {
		m.tilesWalls[i] = make([]uint8, width)
		for j := 0; j < int(width); j++ {
			m.tilesWalls[i][j], err = r.ReadUint8()
			if err != nil {
				return err
			}
		}
	}

	collidersLen, err := r.ReadInt32()
	if err != nil {
		return err
	}

	m.colliders = make([]*model.Collider, collidersLen)

	for i := 0; i < int(collidersLen); i++ {
		m.colliders[i] = &model.Collider{}
		if err := m.colliders[i].Decode(r); err != nil {
			return err
		}
	}

	return nil
}
