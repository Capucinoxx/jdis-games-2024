import Phaser from 'phaser';
import { CELL_WIDTH } from '../config';


class GridManager {
  private scene: Phaser.Scene;
  private grid_graphics: Phaser.GameObjects.Graphics;
  private walls: Array<Array<Position>> | null; 
  private cells: number[][] | null;

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.walls = null;
    this.cells = null;
    this.grid_graphics = this.scene.add.graphics();
  }

  public set map(values: { cells: number[][], colliders: Array<Array<Position>>}) {
    this.walls = values.colliders;
    this.cells = values.cells;
    this.clear();
    this.draw_walls();
  }

  private draw_walls() {
    this.grid_graphics.lineStyle(2, 0x050505, 1);

    if ((!this.walls || this.walls.length === 0) && this.cells) {
      const width = this.cells.length * CELL_WIDTH;
      this.grid_graphics.beginPath();
      this.grid_graphics.moveTo(0, 0);
      this.grid_graphics.lineTo(width, 0);
      this.grid_graphics.lineTo(width, width);
      this.grid_graphics.lineTo(0, width);
      this.grid_graphics.lineTo(0, 0);
      this.grid_graphics.closePath();
      this.grid_graphics.strokePath();
    
      return;
    }

    if (!this.walls)
      return;

    
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
