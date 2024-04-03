import { Application, Container } from 'pixi.js';
import { Player } from '../models/player';

class PlayerManager {
  private container: Container;
  private players: Player[];

  constructor(app: Application) {
    this.players = [];

    this.container = new Container();
    app.stage.addChild(this.container);
  }

  /**
   * append Ajoute des joueurs à la liste des joueurs gérés par le gestionnaire.
   * Lorsque des joueurs sont ajoutés, ils sont également ajoutés à la scène.
   * 
   * @param players {Player[]} Les joueurs à ajouter.
   */
  public append(...players: Player[]): void {
    players.forEach(player => {
      this.players.push(player);
      this.container.addChild(player.graphics);
    }); 
  }

  /**
   * update Met à jour les joueurs gérés par le gestionnaire.
   * 
   * @param deltaTime {number} Le temps écoulé depuis la dernière mise à jour.
   */
  public update(deltaTime: number): void {
    this.players.forEach(player => player.update(deltaTime));
  }

  /**
   * draw Force le rendu des joueurs gérés par le gestionnaire.
   */
  public draw(): void {
    this.container.removeChildren();
    this.players.forEach(player => this.container.addChild(player.graphics));
  }
};

export { PlayerManager };