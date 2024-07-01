package config


const (
  TicksPerRound = 5 * 60 * 3
  TicksPointRushStage = 4 * 60 * 3

  // --- MAP CONSTANTS
  // ================================
  
  // MapWidth defines CELLS in x and y.
  MapWidth = 10

  // CellWidth defines width of cell
  CellWidth = 10

  SubsquareRatio (float64) = 1.0/9.0

  // SubsquareWidth defines width of sub square
  SubsquareWidth (float64) = 10.0 * SubsquareRatio


  // --- PLAYER CONSTANTS
  // ================================

  // PlayerSize defines the size of a player.
  PlayerSize = 1

  // PlayerHealth is the starting health of a player.
  PlayerHealth = 100

  // PlayerSpeed defines the distance travelled per tick.
  PlayerSpeed = 1.15

  RespawnTime = 5 


  // --- PORJECTILE CONSTANTS
  // ================================

  // ProjectileSize defines the sizes of a projectile. 
  ProjectileSize = 0.35

  // ProjectileDmg defines the damage suffered by a player when hit by a projectile
  ProjectileDmg = 30

  // ProjectileSpeed defines
  ProjectileSpeed = 1.75


  // --- BLADE CONSTANTS
  // ================================

  // BladeSize defines the size of a blade.
  BladeSize = 1.5

  // BladeDmg defines the damage suffered by a player when hit by a blade.
  BladeDmg = 20

  // BladeRotationSpeed defines the speed of rotation of a blade (in deg).
  BladeRotationSpeed = 230

  // BladeDistance defines the distance between the player and the blade.
  BladeDistance = 1.5

  // --- SCORER CONSTANTS
  // ================================
  
  // CoinSize defines the size of coin.
  CoinSize = 0.5

  // CoinValue defines the value when player take a coin.
  CoinValue (int32) = 50

  NumCoins = 30
)
