import { Graphics, PointData } from 'pixi.js';
import { Vector } from './position';

/**
 * Player représente un joueur dans le jeu.
 */
class Player {
  static readonly RADIUS: number = 20;

  private name: string;
  private color: string;
  private sides: number;
  private destination: Vector;
  private speed: number;
  public graphics: Graphics;

  constructor(name: string, color: string, sides: number, pos: Vector) {
    this.name = name;
    this.color = color;
    this.sides = sides;
    this.speed = 1;
    this.graphics = new Graphics();
    this.destination = pos;

    this.draw();
  }

  /**
   * Dessine le polygone représentant le joueur.
   */
  private draw(): void {
    this.graphics.clear();
    this.graphics.poly(this.points(), false);
    this.graphics.fill(this.color);
  }

  /**
   * @returns Les points du polygone représentant le joueur.
   */
  public points(): Vector[] {
    let points: Vector[] = [];
    for (let i = 0; i < this.sides; i++) {
      const angle = Math.PI * 2 * i / this.sides;
      points.push(new Vector(Math.cos(angle) * Player.RADIUS, Math.sin(angle) * Player.RADIUS));
    }
    return points;
  }

  /**
   * Définit la destination du joueur.
   * @param v {Vector} La destination du joueur.
   */
  public set_destination(v: Vector): void {
    this.destination = v;
  }

  /**
   * Met à jour la position du joueur. On déplace le joueur vers sa destination
   * selon le temps écoulé depuis la dernière mise à jour.
   * 
   * @param dt {number} Le temps écoulé depuis la dernière mise à jour.
   */
  public update(dt: number): void {
    const dx = this.destination.x - this.graphics.x;
    const dy = this.destination.y - this.graphics.y;
    const distance = Math.sqrt(dx * dx + dy * dy);
    const speed = Math.min(distance, this.speed * dt);

    this.graphics.x += (dx / distance * speed) || 0;
    this.graphics.y += (dy / distance * speed) || 0;
  }
};