import Phaser from 'phaser';

export class Player extends Phaser.GameObjects.Container {
    public uuid: string;
    public speed: number;
    public target_x: number;
    public target_y: number;
    private graphics: Phaser.GameObjects.Graphics;
    private body_sprite: Phaser.Physics.Arcade.Sprite;
    
    constructor(scene: Phaser.Scene, x: number, y: number, uuid: string) {
        super(scene, x, y);

        this.uuid = uuid;
        this.speed = 100;
        this.target_x = x;
        this.target_y = y;

        this.body_sprite = this.scene.physics.add.sprite(0, 0, this.uuid);
        this.body_sprite.setVisible(false);

        this.graphics = new Phaser.GameObjects.Graphics(this.scene);
        this.graphics.fillStyle(0xff0000);
        this.graphics.fillCircle(0, 0, 5);
        this.add(this.graphics);

        this.scene.add.existing(this);
        this.scene.physics.add.existing(this);

        this.body = this.body_sprite.body;
        this.body!.setCircle(5);
        this.body_sprite.setPosition(x, y);
    }

    public update_position(x: number, y: number): void {
        this.body_sprite.setPosition(this.target_x, this.target_y);
        this.target_x = x;
        this.target_y = y;
    }

    public move_to_target(delta: number): void {
      if (!(this.body instanceof Phaser.Physics.Arcade.Body))
        return;

      const distance = Phaser.Math.Distance.Between(this.x, this.y, this.target_x, this.target_y);
      if (distance < 1) {
        this.body.setVelocity(0, 0);
        return;
      }

      const angle = Phaser.Math.Angle.Between(this.x, this.y, this.target_x, this.target_y);
      const move_distance = Math.min(distance, this.speed * delta / 1000 * 3.33);
      this.body_sprite.x += Math.cos(angle) * move_distance;
      this.body_sprite.y += Math.sin(angle) * move_distance;
      this.x = this.body_sprite.x;
      this.y = this.body_sprite.y;

      this.body.setVelocityX(Math.cos(angle) * this.speed);
      this.body.setVelocityY(Math.sin(angle) * this.speed);
    }

    public get target_pos(): [number, number] {
        return [this.target_x, this.target_y];
    }
}
