package config


const (
  // --- PLAYER CONSTANTS
  // ================================

  // PlayerSize defines the size of a player.
  PlayerSize = 1

  // PlayerHealth is the starting health of a player.
  PlayerHealth = 100

  // PlayerSpeed defines the distance travelled per tick.
  PlayerSpeed = 1.15


  // --- PORJECTILE CONSTANTS
  // ================================

  // ProjectileSize defines the sizes of a projectile. 
  ProjectileSize = 0.35

  // ProjectileDmg defines the damage suffered by a player when hit by a projectile
  ProjectileDmg = 30

  // ProjectileSpeed defines
  ProjectileSpeed = 1.75


  // --- SCORER CONSTANTS
  // ================================
  
  // CoinSize defines the size of coin.
  CoinSize = 0.5

  // CoinValue defines the value when player take a coin.
  CoinValue (int) = 50
)
