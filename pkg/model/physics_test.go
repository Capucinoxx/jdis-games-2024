package model

import "testing"

func TestPolygonsIntersect(t *testing.T) {
	tests := map[string]struct {
		a, b Polygon
		want bool
	}{
		"Polygons intersect": {
			a:    Polygon{vertices: []*Point{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}}},
			b:    Polygon{vertices: []*Point{{X: 2, Y: 2}, {X: 6, Y: 2}, {X: 6, Y: 6}, {X: 2, Y: 6}}},
			want: true,
		},
		"Polygons don't intersect": {
			a:    Polygon{vertices: []*Point{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}}},
			b:    Polygon{vertices: []*Point{{X: 3, Y: 3}, {X: 5, Y: 3}, {X: 5, Y: 5}, {X: 3, Y: 5}}},
			want: false,
		},
		"Polygon and line intersect": {
			a:    Polygon{vertices: []*Point{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}}},
			b:    Polygon{vertices: []*Point{{X: 2, Y: 2}, {X: 6, Y: 2}}},
			want: true,
		},
		"Polygon and line don't intersect": {
			a:    Polygon{vertices: []*Point{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}}},
			b:    Polygon{vertices: []*Point{{X: 5, Y: 5}, {X: 7, Y: 5}}},
			want: false,
		},
		"Line and line intersect": {
			a:    Polygon{vertices: []*Point{{X: 0, Y: 0}, {X: 4, Y: 4}}},
			b:    Polygon{vertices: []*Point{{X: 0, Y: 4}, {X: 4, Y: 0}}},
			want: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if PolygonsIntersect(tt.a, tt.b) != tt.want {
				t.Errorf("PolygonsIntersect() != %v", tt.want)
			}
		})
	}
}
