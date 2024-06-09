import Phaser from 'phaser';

const BULLET_RADIUS = 10;
const BULLET_SPEED = 10;

class Bullet extends Phaser.GameObjects.Graphics {
  private destination: Phaser.Math.Vector2;
  
  constructor(scene: Phaser.Scene, x: number, y: number, destination: { x: number, y: number }) {
    super(scene);
    this.scene = scene;
    this.destination = new Phaser.Math.Vector2(destination.x, destination.y);

    this.setPosition(x, y);
    this.fillStyle(0x00ff00, 1);

    this.fillCircle(0, 0, BULLET_RADIUS);
    this.scene.add.existing(this);
    scene.physics.world.enable(this);
  }

  update(time: number, delta: number): void {
    const angle = Phaser.Math.Angle.Between(this.x, this.y, this.destination.x, this.destination.y);
    (this.body as Phaser.Physics.Arcade.Body).setVelocity(Math.cos(angle) * BULLET_SPEED, Math.sin(angle) * BULLET_SPEED);

    if (Phaser.Math.Distance.Between(this.x, this.y, this.destination.x, this.destination.y) < BULLET_RADIUS)
      this.destroy_bullet();
  }

  destroy_bullet() { this.destroy(); }
}

export { Bullet };
