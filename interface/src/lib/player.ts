import Phaser from 'phaser';

import { IsometricCoordinates, CartesianCoordinates, cartesian_to_isometric, isometric_to_cartesian } from '.';

export class Player extends Phaser.GameObjects.Container {
    public uuid: string;
    public speed: number;
    private target_iso: IsometricCoordinates;
    private zero: CartesianCoordinates;
    private graphics: Phaser.GameObjects.Graphics;
    private body_sprite: Phaser.Physics.Arcade.Sprite;
    
    constructor(scene: Phaser.Scene, cart: CartesianCoordinates, uuid: string) {
        super(scene);

        this.uuid = uuid;
        this.speed = 100;
        this.target_iso = cart;
        this.zero = cart;

        this.body_sprite = this.scene.physics.add.sprite(cart.x, cart.y, uuid);
        this.body_sprite.setVisible(false);

        this.graphics = new Phaser.GameObjects.Graphics(this.scene);
        this.graphics.fillStyle(0xff0000);
        this.graphics.fillCircle(0, 0, 5);
        this.add(this.graphics);

        this.scene.add.existing(this);
        this.scene.physics.add.existing(this);

        this.body = this.body_sprite.body;
        this.body!.setCircle(5);
        this.body_sprite.setPosition(cart.x, cart.y);
    }

    public update_position(cart: CartesianCoordinates): void {
      let iso = cartesian_to_isometric(cart);
      iso.x += this.zero.x;
      iso.y += this.zero.y;
      
      this.body_sprite.setPosition(this.target_iso.x, this.target_iso.y);
      this.target_iso = iso;
    }

    public move_to_target(delta: number): void {
      if (!(this.body instanceof Phaser.Physics.Arcade.Body))
        return;

      const distance = Phaser.Math.Distance.Between(this.x, this.y, this.target_iso.x, this.target_iso.y);
      if (distance < 2) {
        this.body.setVelocity(0, 0);
        this.body_sprite.setPosition(this.target_iso.x, this.target_iso.y);
        return;
      }

      const angle = Phaser.Math.Angle.Between(this.x, this.y, this.target_iso.x, this.target_iso.y);
      const move_distance = Math.min(distance, this.speed * delta / 1000 * 3.33);
      this.body_sprite.x += Math.cos(angle) * move_distance;
      this.body_sprite.y += Math.sin(angle) * move_distance;

      this.body.setVelocityX(Math.cos(angle) * this.speed);
      this.body.setVelocityY(Math.sin(angle) * this.speed);

      this.x = this.body_sprite.x;
      this.y = this.body_sprite.y;
    }

    public get target_pos(): CartesianCoordinates {
      return isometric_to_cartesian({ 'x': this.target_iso.x - this.zero.x, 'y': this.target_iso.y - this.zero.y });
    }
}
