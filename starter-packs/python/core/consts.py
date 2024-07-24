class Consts:
    class Map:
        """
        (fr) Constantes relatives à la cate du jeu.
        (en) Constants related to the game map.

        Attributes:
            WIDTH       (int) : (fr) La largeur de la carte en nombre de cellules. 
                                (en) The width of the map in number of cells.

            HEIGHT      (int) : (fr) La hauteur de la carte en nombre de cellules.
                                (en) The height of the map in number of cells.

            CELL_WIDTH  (int) : (fr) La largeur de chaque cellule sur la carte.
                                (en) The width of each cell in the map.

            CELL_HEIGHT (int) : (fr) La hauteur de chaque cellule sur la carte.
                                (en) The height of each cell in the map.
        """
        WIDTH       = 10
        HEIGHT      = 10
        CELL_WIDTH  = 10
        CELL_HEIGHT = 10


    class Player:
        """
        (fr) Constantes relatives au joueur.
        (en) Constants related to the player.

        Attributes:
            SIZE         (float) : (fr) La taille du joeur.
                                   (en) The size of the player.

            SPEED        (float) : (fr) La vitesse du joueur (par seconde).
                                   (en) The speed of the player (per second).

            MAX_HEALTH   (int)   : (fr) La vie maximale du joueur.
                                   (en) The maximum health of the player.

            RESPAWN_TIME (float) : (fr) Le temps nécessaire pour que le joueur réapparaisse après avoir été éliminé.
                                   (en) The time it takes for the player to respawn after being eliminated.
        """
        SIZE         = 1.0
        SPEED        = 1.15
        MAX_HEALTH   = 100
        RESPAWN_TIME = 5.0


    class Projectile:
        """
        (fr) Constantes relatives aux projectiles.
        (en) Constants related to projectiles.

        Attributes:
            SIZE    (float) : (fr) La taille du projectile (rayon = SIZE / 2).
                              (en) The size of the projectile (radius = SIZE / 2).

            SPEED   (float) : (fr) La vitesse du projectile (par seconde).
                              (en) The speed of the projectile (per second).
                            
            DAMAGE  (int)   : (fr) Les dégats infligés par le projectile.
                              (en) The damage dealt by the projectile.

            TTL     (float) : (fr) La durée de vie du projectile (en secondes).
                              (en) The time-to-live of the projectile (in seconds).
        """
        SIZE    = 0.35
        SPEED   = 3.0
        DAMAGE  = 30
        TTL     = 5.0


    class Blade:
        """
        (fr) Constantes relatives aux lames.
        (en) Constants related to blades.

        Attributes:
            LENGTH    (float) : (fr) La longueur de la lame.
                                (en) The length of the blade.

            THICKNESS (float) : (fr) L'épaisseur de la lame.
                                (en) The thickness of the blade.
                            
            DAMAGE    (int)   : (fr) Les dégats infligés par la lame.
                                (en) The damage dealt by the blade.
        """
        LENGTH    = 2.0
        THICKNESS = 0.25
        DAMAGE    = 20


    class Coin:
        """
        (fr) Constantes relatives aux pièces.
        (en) Constants related to coins.

        Attributes:
            SIZE     (float) : (fr) La taille d'une pièce (rayon = SIZE / 2).
                               (en) The size of the coin (radius = SIZE / 2).

            VALUE    (int)   : (fr) La valeur d'une pièce lorsque ramassée.
                               (en) The value of coin when collected. 

            QUANTITY (int)   : (fr) La quantité de pièces dans la map.
                               (en) The quantity of coins in the map.
        """
        SIZE     = 0.5
        VALUE    = 40
        QUANTITY = 30


    class Treasure:
        """
        (fr) Constantes relative au trésor.
        (en) Constants related to trasure.

        Attributes:
            SIZE     (float) : (fr) La taille d'un trésor (rayon = SIZE / 2).
                               (en) The size of the treasure (radius = SIZE / 2).

            VALUE    (int)   : (fr) La valeur du trésor lorsque ramassée.
                               (en) The value of treasure when collected.   
        """
        SIZE  = 4.0
        VALUE = 1_200




