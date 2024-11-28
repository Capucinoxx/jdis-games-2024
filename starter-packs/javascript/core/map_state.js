const ColliderType = Object.freeze({
  Wall: 0,
  Projectile: 1,
  toString(value) {
    return Object.keys(this).find((key) => this[key] ===value) || '';
  },
});

class Point {
  constructor(x = 0.0, y = 0.0) {
    this.x = x;
    this.y = y;
  }

  toString() { return `{"x":${this.x},"y":${this.y}}`; }
};


class Collider {
  constructor(colliderType = ColliderType.Projectile, positions = []) {
    this.collider_type = colliderType;
    this.positions = positions;
  }

  toString() {
    const pos_str = this.positions.map((pos) => pos.toString()).join(",");
    return `{"collider_type":${this.collider_type},"positions":[${pos_str}]}`;
  }
};

class MapState {
  constructor(size = 0, discrete_grid = null, walls = [], save = null) {
    this.size = size;
    this.discrete_grid = discrete_grid || [];
    this.walls = walls;
    this.save = save || new Uint8Array(0);
  }

  toString() {
    const walls_str = this.walls.map((wall) => wall.toString()).join(",");
    const save_str = Array.from(this.save)
      .map((byte) => `0x${byte.toString(16).padStart(2, '0')}`)
      .join(' ');

    return `{"size": ${this.size}, "discrete_grid": ${JSON.stringify(this.discrete_grid)}, "walls": [${walls_str}], "save": ${save_str}}`;
  }
};

export { ColliderType, Point, Collider, MapState };
