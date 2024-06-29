import Phaser from 'phaser';

class MovableObject extends Phaser.GameObjects.Container {
  private speed: number;
  protected destination: Phaser.Math.Vector2;

  constructor(scene: Phaser.Scene, x: number, y: number, destination: Phaser.Math.Vector2, speed: number, deps: Array<Phaser.GameObjects.GameObject> = []) {
    super(scene, x, y, deps);
    this.speed = speed;
    this.destination = destination;

    this.scene.physics.world.enable(this);
    scene.add.existing(this);
  }


  public move(dt: number): void {
    const dist = this.speed * (dt / 1000);
    const direction = new Phaser.Math.Vector2(this.destination.x - this.x, this.destination.y - this.y);
    const cur_dist = direction.length();

    direction.normalize();

    if (cur_dist > dist) {
      this.x += direction.x * dist;
      this.y += direction.y * dist;
    } else {
      this.setPosition(this.destination.x, this.destination.y);
      this.on_arrival();
    }
  };

  protected on_arrival(): void {
    this.destroy();
  }
};

export { MovableObject };
