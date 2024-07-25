package consts

const (
	// --- GAME CONSTANTS
	// ================================

	// Tickrate defines the number of ticks per second.
	Tickrate = 30

	// TicksPerRound defines the number of ticks per round.
	TicksPerRound = 5 * 60 * 3 * 10

	// TicksPointRushStage defines the number of ticks for the point rush stage.
	TicksPointRushStage = 4 * 60 * 3 * 10

	// --- MAP CONSTANTS
	// ================================

	// MapWidth defines the width of the map in cells.
	MapWidth = 10

	// CellWidth defines the width of each cell.
	CellWidth = 10

	// NumSubsquare defines the number of subsquares within a cell.
	NumSubsquare = 9.0

	// SubsquareRatio defines the ratio of a subsquare to the cell.
	SubsquareRatio float64 = 1.0 / NumSubsquare

	// SubsquareWidth defines the width of a subsquare.
	SubsquareWidth float64 = 10.0 * SubsquareRatio

	// --- PLAYER CONSTANTS
	// ================================

	// PlayerSize defines the size of a player.
	PlayerSize = 1

	// PlayerHealth is the starting health of a player.
	PlayerHealth = 100

	// PlayerSpeed defines the distance traveled per tick.
	PlayerSpeed = 1.15

	// RespawnTime defines the time (in seconds) for a player to respawn after being eliminated.
	RespawnTime = 5

	// --- PROJECTILE CONSTANTS
	// ================================

	// ProjectileSize defines the size of a projectile.
	ProjectileSize = 0.35

	// ProjectileDmg defines the damage suffered by a player when hit by a projectile.
	ProjectileDmg = 15

	// ProjectileSpeed defines the speed of a projectile.
	ProjectileSpeed = 3.0

	// ProjectileTTL defines the time to live of a projectile (in seconds).
	ProjectileTTL = 5.0

	// --- BLADE CONSTANTS
	// ================================

	// BladeSize defines the size of a blade.
	BladeSize = PlayerSize/2.0 + 1.5

	// BladeDmg defines the damage suffered by a player when hit by a blade.
	BladeDmg = 4

	// BladeRotationSpeed defines the speed of rotation of a blade (in degrees).
	BladeRotationSpeed = 230

	// --- SCORER CONSTANTS
	// ================================

	// CoinSize defines the size of a coin.
	CoinSize = 0.5

	// CoinValue defines the value when a player collects a coin.
	CoinValue int32 = 40

	// NumCoins defines the total number of coins available in the game.
	NumCoins = 30

	// BigCoinSize defines the size of a big coin.
	BigCoinSize = 4

	// BigCoinValue defines the value when a player collects a big coin.
	BigCoinValue int32 = NumCoins * CoinValue

	// --- SCORE CONSTANTS
	// ================================

	// ScoreOnHitWithProjectile defines the score awarded when hitting an opponent with a projectile.
	ScoreOnHitWithProjectile = 15

	// ScoreOnHitWithBlade defines the score awarded when hitting an opponent with a blade.
	ScoreOnHitWithBlade = 10
)
