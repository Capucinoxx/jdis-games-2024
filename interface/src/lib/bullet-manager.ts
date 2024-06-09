import Phaser from 'phaser';
import { Bullet } from './bullet';
import '../types/index.d.ts';

class BulletManager {
  private scene: Phaser.Scene;
  private bullets: Phaser.Physics.Arcade.Group; 
  private cache: Map<string, Bullet>;

  constructor(scene: Phaser.Scene) {
    this.scene = scene;
    this.bullets = this.scene.physics.add.group({ classType: Bullet, runChildUpdate: true });
    this.cache = new Map<string, Bullet>();
  }

  public sync(payload: Projectile[]): void {
    const projectiles = new Map<string, Projectile>();
    payload.forEach((p) => projectiles.set(p.id, p));

    this.cache.forEach((bullet, uuid) => {
      if (!projectiles.has(uuid))
        bullet.destroy_bullet();
    });

    projectiles.forEach(({ pos, dest }, uuid) => {
      if (!this.cache.has(uuid)) {
        dest = { x: dest.x * 30, y: dest.y * 30 };
        const bullet = new Bullet(this.scene, pos.x * 30, pos.y * 30, dest);
        this.bullets.add(bullet);
        this.cache.set(uuid, bullet);
      }
    });
  }
};

export { BulletManager };
