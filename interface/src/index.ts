import Phaser from 'phaser';
import { MainScene } from './scenes';
import { load_wasm } from './wasm/loader';

const config: Phaser.Types.Core.GameConfig = {
  type: Phaser.AUTO,
  width: 800,
  height: 600,
  physics: {
    default: 'arcade',
    arcade: {
      debug: false
    }
  },
  scene: [MainScene]
};

load_wasm();

const game = new Phaser.Game(config);