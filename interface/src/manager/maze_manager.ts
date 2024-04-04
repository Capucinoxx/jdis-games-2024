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

  /**
   * set Met à jour la liste de polygones à afficher. Cela
   * force le re-rendu des polygones. 
   * 
   * @param colliders {Polygon[]} Les polygones à afficher.
   */
  public set(colliders: Polygon[]): void {
    this.colliders = colliders;
    this.draw();
  }

  /**
   * draw Force le rendu des polygones.
   */
  public draw(): void {
    this.container.removeChildren();
    this.colliders.forEach(collider => {
      this.container.addChild(collider.graphics);
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