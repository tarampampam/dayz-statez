package mod

import (
	"bytes"
	"strconv"
)

type DataFileName string

const (
	AnimalsFileName     DataFileName = "animals.json"     // ~80 bytes on empty server
	AreasFileName       DataFileName = "areas.json"       // ~1.5 KiB on empty server
	EventsFileName      DataFileName = "events.json"      // ~170 bytes on empty server
	HitsFileName        DataFileName = "hits.json"        // ~30 bytes on empty server
	PlayersFileName     DataFileName = "players.json"     // ~1 KiB for 1 player
	TerritoriesFileName DataFileName = "territories.json" // 2 bytes on empty server
	VehiclesFileName    DataFileName = "vehicles.json"    // ~17 KiB on empty server
	WeaponsFileName     DataFileName = "weapons.json"     // ~600 bytes on empty server
	WorldFileName       DataFileName = "world.json"       // ~1.5 KiB on empty server
	ZombiesFileName     DataFileName = "zombies.json"     // ~180 bytes on empty server
)

// esBool is a JSON-flexible boolean. EnforceScript (the DayZ mod runtime) serializes bool as the number
// 0/1, but other producers may emit a real JSON bool or a quoted string.
type esBool bool

// UnmarshalJSON accepts a JSON bool (true/false), number (0/1), or string ("true"/"1"/"0"/...). JSON null -> false.
func (b *esBool) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)

	// empty shit
	if len(data) == 0 {
		*b = false
		return nil
	}

	// null
	if bytes.Equal(data, []byte("null")) {
		return nil // encoding/json convention: null is a no-op
	}

	// bool-like value
	if v, err := strconv.ParseBool(string(data)); err == nil {
		*b = esBool(v)
		return nil
	}

	// give a chance for float/numeric-like value
	f, err := strconv.ParseFloat(string(data), 64)
	if err == nil {
		*b = f >= 1
	}

	return nil
}

// MarshalJSON always emits a real JSON boolean (true/false).
func (b esBool) MarshalJSON() ([]byte, error) {
	if b {
		return []byte("true"), nil
	}

	return []byte("false"), nil
}

// Animals is the DTO for AnimalsFileName file.
type Animals struct {
	// ByType maps canonical animal type -> live count, e.g. "wolf", "deer", "cow", "pig", "sheep", "goat", "bear".
	// Keys are derived from the animal's steak/pelt config, lowercased. Dynamic key set.
	ByType map[string]int `json:"by_type"`
}

type (
	// Areas is the DTO for AreasFileName file.
	Areas []Area

	// Area is a single contaminated/effect area.
	Area struct {
		// Type is the zone type from engine enum eZoneType: "STATIC" (preset gas zones) or "DYNAMIC"
		// (artillery/toxic events).
		Type string `json:"type"`
		// ClassName is the area entity class, e.g. "ContaminatedArea_Static", "ContaminatedArea_Dynamic".
		ClassName string `json:"class_name"`
		// PosX/PosY are the area center on the horizontal map plane.
		PosX float64 `json:"pos_x"`
		PosY float64 `json:"pos_y"`
		// Radius is the area radius in meters.
		Radius float64 `json:"radius"`
		// Insiders is the number of players currently inside the area.
		Insiders int `json:"insiders"`
	}
)

// Events is the DTO for EventsFileName file.
// Counts every engine EventType seen on the server.
type Events struct {
	// Events maps event name -> cumulative occurrences since server start. Keys are human-readable engine event names
	// (e.g. "ClientConnected", "ClientDisconnected", "PlayerDeath", "ChatMessage", "MPSessionStart"), or "unknown_<id>"
	// for event ids without a known name.
	Events map[string]int `json:"events"`
}

// Hits is the DTO for HitsFileName file.
type Hits struct {
	// PlayerHit maps ammo type -> cumulative qualifying hits on players.
	PlayerHit map[string]int `json:"player_hit"`
	// CreatureHit maps ammo type -> cumulative qualifying hits on creatures.
	CreatureHit map[string]int `json:"creature_hit"`
}

type (
	// Players is the DTO for PlayersFileName file.
	Players []Player

	// Player is a single player snapshot.
	Player struct {
		// SteamID is the player's Steam ID (PlayerIdentity.GetPlainId()).
		SteamID string `json:"steam_id"`
		// GUID is the DayZ identity id / BE GUID basis (PlayerIdentity.GetId()).
		GUID string `json:"guid"`
		// Name is the player display name.
		Name string `json:"name"`
		// BloodType is the human-readable blood type incl. Rh sign, one of:
		// "A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-".
		BloodType string `json:"blood_type"`
		// Gender is "male" or "female".
		Gender string `json:"gender"`
		// CharacterType is the survivor variant with the "SurvivorM_"/"SurvivorF_" prefix stripped by the mod,
		// e.g. "Judy", "Boris". Derived from GetType().
		CharacterType string `json:"character_type"`

		// IsLoaded reports whether the character finished loading (IsPlayerLoaded()).
		IsLoaded esBool `json:"is_loaded"`
		// Health is normalized global health, 0..1.
		Health float64 `json:"health"`
		// Blood is normalized blood level, 0..1.
		Blood float64 `json:"blood"`
		// Shock is normalized shock/consciousness health, 0..1.
		Shock float64 `json:"shock"`
		// Energy is normalized energy stat, 0..1.
		Energy float64 `json:"energy"`
		// Water is normalized hydration stat, 0..1.
		Water float64 `json:"water"`
		// Toxicity is normalized toxicity stat, 0..1.
		Toxicity float64 `json:"toxicity"`
		// TemperatureCelsius is the player's body temperature, C (~36 when healthy).
		TemperatureCelsius float64 `json:"temperature_celsius"`
		// WeightGrams is the carried load in grams (GetPlayerLoad()).
		WeightGrams float64 `json:"weight_grams"`
		// Wetness is the wetness stat value.
		Wetness float64 `json:"wetness"`
		// LifetimeSeconds is seconds since this character's metrics were created (≈ seconds the character entity has
		// existed this session).
		LifetimeSeconds float64 `json:"lifetime_seconds"`

		// PosX/PosY are the horizontal map plane (engine X / engine Z).
		PosX float64 `json:"pos_x"`
		PosY float64 `json:"pos_y"`
		// PosAlt is altitude/height (engine Y), meters.
		PosAlt float64 `json:"pos_alt"`
		// PosAngle is yaw heading in degrees.
		PosAngle float64 `json:"pos_angle"`

		// PingMin/PingMax are min/max ping (ms) observed over the collection window; 0 if unavailable.
		PingMin int `json:"ping_min"`
		PingMax int `json:"ping_max"`
		// Throttle is the network output throttle fraction observed in the window.
		Throttle float64 `json:"throttle"`

		// AgentsActive is the number of active agents/pathogens (set bits of GetAgents()).
		AgentsActive int `json:"agents_active"`
		// BleedingSources is the number of active bleeding sources (set bits of GetBleedingBits()).
		BleedingSources int `json:"bleeding_sources"`
		// ImmunityBoosted reports whether immunity is currently boosted.
		ImmunityBoosted esBool `json:"immunity_boosted"`
		// IsUnconscious reports whether the player is unconscious.
		IsUnconscious esBool `json:"is_unconscious"`
		// IsRestrained reports whether the player is restrained (handcuffed).
		IsRestrained esBool `json:"is_restrained"`
		// Godmode reports whether damage is disabled (!GetAllowDamage()).
		Godmode esBool `json:"godmode"`

		// StatPlaytimeSeconds is lifetime playtime from server analytics, seconds.
		StatPlaytimeSeconds float64 `json:"stat_playtime_seconds"`
		// StatDistanceMeters is lifetime distance traveled, meters.
		StatDistanceMeters float64 `json:"stat_distance_meters"`
		// StatLongestHitMeters is the longest survivor hit distance, meters.
		StatLongestHitMeters float64 `json:"stat_longest_hit_meters"`
		// StatPlayersKilled is lifetime players killed (analytics).
		StatPlayersKilled int `json:"stat_players_killed"`
		// StatInfectedKilled is lifetime infected killed (analytics).
		StatInfectedKilled int `json:"stat_infected_killed"`
	}
)

type (
	// Territories is the DTO for TerritoriesFileName file.
	Territories []Territory

	// Territory is a single base/territory flag.
	Territory struct {
		// PosX/PosY are the flag's horizontal map-plane position.
		PosX float64 `json:"pos_x"`
		PosY float64 `json:"pos_y"`
		// LifetimeFraction is the remaining refresh-timer fraction, 0..1 (GetRefresherTime01()); 1.0 = freshly refreshed.
		LifetimeFraction float64 `json:"lifetime_fraction"`
		// RefresherRadius is the flag refresh radius in meters (GameConstants.REFRESHER_RADIUS).
		RefresherRadius float64 `json:"refresher_radius"`
	}
)

type (
	// Vehicles is the DTO for VehiclesFileName file.
	Vehicles []Vehicle

	// Vehicle is a single vehicle snapshot.
	Vehicle struct {
		// ClassName is the entity class, e.g. "Hatchback_02_Blue", "Boat_01_Black".
		ClassName string `json:"class_name"`
		// VehicleType is the engine vehicle category (from GetVehicleType()) with the "VehicleType" prefix stripped by
		// the mod: "Car", "Boat", "Helicopter" (Expansion vehicles may report their own category).
		VehicleType string `json:"vehicle_type"`
		// Hash is a stable non-negative hash of the vehicle's persistent ID (0..2^31-1);
		// -1 when the persistent ID is unavailable.
		Hash int `json:"hash"`

		// Health is normalized global health, 0..1.
		Health float64 `json:"health"`
		// CrewOccupied is the number of occupied crew seats.
		CrewOccupied int `json:"crew_occupied"`
		// SpeedMS is the current speed in m/s (small values snapped to 0).
		SpeedMS float64 `json:"speed_ms"`
		// EngineOn reports whether the engine is running.
		EngineOn esBool `json:"engine_on"`
		// FuelFraction is tank fill level, 0..1; -1.0 when the vehicle has no fuel system / fuel is not applicable.
		FuelFraction float64 `json:"fuel_fraction"`

		// PosX/PosY horizontal map plane; PosAlt height (m); PosAngle yaw degrees.
		PosX     float64 `json:"pos_x"`
		PosY     float64 `json:"pos_y"`
		PosAlt   float64 `json:"pos_alt"`
		PosAngle float64 `json:"pos_angle"`
	}
)

// Weapons is the DTO for WeaponsFileName file.
type Weapons struct {
	// Shots maps weapon type -> cumulative shots fired since server start.
	Shots map[string]int `json:"shots"`
	// WorldCount maps weapon type -> current live weapon-entity count.
	WorldCount map[string]int `json:"world_count"`
	// PlayerKills maps kill-source -> cumulative players killed by that source.
	// See CreatureKills for the source vocabulary.
	PlayerKills map[string]int `json:"player_kills"`
	// CreatureKills maps kill-source -> cumulative creatures (infected/animals) killed. Source keys are a weapon type,
	// or one of the synthetic tokens: "zombie", "animal", "transport", "grenade", "trap", "tool", "player_fist",
	// "player_melee_item_heavy|two_hands|one_hand", "none", "object", "unknown"
	// (and "eai_*" variants on EXPANSIONMODAI builds).
	CreatureKills map[string]int `json:"creature_kills"`
}

type (
	// World is the DTO for WorldFileName file. It contains a single snapshot of the world state at the collection time.
	World struct {
		// TS is the snapshot timestamp in Unix-style seconds derived from the
		// in-game world date/time (minute precision), NOT the real wall clock.
		TS int `json:"ts"`
		// World is the lowercased world/map name, e.g. "chernarusplus", "enoch".
		World string `json:"world"`

		// UptimeSeconds is server process uptime in seconds (engine tick time).
		UptimeSeconds float64 `json:"uptime_seconds"`
		// GameTimeUnix is the in-game world time as Unix-style seconds (same basis
		// as TS, sampled at collection time).
		GameTimeUnix int `json:"game_time_unix"`
		// TimeOfDay is the daytime phase: 0=dawn, 1=day, 2=dusk, 3=night.
		TimeOfDay int `json:"time_of_day"`

		// PlayersOnline is the current connected player count.
		PlayersOnline int `json:"players_online"`
		// MaxPlayers is the server slot limit (serverDZ.cfg maxPlayers).
		MaxPlayers int `json:"max_players"`
		// Infected is the current live zombie (infected) count.
		Infected int `json:"infected"`
		// InfectedDead is the current infected-corpse count (killed, not yet cleaned).
		InfectedDead int `json:"infected_dead"`
		// Animals is the current live animal count.
		Animals int `json:"animals"`
		// AnimalsDead is the current animal-corpse count.
		AnimalsDead int `json:"animals_dead"`
		// PathGraphUpdates is the cumulative count of navmesh region updates caused
		// by placing objects (UpdatePathgraphRegionByObject) since server start.
		PathGraphUpdates int `json:"path_graph_updates"`

		Weapons             int `json:"weapons"`              // live weapon entities
		Suppressors         int `json:"suppressors"`          // live suppressor items
		Optics              int `json:"optics"`               // live optic items
		AmmoBoxes           int `json:"ammo_boxes"`           // live ammo boxes (Box_Base named *AmmoBox*)
		Magazines           int `json:"magazines"`            // live magazines
		Ammo                int `json:"ammo"`                 // live ammo piles
		Explosives          int `json:"explosives"`           // live explosive items
		ExplosivesDetonated int `json:"explosives_detonated"` // cumulative detonations since start

		Items          int `json:"items"`           // all ItemBase entities
		Bottles        int `json:"bottles"`         // Bottle_Base entities
		Clothes        int `json:"clothes"`         // Clothing entities
		Boxes          int `json:"boxes"`           // non-ammo Box_Base entities
		Cars           int `json:"cars"`            // live cars
		Boats          int `json:"boats"`           // live boats
		Helicopters    int `json:"helicopters"`     // live helicopters
		TerritoryFlags int `json:"territory_flags"` // live territory (base) flags
		BaseBuildings  int `json:"base_buildings"`  // BaseBuildingBase entities
		Tents          int `json:"tents"`           // TentBase entities
		Stashes        int `json:"stashes"`         // UndergroundStash entities
		Gardens        int `json:"gardens"`         // GardenPlot entities
		Containers     int `json:"containers"`      // DeployableContainer_Base entities
		CarParts       int `json:"car_parts"`       // CarDoor entities
		CarWheels      int `json:"car_wheels"`      // CarWheel items

		EffectAreas       int `json:"effect_areas"`       // live effect areas (static + dynamic contaminated zones)
		Corpses           int `json:"corpses"`            // dead-player corpses tracked by the mission
		ArtilleryBarrages int `json:"artillery_barrages"` // cumulative dynamic contaminated-area spawns since start
		CrashSites        int `json:"crash_sites"`        // live crash sites (CrashBase)
		GiftsUnpacked     int `json:"gifts_unpacked"`     // cumulative unpacked gifts since start

		ChatMessages         int `json:"chat_messages"`
		PlayersSpawns        int `json:"players_spawns"`
		PlayersDeaths        int `json:"players_deaths"`
		InfectedDeaths       int `json:"infected_deaths"`
		AnimalsDeaths        int `json:"animals_deaths"`
		CarsDestroyed        int `json:"cars_destroyed"`
		BoatsDestroyed       int `json:"boats_destroyed"`
		HelicoptersDestroyed int `json:"helicopters_destroyed"`

		EAI          int `json:"eai,omitempty"`            // live Expansion AI survivors
		EAINPC       int `json:"eai_npc,omitempty"`        // live Expansion AI NPCs
		EAIDeaths    int `json:"eai_deaths,omitempty"`     // cumulative Expansion AI deaths
		EAINPCDeaths int `json:"eai_npc_deaths,omitempty"` // cumulative Expansion AI NPC deaths

		FPS     WorldFPS     `json:"fps"`
		Weather WorldWeather `json:"weather"`
		Food    WorldFood    `json:"food"`
		Mining  WorldMining  `json:"mining"`
		Fishing WorldFishing `json:"fishing"`
	}

	// WorldFPS is the server frame-rate block of WorldFileName. "window" values span the last collection interval;
	// "current" is the latest ~1s average.
	WorldFPS struct {
		Current       float64 `json:"current"`        // latest ~1s average FPS
		WindowMin     float64 `json:"window_min"`     // min 1s-average FPS in the interval
		WindowMax     float64 `json:"window_max"`     // max 1s-average FPS in the interval
		WindowAvg     float64 `json:"window_avg"`     // mean of 1s-average FPS samples
		WindowSamples int     `json:"window_samples"` // number of 1s samples in the interval
		Limit         int     `json:"limit"`          // configured FPS cap (CLI -limitFPS); 0 = unlimited
	}

	// WorldWeather is the weather block of WorldFileName. Temperatures are C.
	WorldWeather struct {
		Temperature     float64 `json:"temperature"`      // temperature at map center incl. weather modifiers, C
		TemperatureBase float64 `json:"temperature_base"` // base environment temperature (no local/weather modifiers), C
		WindSpeed       float64 `json:"wind_speed"`       // wind speed, m/s
		WindDirection   float64 `json:"wind_direction"`   // wind direction, degrees 0..360
		Rain            float64 `json:"rain"`             // rain intensity, 0..1
		Snow            float64 `json:"snow"`             // snowfall intensity, 0..1
		Fog             float64 `json:"fog"`              // fog density, 0..1
		Clouds          float64 `json:"clouds"`           // overcast, 0..1
	}

	// WorldFood is the food-item breakdown of WorldFileName (live entity counts per category). Total counts all edibles
	// except bottles; the rest are subsets.
	WorldFood struct {
		Total        int `json:"total"`
		Other        int `json:"other"`
		Fruit        int `json:"fruit"`
		Mushroom     int `json:"mushroom"`
		Corpse       int `json:"corpse"`
		Worm         int `json:"worm"`
		Guts         int `json:"guts"`
		Meat         int `json:"meat"`
		HumanMeat    int `json:"human_meat"`
		Disinfectant int `json:"disinfectant"`
		Pills        int `json:"pills"`
		Drink        int `json:"drink"`
		Snack        int `json:"snack"`
		Candy        int `json:"candy"`
		CannedSmall  int `json:"canned_small"`
		CannedMedium int `json:"canned_medium"`
		CannedBig    int `json:"canned_big"`
		Jar          int `json:"jar"`
	}

	// WorldMining is the cumulative count of harvested resources since server start.
	WorldMining struct {
		Bushes int `json:"bushes"` // bushes mined
		Rocks  int `json:"rocks"`  // rocks mined
		Trees  int `json:"trees"`  // trees mined
	}

	// WorldFishing is cumulative fishing activity since server start.
	WorldFishing struct {
		Attempts int `json:"attempts"` // fishing actions started
		Catches  int `json:"catches"`  // successful catches
	}
)

// Zombies is the DTO for ZombiesFileName file.
type Zombies struct {
	// ByType maps canonical infected type -> live count. Keys are derived from the aiAgentTemplate (e.g. "msoldier",
	// "soldier_heavy", "default_male", "default_female", "military", "runner", or "nbc" for contamination-resistant
	// infected). Dynamic key set.
	ByType map[string]int `json:"by_type"`
	// ByState maps mind state -> live count. Keys: "calm", "disturbed", "alerted", "chase", "fight", and "dead"
	// (killed, pending cleanup).
	ByState map[string]int `json:"by_state"`
}
