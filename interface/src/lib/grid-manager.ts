import Phaser from 'phaser';
import { CELL_WIDTH, NUM_CELLS } from '../config';


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
    this.draw_background();
  }

  public set map(values: { cells: number[][], colliders: Array<Array<Position>>}) {
    this.walls = values.colliders;
    this.cells = values.cells;
    this.clear();
    this.draw_walls();
  }

  private draw_background() {
    let graphics = this.scene.make.graphics();
    const width = CELL_WIDTH * NUM_CELLS;

    for (let i = 0; i <= width; i += CELL_WIDTH) {
        if (i === 0 || i === width) {
            graphics.lineStyle(5, 0xc4c3ca, 0.95);
        } else if ((i / CELL_WIDTH) % 2 === 0) {
            graphics.lineStyle(3, 0xc4c3ca, 0.6);
        } else {
            graphics.lineStyle(1, 0xc4c3ca, 0.1);
        }
        graphics.beginPath();
        graphics.moveTo(i, 0);
        graphics.lineTo(i, width);
        graphics.closePath();
        graphics.strokePath();
    }

    for (let j = 0; j <= width; j += CELL_WIDTH) {
        if (j === 0 || j === width) {
          graphics.lineStyle(5, 0xc4c3ca, 0.95);
        } else if ((j / CELL_WIDTH) % 2 === 0) {
          graphics.lineStyle(3, 0xc4c3ca, 0.6); 
        } else {
          graphics.lineStyle(1, 0xc4c3ca, 0.1);
        }
        graphics.beginPath();
        graphics.moveTo(0, j);
        graphics.lineTo(width, j);
        graphics.closePath();
        graphics.strokePath();
    }

    graphics.generateTexture('grid', width, width);
    const sprite = this.scene.add.sprite(0, 0, 'grid');  
    sprite.setOrigin(0, 0);
  }

  private draw_walls() {
    if (!this.walls)
      return;

    this.grid_graphics.lineStyle(2, 0xbad7f7, 1);
    
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
