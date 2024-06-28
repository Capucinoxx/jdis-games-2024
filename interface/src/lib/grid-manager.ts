import Phaser from 'phaser';
import { CELL_WIDTH as TILE_WIDTH } from '../config';

const GRID_COLOR = 0x151515;
const WALL_COLOR = 0x00ff00;


class GridManager {
  private scene: Phaser.Scene;
  private grid_graphics: Phaser.GameObjects.Graphics;
  private grid_values: number[][] | null;
  private walls: Array<Array<Position>> | null; 

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.grid_values = null;
    this.walls = null;
    this.grid_graphics = this.scene.add.graphics();
  }

  public set map(values: { cells: number[][], colliders: Array<Array<Position>>}) {
    this.grid_values = values.cells;
    this.walls = values.colliders;
    this.draw_grid();
    this.draw_walls();
  }


  private draw_grid(): void {
    if (!this.grid_values)
      return;

    const rows = this.grid_values.length;
    const cols = this.grid_values[0].length;


    this.grid_graphics.clear();
    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {
/**
        this.grid_graphics.fillStyle(0x00a759, this.grid_values[y][x] * 0.05);
        this.grid_graphics.fillPoints([
          { x: x * TILE_WIDTH, y: y * TILE_WIDTH },
          { x: (x + 1) * TILE_WIDTH, y: y * TILE_WIDTH },
          { x: (x + 1) * TILE_WIDTH, y: (y + 1) * TILE_WIDTH },
          { x: x * TILE_WIDTH, y: (y + 1) * TILE_WIDTH },
          { x: x * TILE_WIDTH, y: y * TILE_WIDTH },
        ], true);
        */
      }
      
    }
  }

  private draw_walls() {
    if (!this.walls)
      return;

    this.grid_graphics.lineStyle(1, 0x050505, 1);
    this.walls.forEach((wall) => {
      if (wall.length != 2)
        return;

      this.grid_graphics.beginPath();
      this.grid_graphics.moveTo(wall[0].x, wall[0].y);
      this.grid_graphics.lineTo(wall[1].x, wall[1].y);
      this.grid_graphics.closePath();
      this.grid_graphics.strokePath();
    });
  }

  public clear(): void {
    this.grid_graphics.clear();
  }
};

export { GridManager };
