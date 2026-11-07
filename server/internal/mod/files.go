package mod

type FileName string

const (
	AnimalsFileName     FileName = "animals.json"     // ~80 bytes on empty server
	AreasFileName       FileName = "areas.json"       // ~1.5 KiB on empty server
	EventsFileName      FileName = "events.json"      // ~170 bytes on empty server
	HitsFileName        FileName = "hits.json"        // ~30 bytes on empty server
	PlayersFileName     FileName = "players.json"     // ~1 KiB for 1 player
	TerritoriesFileName FileName = "territories.json" // 2 bytes on empty server
	VehiclesFileName    FileName = "vehicles.json"    // ~17 KiB on empty server
	WeaponsFileName     FileName = "weapons.json"     // ~600 bytes on empty server
	WorldFileName       FileName = "world.json"       // ~1.5 KiB on empty server
	ZombiesFileName     FileName = "zombies.json"     // ~180 bytes on empty server
)
