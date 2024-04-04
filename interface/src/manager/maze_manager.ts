import { Application, Container } from 'pixi.js';
import { Polygon, Vector } from '../models/position';

class MazeManager {
  private container: Container;
  private colliders: Polygon[];

  private scale: number = 1;

  constructor(app: Application) {
    this.container = new Container();
    app.stage.addChild(this.container);
  }

  public set(colliders: Polygon[]): void {
    this.colliders = colliders;
    this.draw();
  }

  public draw(): void {
    this.container.removeChildren();
    this.colliders.forEach(collider => {
      this.container.addChild(collider.graphics);
    });
  }

  public multiply_scale(factor: number): void {
    this.scale *= factor;
    this.colliders.map(collider => {
      collider.graphics.scale.set(this.scale);
    });
  }

  public get view(): Container {
    return this.container;
  }

  public get size(): Vector {
    return new Vector(this.container.width, this.container.height);
  }
};

export { MazeManager };