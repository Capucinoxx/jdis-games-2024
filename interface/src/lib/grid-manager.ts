import Phaser from 'phaser';


class GridManager {
  private scene: Phaser.Scene;
  private grid_graphics: Phaser.GameObjects.Graphics;
  private walls: Array<Array<Position>> | null; 

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.walls = null;
    this.grid_graphics = this.scene.add.graphics();
  }

  public set map(values: { cells: number[][], colliders: Array<Array<Position>>}) {
    this.walls = values.colliders;
    this.clear();
    this.draw_walls();
  }

  private draw_walls() {
    if (!this.walls)
      return;

    this.grid_graphics.lineStyle(2, 0x050505, 1);
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
