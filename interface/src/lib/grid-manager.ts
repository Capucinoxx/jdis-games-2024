import Phaser from 'phaser';

const TILE_COLORS = [
  0x5D5D5D,
  0x4C4C4C,
  0x3C3C3C,
  0x2E2E2E,
  0x222222
];

const GRID_COLOR = 0x151515;

const TILE_WIDTH = 200;
const TILE_HEIGHT = 100;
const HALF_TILE_WIDTH = TILE_WIDTH / 2;
const HALF_TILE_HEIGHT = TILE_HEIGHT / 2;

class GridManager {
  private scene: Phaser.Scene;
  private grid_graphics: Phaser.GameObjects.Graphics;
  private grid_values: number[][] | null;

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.grid_values = null;
    this.grid_graphics = this.scene.add.graphics();
  }

  public set tiles(values: number[][]) {
    this.grid_values = values;
    this.draw_grid();
  }

  private draw_grid(): void {
    if (!this.grid_values)
      return;

    const rows = this.grid_values.length;
    const cols = this.grid_values[0].length;

    const base_x = this.scene.cameras.main.centerX - cols * HALF_TILE_WIDTH;
    const base_y = this.scene.cameras.main.centerY - rows * HALF_TILE_HEIGHT;

    this.grid_graphics.clear();
    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {
        const ix = base_x + (x - y) * HALF_TILE_WIDTH;
        const iy = base_y + (x + y) * HALF_TILE_HEIGHT;

        const color = TILE_COLORS[this.grid_values[y][x]] || 0xff0000;
        this.grid_graphics.fillStyle(color, 1.0);
        this.grid_graphics.fillPoints([
          { x: ix, y: iy + HALF_TILE_HEIGHT },
          { x: ix + HALF_TILE_WIDTH, y: iy },
          { x: ix + TILE_WIDTH, y: iy + HALF_TILE_HEIGHT },
          { x: ix + HALF_TILE_WIDTH, y: iy + TILE_HEIGHT },
          { x: ix, y: iy + HALF_TILE_HEIGHT },
        ], true);

        this.grid_graphics.lineStyle(1, GRID_COLOR, 1);
        this.grid_graphics.strokePoints([
          { x: ix, y: iy + HALF_TILE_HEIGHT },
          { x: ix + HALF_TILE_WIDTH, y: iy },
          { x: ix + TILE_WIDTH, y: iy + HALF_TILE_HEIGHT },
          { x: ix + HALF_TILE_WIDTH, y: iy + TILE_HEIGHT },
          { x: ix, y: iy + HALF_TILE_HEIGHT },
        ]);
      }
    }
  }
};

export { GridManager };