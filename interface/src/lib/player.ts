import Phaser from 'phaser';
import { PLAYER_SIZE, PLAYER_SPEED } from '../config';


class Player extends Phaser.GameObjects.Container {
  private dest: Phaser.Math.Vector2;
  private name_graphics: Phaser.GameObjects.Text;
  private rect_graphics: Phaser.GameObjects.Rectangle;
  private username: string;

  constructor(scene: Phaser.Scene, x: number, y: number, name: string, color: number) {
    const rect_el = scene.add.rectangle(0, 0, PLAYER_SIZE, PLAYER_SIZE, color);
    const name_el = scene.add.text(0, -PLAYER_SIZE / 2 - 10, name, { fontSize: '16px', color: 'red' }).setOrigin(0.5);

    super(scene, x, y, [rect_el, name_el]);
  
    this.rect_graphics = rect_el;
    this.name_graphics = name_el;
    this.username = name;
    this.dest = new Phaser.Math.Vector2(x, y);
  }

  public get id(): string { return this.username };

  public move(dt: number): void {
    const distance = PLAYER_SPEED * (dt / 1000);
    const direction = this.dest.clone().subtract(new Phaser.Math.Vector2(this.x, this.y)).normalize();
    const movement = direction.scale(distance);

    if (Phaser.Math.Distance.Between(this.x, this.y, this.dest.x, this.dest.y) > distance) {
      this.x += movement.x;
      this.y += movement.y;
    } else {
    this.setPosition(this.dest.x, this.dest.y);
    }
  }

  public set_movement(pos: Phaser.Math.Vector2, dest: Phaser.Math.Vector2): void {
    this.dest = dest;
    this.setPosition(pos.x, pos.y);
  }
}

export { Player }; 
