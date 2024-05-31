import Phaser from 'phaser';
import { GameManager } from '../../lib';


class MainScene extends Phaser.Scene {
  private manager: GameManager | undefined;

  constructor() { super({ key: 'MainScene' }); }

  create() {
    this.manager = new GameManager(this);

    this.time.addEvent({
      delay: 1000,
      callback: this.test_player_movement,
      callbackScope: this,
      loop: true
    });

    this.cameras.main.setBounds(0, 0, 8000, 6000);
    this.physics.world.setBounds(0, 0, 8000, 6000);

    this.input.on('wheel', (pointer: Phaser.Input.Pointer): void => {
      const dy = pointer.deltaY;
      
      let new_zoom = this.cameras.main.zoom + (dy > 0 ? -.1 : .1);
      if (new_zoom > 0.3 && new_zoom <= 1)
        this.cameras.main.zoom = new_zoom;
    });

    this.input.on('pointerdown', (pointer: Phaser.Input.Pointer): void => {
      this.cameras.main.pan(pointer.x, pointer.y, 500);
    });

    window.addEventListener('resize', () => this.scale.resize(window.innerWidth, window.innerHeight));
  }


  /**
   * TODO: to be removed when testing connectivity with the server.
   */
  test_player_movement() {
    const calculate_next_position = (curr_x: number, curr_y: number): [number, number] => {
      let x_options = [0, 100, -100].filter(opt => (curr_x + opt >= 0) && (curr_x + opt <= 200));
      let y_options = [0, 100, -100].filter(opt => (curr_y + opt >= 0) && (curr_y + opt <= 200));
    
      let new_x = curr_x + x_options[Math.floor(Math.random() * x_options.length)];
      let new_y = curr_y + y_options[Math.floor(Math.random() * y_options.length)];
    
      return [new_x, new_y];
    };
  
    const players = ['1234', '5678'].map(uuid => {
      const player = this.manager!.get_player(uuid);
      let new_x = this.cameras.main.centerX - (2 * 200) / 2 + 100, new_y = this.cameras.main.centerY - (2 * 100) / 2;
      

      if (player) {
        const { x, y } = player.target_pos;
        [new_x, new_y] = calculate_next_position(x, y);
      }
                
    
      return { 'uuid': uuid, 'x': new_x, 'y': new_y };
    }).filter(player => player !== null);
  
    this.manager!.update_from_payload({ 'players': players });
  }

  update(_: number, delta: number) {
    this.manager!.update_players_movement(delta);
  }
};

export { MainScene };
