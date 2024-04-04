import { Container, ObservablePoint } from 'pixi.js';
import { clamp } from '../utils';

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

  /**
   * update Met à jours la position des containers pour faire
   * en sorte que la source soit au centre de l'écran.
   */
  public update(): void {
    let targetX = -this.source.position.x + this.screenWidth / 2;
    let targetY = -this.source.position.y + this.screenHeight / 2;

    targetX = clamp(targetX, -(this.mapWidth - this.screenWidth), 0);
    targetY = clamp(targetY, -(this.mapHeight - this.screenHeight), 0);

    this.containers.forEach(container => {
      container.position.set(targetX, targetY);
    });
  }

  public set focus(source: Position) {
    this.source = source;
  }
}

export { Camera };
