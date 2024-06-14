import Phaser from 'phaser';
import { CELL_WIDTH as TILE_WIDTH } from './config';

const TILE_COLORS = [
  0x5D5D5D,
  0x4C4C4C,
  0x3C3C3C,
  0x2E2E2E,
  0x222222
];

const GRID_COLOR = 0x151515;


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


    this.grid_graphics.clear();
    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {

        const color = TILE_COLORS[this.grid_values[y][x]] || 0xff0000;
        this.grid_graphics.fillStyle(color, 1.0);
        this.grid_graphics.fillPoints([
          { x: x * TILE_WIDTH, y: y * TILE_WIDTH },
          { x: (x + 1) * TILE_WIDTH, y: y * TILE_WIDTH },
          { x: (x + 1) * TILE_WIDTH, y: (y + 1) * TILE_WIDTH },
          { x: x * TILE_WIDTH, y: (y + 1) * TILE_WIDTH },
          { x: x * TILE_WIDTH, y: y * TILE_WIDTH },
        ], true);
      }
    }
  }
};

export { GridManager };
