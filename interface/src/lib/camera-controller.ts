class CameraController {
  private camera: Phaser.Cameras.Scene2D.Camera;
  private cursors: Phaser.Types.Input.Keyboard.CursorKeys;
  private wasd: { [key: string]: Phaser.Input.Keyboard.Key };
  private shift: Phaser.Input.Keyboard.Key;
  private current_following: Phaser.GameObjects.GameObject | null = null;

  constructor(camera: Phaser.Cameras.Scene2D.Camera, input: Phaser.Input.InputPlugin) {
    if (!input.keyboard)
      throw new Error('');

    this.camera = camera;
    this.camera.setZoom(1);
    
    this.cursors = input.keyboard.createCursorKeys();
    this.wasd = input.keyboard.addKeys('W,S,A,D') as { [key: string]: Phaser.Input.Keyboard.Key };
    this.shift = input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.SHIFT);

    input.on('wheel', (_: Phaser.Input.Pointer, __: any, ___: number, dy: number, ____: number) => {
      if (dy > 0) this.camera.zoom = Math.max(0.5, this.camera.zoom - 0.1);
      else if (dy < 0) this.camera.zoom = Math.min(1.1, this.camera.zoom + 0.1);
    });
  }

  public follow(target: Phaser.GameObjects.GameObject) {
    this.camera.startFollow(target);
    this.current_following = target;
  }

  public unfollow() {
    if (this.current_following) {
      this.camera.stopFollow();
      this.current_following = null;
    }
  }

  public update(dt: number) {
    const speed = 500 * (this.shift.isDown ? 2 : 1);
    const movement = speed * (dt / 1000);
    let on_movement = false;

    if (this.cursors.left.isDown || this.wasd.A.isDown) {
      this.camera.scrollX -= movement;
      this.suspend_following();
      on_movement = true;
    } else if (this.cursors.right.isDown || this.wasd.D.isDown) {
      this.camera.scrollX += movement;
      this.suspend_following();
      on_movement = true;
    }

    if (this.cursors.up.isDown || this.wasd.W.isDown) {
      this.camera.scrollY -= movement;
      this.suspend_following();
      on_movement = true;
    } else if (this.cursors.down.isDown || this.wasd.S.isDown) {
      this.camera.scrollY += movement;
      this.suspend_following();
      on_movement = true;
    }

    if (this.any_key_justup() && !on_movement && this.current_following !== null)
      this.follow(this.current_following);
  }

  private suspend_following() {
    this.camera.stopFollow();
  }

  private any_key_justup(): boolean {
    return Phaser.Input.Keyboard.JustUp(this.cursors.left) || Phaser.Input.Keyboard.JustUp(this.cursors.right) ||
           Phaser.Input.Keyboard.JustUp(this.cursors.up) || Phaser.Input.Keyboard.JustUp(this.cursors.down) ||
           Phaser.Input.Keyboard.JustUp(this.wasd.A) || Phaser.Input.Keyboard.JustUp(this.wasd.D) ||
           Phaser.Input.Keyboard.JustUp(this.wasd.W) || Phaser.Input.Keyboard.JustUp(this.wasd.S);
  }
};

export { CameraController };