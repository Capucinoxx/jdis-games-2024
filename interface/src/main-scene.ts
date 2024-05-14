import Phaser from 'phaser';
import { GameManager } from './game-manager';

class MainScene extends Phaser.Scene {
  private manager: GameManager;
  private grid_graphics: Phaser.GameObjects.Graphics;

  constructor() { super({ key: 'MainScene' }); }

  create() {
    this.manager = new GameManager(this);
    this.draw_grid();

    this.time.addEvent({
      delay: 300,
      callback: this.test_player_movement,
      callbackScope: this,
      loop: true
    });

    this.cameras.main.setBounds(0, 0, 800, 600);
    this.physics.world.setBounds(0, 0, 800, 600);
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
   * TODO: to be removed when testing connectivity with the server.
   */
  test_player_movement() {
    const calculate_next_position = (curr_x: number, curr_y: number): [number, number] => {
      let x_options = [0, 100, -100].filter(opt => (curr_x + opt >= 0) && (curr_x + opt <= 600));
      let y_options = [0, 100, -100].filter(opt => (curr_y + opt >= 0) && (curr_y + opt <= 600));
    
      let new_x = curr_x + x_options[Math.floor(Math.random() * x_options.length)];
      let new_y = curr_y + y_options[Math.floor(Math.random() * y_options.length)];
    
      return [new_x, new_y];
    };
  
    const players = ['1234', '5678'].map(uuid => {
      const player = this.manager.get_player(uuid);
      let new_x = 0, new_y = 0;

      if (player)
        [new_x, new_y] = calculate_next_position(player.target_x, player.target_y);
    
      return { 'uuid': uuid, 'x': new_x, 'y': new_y };
    }).filter(player => player !== null);
  
    this.manager.update_from_payload({ 'players': players });
  }

  update(_: number, delta: number) {
    this.manager.update_players_movement(delta);
  }
};

export { MainScene };
