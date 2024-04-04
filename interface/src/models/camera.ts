import { Container, ObservablePoint } from 'pixi.js';

interface Position {
  position: ObservablePoint;
}

class Camera {
  private containers: Container[];
  private source: Position;
  private mapWidth: number;
  private mapHeight: number;
  private screenWidth: number;
  private screenHeight: number;

  constructor(containers: Container[], source: Position, screenWidth: number, screenHeight: number, mapWidth: number, mapHeight: number) {
    this.containers = containers;
    this.source = source;
    this.screenWidth = screenWidth;
    this.screenHeight = screenHeight;
    this.mapWidth = mapWidth;
    this.mapHeight = mapHeight;
  }

  private clamp(value: number, min: number, max: number): number {
    return Math.max(min, Math.min(max, value));
  }

  public update(): void {
    let targetX = -this.source.position.x + this.screenWidth / 2;
    let targetY = -this.source.position.y + this.screenHeight / 2;

    targetX = this.clamp(targetX, -(this.mapWidth - this.screenWidth), 0);
    targetY = this.clamp(targetY, -(this.mapHeight - this.screenHeight), 0);

    console.log(targetX, targetY);

    this.containers.forEach(container => {
      container.position.set(targetX, targetY);
    });
  }

  public set focus(source: Position) {
    this.source = source;
  }
}

export { Camera };
