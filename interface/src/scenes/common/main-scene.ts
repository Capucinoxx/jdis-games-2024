import Phaser from 'phaser';
import { GameManager } from '../../lib';

class MainScene extends Phaser.Scene {
  private manager: GameManager;
  private grid_graphics: Phaser.GameObjects.Graphics;

  constructor() { super({ key: 'MainScene' }); }

  create() {
    this.manager = new GameManager(this, this.cameras.main.centerX);
    this.draw_iso_grid();
    this.draw_grid();

    this.time.addEvent({
      delay: 1000,
      callback: this.test_player_movement,
      callbackScope: this,
      loop: true
    });

    this.cameras.main.setBounds(0, 0, 800, 600);
    this.physics.world.setBounds(0, 0, 800, 600);

    this.input.on('wheel', (pointer: Phaser.Input.Pointer): void => {
      const dy = pointer.deltaY;
      
      let new_zoom = this.cameras.main.zoom + (dy > 0 ? -.1 : .1);
      if (new_zoom > 0.3 && new_zoom <= 1)
        this.cameras.main.zoom = new_zoom;
    });
  }

  /**
   * TODO: to be removed when the tiles are implemented.
   */
  draw_grid() {
    const grid_size = 100;
    const rows = 6;
    const cols = 8;
    this.grid_graphics = this.add.graphics({ lineStyle: { width: 2, color: 0x0000ff } });

    for (let i = 0; i <= cols; i++)
      this.grid_graphics.lineBetween(i * grid_size, 0, i * grid_size, 600);
    for (let j = 0; j <= rows; j++)
      this.grid_graphics.lineBetween(0, j * grid_size, 800, j * grid_size);
  }

  /**
   * TODO: to be removed when the tiles are implemented.
   */
  draw_iso_grid() {
    const tile_width = 200;
    const tile_height = 100;

    const rows = 2;
    const cols = 2;

    this.grid_graphics = this.add.graphics();
    this.grid_graphics.lineStyle(1, 0x0000ff, 0.5);

    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {
        const base_x = this.cameras.main.centerX - (cols * tile_width) / 2;
        const base_y = this.cameras.main.centerY - (rows * tile_height) / 2;

        let ix = base_x + (x - y) * tile_width / 2;
        let iy = base_y + (x + y) * tile_height / 2;

        this.grid_graphics.strokePoints([
          { x: ix, y: iy + tile_height / 2 },
          { x: ix + tile_width / 2, y: iy },
          { x: ix + tile_width, y: iy + tile_height / 2 },
          { x: ix + tile_width / 2, y: iy + tile_height },
          { x: ix, y: iy + tile_height / 2 }
      ], true);
      }
    }
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
      const player = this.manager.get_player(uuid);
      let new_x = this.cameras.main.centerX - (2 * 200) / 2 + 100, new_y = this.cameras.main.centerY - (2 * 100) / 2;
      

      if (player) {
        const { x, y } = player.target_pos;
        [new_x, new_y] = calculate_next_position(x, y);
      }
        
        
    
      return { 'uuid': uuid, 'x': new_x, 'y': new_y };
    }).filter(player => player !== null);
  
    this.manager.update_from_payload({ 'players': players });
  }

  update(_: number, delta: number) {
    this.manager.update_players_movement(delta);
  }
};

export { MainScene };
