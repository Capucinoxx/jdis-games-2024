import Phaser from 'phaser';
import { MainScene } from './scenes';
import { load_wasm } from './wasm/loader';
import { handle_modals } from './modal';
import { handle_forms } from './form';
import { generate_particles } from './particle';

generate_particles('particles');
handle_forms();

const config: Phaser.Types.Core.GameConfig = {
  type: Phaser.AUTO,
  width: 800,
  height: 800,
  backgroundColor: '#2a2b38',
  parent: 'game',
  physics: {
    default: 'arcade',
    arcade: {
      debug: false
    }
  },
  scene: [MainScene],
};

// load_wasm().then(() => {
  const game = new Phaser.Game(config);
  handle_modals(game);
// });
