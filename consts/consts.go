package consts

const (
	// --- GAME CONSTANTS
	// ================================

	// Tickrate defines the number of ticks per second.
	Tickrate = 30

	// TicksPerRound defines the number of ticks per round.
	TicksPerRound = 5 * 60 * 3

	// TicksPointRushStage defines the number of ticks for the point rush stage.
	TicksPointRushStage = 4 * 60 * 3

	// --- MAP CONSTANTS
	// ================================

	// MapWidth defines CELLS in x and y.
	MapWidth = 10

	// CellWidth defines width of cell
	CellWidth = 10

	// NumSubsquare defines the number of subsquares within a cell.
	NumSubsquare = 9.0

	// SubsquareRatio defines the ratio of subsquare to the cell.
	SubsquareRatio (float64) = 1.0 / NumSubsquare

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

	// RespawnTime defines the time (in seconds) for a player to respawn after being eliminated.
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

	// NumCoins defines the total number of coins available in the game.
	NumCoins = 30

	// BigCoinSize defines the size of big coin.
	BigCoinSize = 4

	// BigCoinValue defines the value when player take a big coin.
	BigCoinValue (int32) = 500
)
