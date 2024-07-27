import { MoveTo, ShootAt, Store, SwitchWeapon, BladeRotate, Weapon } from '../core/action.js';

class MyBot {
    constructor() {
        this.name = 'name_of_my_super_cool_bot';
        this.map = [[]];
    }

    /**
     * (fr)
     *      Cette méthode est appelé lorsque votre bot se connecte à la partie
     *      ou lorsqu'une nouvelle partie commence. Cela vous donne les 
     *      informations nécessaires sur la map.
     * @param {Model.MapState} state 
     * @returns {}
     */
    on_start(state) {
        this.map = state.map;
    }

    /**
     * (fr)
     *      Cette méthode est appelée à chaque tick de jeu. Vous pouvez y définir 
     *      le comportement de voter bot. Elle doit retourner une liste d'actions 
     *      qui sera exécutée par le serveur.
     *  
     *      Liste des actions possibles:
     *      - MoveTo({x, y})    permet de diriger son bot, il ira a vitesse
     *                          constante jusqu'à ce point.
     * 
     *      - ShootAt({x, y})   Si vous avez le fusil comme arme, cela va tirer
     *                          à la coordonnée donnée.
     * 
     *      - Store([...])      Permet de storer 100 octets dans le serveur. Lors
     *                          de votre reconnection, ces données vous seront
     *                          redonnées par le serveur.
     * 
     *      - SwtichWeapon(id)  Permet de changer d'arme. Par défaut, votre bot
     *                          n'est pas armé, voici vos choix:
     *                              Weapon.None
     *                              Weapon.Gun
     *                              Weapon.Blade
     *
     * (en)
     *      TODO
     * 
     * @param {Model.GameState} message 
     * @returns{Model.Actions}
     */
    on_tick(message) {
        console.log(`Current tick: ${message.tick}`);

        return [
            new SwitchWeapon(Weapon.Blade)
        ];
    }

    on_end() {

    }
};

export { MyBot };