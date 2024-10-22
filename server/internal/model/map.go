package model

import (
	"container/heap"
	"math"
	"math/rand"
	"time"

	"github.com/capucinoxx/jdis-games-2024/consts"
	"github.com/capucinoxx/jdis-games-2024/pkg/codec"
	"github.com/capucinoxx/jdis-games-2024/pkg/model"
	"github.com/capucinoxx/jdis-games-2024/pkg/utils"
)

type point struct {
	x, y int
}

var directions = []point{
	{-1, 0},
	{1, 0},
	{0, 1},
	{0, -1},
}

type cell struct {
	n, s, e, w bool
}

func (c cell) isWall(pos int) bool {
	switch pos {
	case 0:
		return c.n
	case 1:
		return c.s
	case 2:
		return c.e
	case 3:
		return c.w
	}

	return false
}

type Map struct {
	size         int
	grid         [][]cell
	discreteGrid [][]uint8
	start        model.Point

	spawns [2][]*model.Point
	walls  []*model.Collider
}

func (m *Map) Centroid() model.Point {
	return m.start
}

func (m *Map) primGenerateMaze(start point) {
	walls := [][5]int{}
	visited := make(map[point]bool)
	visited[start] = true

	for i, dir := range directions {
		dx, dy := dir.x, dir.y
		walls = append(walls, [5]int{start.x + dx, start.y + dy, i, start.x, start.y})
	}

	for len(walls) > 0 {
		idx := rand.Intn(len(walls))
		wall := walls[idx]
		walls = append(walls[:idx], walls[idx+1:]...)
		nx, ny, direction, px, py := wall[0], wall[1], wall[2], wall[3], wall[4]

		if nx >= 0 && nx < consts.MapWidth && ny >= 0 && ny < consts.MapWidth && !visited[point{nx, ny}] {
			m.removeWall(point{px, py}, point{nx, ny}, direction)
			visited[point{nx, ny}] = true

			for i, dir := range directions {
				dx, dy := dir.x, dir.y
				walls = append(walls, [5]int{nx + dx, ny + dy, i, nx, ny})
			}
		}
	}
}

func (m *Map) generateColliders() {
	m.walls = []*model.Collider{}
	for i, row := range m.grid {
		for j, c := range row {
			if c.n {
				m.walls = append(m.walls, &model.Collider{
					Points: []*model.Point{
						{X: float64(j * consts.CellWidth), Y: float64(i * consts.CellWidth)},
						{X: float64((j + 1) * consts.CellWidth), Y: float64(i * consts.CellWidth)},
					},
				})
			}

			if c.s {
				m.walls = append(m.walls, &model.Collider{
					Points: []*model.Point{
						{X: float64(j * consts.CellWidth), Y: float64((i + 1) * consts.CellWidth)},
						{X: float64((j + 1) * consts.CellWidth), Y: float64((i + 1) * consts.CellWidth)},
					},
				})
			}

			if c.e {
				m.walls = append(m.walls, &model.Collider{
					Points: []*model.Point{
						{X: float64((j + 1) * consts.CellWidth), Y: float64(i * consts.CellWidth)},
						{X: float64((j + 1) * consts.CellWidth), Y: float64((i + 1) * consts.CellWidth)},
					},
				})
			}

			if c.w {
				m.walls = append(m.walls, &model.Collider{
					Points: []*model.Point{
						{X: float64(j * consts.CellWidth), Y: float64(i * consts.CellWidth)},
						{X: float64(j * consts.CellWidth), Y: float64((i + 1) * consts.CellWidth)},
					},
				})
			}

		}
	}
}

func (m *Map) removeWall(p1, p2 point, direction int) {
	if direction == 0 {
		m.grid[p1.x][p1.y].n = false
		m.grid[p2.x][p2.y].s = false
	} else if direction == 1 {
		m.grid[p1.x][p1.y].s = false
		m.grid[p2.x][p2.y].n = false
	} else if direction == 2 {
		m.grid[p1.x][p1.y].e = false
		m.grid[p2.x][p2.y].w = false
	} else if direction == 3 {
		m.grid[p1.x][p1.y].w = false
		m.grid[p2.x][p2.y].e = false
	} else if direction == -1 {
		m.grid[p1.x][p1.y].n = false
		m.grid[p2.x][p2.y].s = false
		m.grid[p1.x][p1.y].s = false
		m.grid[p2.x][p2.y].n = false
		m.grid[p1.x][p1.y].e = false
		m.grid[p2.x][p2.y].w = false
		m.grid[p1.x][p1.y].w = false
		m.grid[p2.x][p2.y].e = false
	}
}

func (m *Map) subdivise(n int) [][]cell {
	nm := make([][]cell, consts.MapWidth*n)
	for i := range nm {
		nm[i] = make([]cell, consts.MapWidth*n)
	}

	for i, row := range m.grid {
		for j, c := range row {
			for k := 0; k < n; k++ {
				for l := 0; l < n; l++ {
					nm[i*n+k][j*n+l] = cell{
						n: c.n && k == 0,
						s: c.s && k == n-1,
						e: c.e && l == n-1,
						w: c.w && l == 0,
					}
				}
			}
		}
	}

	return nm
}

func (m *Map) countWallsInSubsquares(n int) {
	m.discreteGrid = make([][]uint8, (consts.MapWidth / n))
	for i := 0; i < (consts.MapWidth / n); i++ {
		m.discreteGrid[i] = make([]uint8, (consts.MapWidth / n))
	}

	for i := 0; i < consts.MapWidth; i += n {
		for j := 0; j < consts.MapWidth; j += n {
			count := uint8(0)
			for k := i; k < n+i && k < consts.MapWidth; k++ {
				for l := j; l < n+j && l < consts.MapWidth; l++ {
					if m.grid[k][l].n && k == i {
						count++
					}
					if m.grid[k][l].s {
						count++
					}
					if m.grid[k][l].e {
						count++
					}
					if m.grid[k][l].w && l == j {
						count++
					}
				}
			}
			m.discreteGrid[i/n][j/n] = count
		}
	}
}

type item struct {
	pos      point
	priority int
	index    int
}

type priorityQueue []*item

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (m *Map) dijkstra(start point, grid [][]cell) [][]int {
	height := len(grid)
	width := len(grid[0])

	dist := make([][]int, height)
	for i := range dist {
		dist[i] = make([]int, width)
		for j := range dist[i] {
			dist[i][j] = math.MaxInt32
		}
	}

	dist[start.x][start.y] = 0

	pq := make(priorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &item{start, 0, 0})

	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(*item)
		pos := curr.pos
		currentDist := curr.priority

		if currentDist > dist[pos.x][pos.y] {
			continue
		}

		for i, dir := range directions {
			newPos := point{pos.x + dir.x, pos.y + dir.y}

			if newPos.x < 0 || newPos.x >= height || newPos.y < 0 || newPos.y >= width {
				continue
			}

			if !grid[pos.x][pos.y].isWall(i) {
				newDist := currentDist + 1
				if newDist < dist[newPos.x][newPos.y] {
					dist[newPos.x][newPos.y] = newDist
					heap.Push(&pq, &item{newPos, newDist, 0})
				}
			}
		}
	}

	return dist
}

func (m *Map) getSpawnPoints(distances [][]int, min int, focusedRange int) {
	points := map[int][]*model.Point{}
	m.spawns[0] = []*model.Point{}

	var isLimit = func(n int) bool {
		m := n % int(consts.NumSubsquare)
		return m == 0 || m == int(consts.NumSubsquare)-1
	}

	for i := 0; i < consts.MapWidth; i++ {
		for j := 0; j < consts.MapWidth; j++ {
			center := &model.Point{
				X: float64(j*consts.CellWidth + consts.CellWidth/2),
				Y: float64(i*consts.CellWidth + consts.CellWidth/2),
			}

			for _, dir := range directions {
				x := center.X + float64(dir.x*consts.PlayerSize)*1.5
				y := center.Y + float64(dir.y*consts.PlayerSize)*1.5

				m.spawns[0] = append(m.spawns[0], &model.Point{X: x, Y: y})
			}

			m.spawns[0] = append(m.spawns[0], center)
		}
	}

	for i, row := range distances {
		for j, dist := range row {
			if isLimit(i) || isLimit(j) {
				continue
			}

			if _, ok := points[dist]; !ok {
				points[dist] = []*model.Point{}
			}

			points[dist] = append(points[dist], &model.Point{
				X: float64(j) * consts.SubsquareWidth,
				Y: float64(i) * consts.SubsquareWidth,
			})
		}
	}

	positions := make([]*model.Point, 0, min)

	for i := focusedRange - 1; i <= focusedRange+1 || len(positions) < min; i++ {
		if pts, ok := points[i]; ok {
			positions = append(positions, pts...)
		}
	}

	m.spawns[1] = positions
}

func (m *Map) Setup() {
	spawns := 0
	m.size = consts.MapWidth

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for spawns < 40 {
		grid := make([][]cell, consts.MapWidth)
		for i := range grid {
			grid[i] = make([]cell, consts.MapWidth)
			for j := range grid[i] {
				grid[i][j] = cell{true, true, true, true}
			}
		}
		m.grid = grid

		start := point{rand.Intn(m.size), rand.Intn(m.size)}
		m.primGenerateMaze(start)
		m.generateColliders()

		m.countWallsInSubsquares(2)

		m.start = model.Point{X: float64(start.x*consts.CellWidth + consts.CellWidth/2), Y: float64(start.y*consts.CellWidth + consts.CellWidth/2)}
		distances := m.dijkstra(point{x: start.x * consts.NumSubsquare, y: start.y * consts.NumSubsquare}, m.subdivise(consts.NumSubsquare))
		m.getSpawnPoints(distances, 40, 40)
		spawns = len(m.spawns[1])
	}

	utils.Shuffle(r, m.spawns[0])
	utils.Shuffle(r, m.spawns[1])
}

func (m *Map) Colliders() []*model.Collider {
	return m.walls
}

func (m *Map) Spawns(phase int) []*model.Point {
	return m.spawns[phase]
}

func (m *Map) Size() int {
	return m.size
}

func (m *Map) DiscreteMap() [][]uint8 {
	return m.discreteGrid
}

func (m *Map) Encode(w codec.Writer, withWalls bool) error {
	w.WriteInt8(int8(len(m.discreteGrid)))

	for _, row := range m.discreteGrid {
		for _, cell := range row {
			w.WriteUint8(cell)
		}
	}

	if !withWalls {
		return w.WriteInt32(0)
	}

	w.WriteInt32(int32(len(m.walls)))

	for _, wall := range m.walls {
		wall.Encode(w)
	}

	return nil
}

func (m *Map) Decode(r codec.Reader) error {
	size, err := r.ReadInt8()
	if err != nil {
		return err
	}
	m.size = int(size)

	m.discreteGrid = make([][]uint8, m.size)
	for i := 0; i < m.size; i++ {
		m.discreteGrid[i] = make([]uint8, m.size)
		for j := 0; j < m.size; j++ {
			m.discreteGrid[i][j], err = r.ReadUint8()
			if err != nil {
				return err
			}
		}
	}

	wallsLen, err := r.ReadInt32()
	if err != nil {
		return err
	}

	m.walls = make([]*model.Collider, wallsLen)
	for i := 0; i < int(wallsLen); i++ {
		m.walls[i] = &model.Collider{}
		if err = m.walls[i].Decode(r); err != nil {
			return err
		}
	}

	return nil
}
