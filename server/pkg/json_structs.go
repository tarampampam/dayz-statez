package pkg

// Animals is the live wildlife census - count of alive animals per canonical type. Use for population/spawn-balance
// monitoring and wildlife density dashboards.
type Animals struct {
	// ByType maps canonical animal type -> live count, e.g. "wolf", "deer", "cow", "pig", "sheep", "goat", "bear".
	// Keys are derived from the animal's steak/pelt config, lowercased. Dynamic key set.
	ByType map[string]int `json:"by_type"`
}

// --------------------------------------------------------------------------------------------------------------------

// Areas is the active effect areas - static preset gas zones and dynamic toxic/artillery events - each with its
// type, center position, radius and the number of players currently inside. Use for a live hazard map and
// player-in-zone alerts.
type Areas []struct {
	// Type is the zone type from engine enum eZoneType: "STATIC" (preset gas zones) or "DYNAMIC"
	// (artillery/toxic events).
	Type string `json:"type"`
	// Center on the horizontal map plane.
	Pos struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"pos"`
	// Radius is the area radius in meters.
	Radius float64 `json:"radius"`
	// Insiders is the number of players currently inside the area.
	Insiders int `json:"insiders"`
}

// --------------------------------------------------------------------------------------------------------------------

// Events is the cumulative counters (since server start) of engine events keyed by name - session lifecycle, client
// connect/disconnect, deaths, chat, VON, etc. Use for activity analytics and diagnosing connection/session issues.
type Events struct {
	// Events maps event name -> cumulative occurrences since server start. Keys are human-readable engine event names
	// (e.g. "ClientConnected", "ClientDisconnected", "PlayerDeath", "ChatMessage", "MPSessionStart"), or "unknown_<id>"
	// for event ids without a known name.
	Events map[string]int `json:"events"`
}

// --------------------------------------------------------------------------------------------------------------------

// Hits is the cumulative qualifying hit counts since server start, split into player_hit and creature_hit and broken
// down by ammo type. Use for combat/weapon-usage analytics.
type Hits struct {
	// PlayerHit maps ammo type -> cumulative qualifying hits on players.
	PlayerHit map[string]int `json:"player_hit"`
	// CreatureHit maps ammo type -> cumulative qualifying hits on creatures.
	CreatureHit map[string]int `json:"creature_hit"`
}

// --------------------------------------------------------------------------------------------------------------------

// Players is the Snapshot of every online player - identity (steam id, guid, name), character, position, vitals
// (health/blood/shock/energy/water/toxicity/temperature), state flags (unconscious, restrained, godmode, bleeding,
// etc.), network and lifetime stats (playtime, distance, kills). Use for a live player map.
type Players []struct {
	ID struct {
		// SteamID is the player's Steam ID (PlayerIdentity.GetPlainId()).
		SteamID string `json:"steam_id"`
		// GUID is the DayZ identity id / BE GUID basis (PlayerIdentity.GetId()).
		GUID string `json:"guid"`
		// Name is the player display name.
		Name string `json:"name"`
	} `json:"id"`

	// Character describes the survivor identity/appearance.
	Character struct {
		// Type is the survivor variant with the "SurvivorM_"/"SurvivorF_" prefix stripped by the mod,
		// e.g. "Judy", "Boris". Derived from GetType().
		Type string `json:"type"`
		// Gender is "male" or "female".
		Gender string `json:"gender"`
		// BloodType is the human-readable blood type incl. Rh sign, one of:
		// "A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-".
		BloodType string `json:"blood_type"`
	} `json:"character"`

	// Pos is the player position: x/y horizontal map plane (engine X / engine Z), alt altitude/height (engine Y, m),
	// angle yaw heading in degrees.
	Pos struct {
		X     float64 `json:"x"`
		Y     float64 `json:"y"`
		Alt   float64 `json:"alt"`
		Angle float64 `json:"angle"`
	} `json:"pos"`

	// Vitals are normalized life stats (0..1) plus body temperature.
	Vitals struct {
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
		// Wetness is the wetness stat value.
		Wetness float64 `json:"wetness"`
		// TemperatureCelsius is the player's body temperature, C (~36 when healthy).
		TemperatureCelsius float64 `json:"temperature_celsius"`
	} `json:"vitals"`

	// State are condition flags and active-condition counters.
	State struct {
		// IsLoaded reports whether the character finished loading (IsPlayerLoaded()).
		IsLoaded bool `json:"is_loaded"`
		// IsUnconscious reports whether the player is unconscious.
		IsUnconscious bool `json:"is_unconscious"`
		// IsRestrained reports whether the player is restrained (handcuffed).
		IsRestrained bool `json:"is_restrained"`
		// Godmode reports whether damage is disabled (!GetAllowDamage()).
		Godmode bool `json:"godmode"`
		// ImmunityBoosted reports whether immunity is currently boosted.
		ImmunityBoosted bool `json:"immunity_boosted"`
		// AgentsActive is the number of active agents/pathogens (set bits of GetAgents()).
		AgentsActive int `json:"agents_active"`
		// BleedingSources is the number of active bleeding sources (set bits of GetBleedingBits()).
		BleedingSources int `json:"bleeding_sources"`
	} `json:"state"`

	// Network quality observed over the collection window.
	Network struct {
		// Ping is the min/max ping (ms) observed over the window; 0 if unavailable.
		Ping struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"ping"`
		// Throttle is the network output throttle fraction observed in the window.
		Throttle float64 `json:"throttle"`
	} `json:"network"`

	// Stats are lifetime analytics counters (persist across sessions).
	Stats struct {
		// PlaytimeSeconds is lifetime playtime from server analytics, seconds.
		PlaytimeSeconds float64 `json:"playtime_sec"`
		// DistanceMeters is lifetime distance traveled, meters.
		DistanceMeters float64 `json:"distance_meters"`
		// LongestHitMeters is the longest survivor hit distance, meters.
		LongestHitMeters float64 `json:"longest_hit_meters"`
		// PlayersKilled is lifetime players killed.
		PlayersKilled int `json:"players_killed"`
		// InfectedKilled is lifetime infected killed.
		InfectedKilled int `json:"infected_killed"`
	} `json:"stats"`

	// WeightGrams is the carried inventory load in grams (GetPlayerLoad()).
	WeightGrams float64 `json:"weight_grams"`
	// LifetimeSeconds is seconds since this character's metrics were created (≈ seconds the character entity has
	// existed this session).
	LifetimeSeconds float64 `json:"lifetime_seconds"`
}

// --------------------------------------------------------------------------------------------------------------------

// Territories is the base territory flags - position, remaining refresh-timer fraction and refresh radius. Use to
// map bases and track which are kept alive vs decaying toward cleanup.
type Territories []struct {
	// Pos is the flag's horizontal map-plane position.
	Pos struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"pos"`

	// Refresher is the base-flag refresh state.
	Refresher struct {
		// LifetimeFraction is the remaining refresh-timer fraction, 0..1 (GetRefresherTime01()); 1.0 = freshly refreshed.
		LifetimeFraction float64 `json:"lifetime_fraction"`
		// Radius is the flag refresh radius in meters (GameConstants.REFRESHER_RADIUS).
		Radius float64 `json:"radius"`
	} `json:"refresher"`
}

// --------------------------------------------------------------------------------------------------------------------

type (
	// Vehicles contains every tracked vehicle - class/type, position, health, crew occupied/capacity, engine, lights,
	// fuel & radiator-coolant levels, cargo item count and per-part health. Use for a live vehicle map, fleet
	// condition monitoring and spotting free vs player-claimed cars (empty cargo + full hp + still at a spawn
	// point = unclaimed).
	Vehicles []struct {
		// Hash is a stable non-negative hash of the vehicle's persistent ID (0..2^31-1);
		// -1 when the persistent ID is unavailable.
		Hash int `json:"hash"`

		// Class is the vehicle classification.
		Class struct {
			// Name is the entity class, e.g. "Hatchback_02_Blue", "Boat_01_Black".
			Name string `json:"name"`
			// Type is the engine vehicle category (from GetVehicleType()) with the "VehicleType" prefix stripped by
			// the mod: "Car", "Boat", "Helicopter" (Expansion vehicles may report their own category).
			Type string `json:"type"`
		} `json:"class"`

		// Pos is the vehicle position: x/y horizontal map plane, alt height (m), angle yaw degrees.
		Pos struct {
			X     float64 `json:"x"`
			Y     float64 `json:"y"`
			Alt   float64 `json:"alt"`
			Angle float64 `json:"angle"`
		} `json:"pos"`

		// State is the vehicle's dynamic status.
		State struct {
			// Health is the normalized global health of the vehicle BODY only (chassis/panels/engine zone), 0..1.
			// Attached parts carry their own condition in Parts.
			Health float64 `json:"health"`
			// SpeedMS is the current speed in m/s (small values snapped to 0).
			SpeedMS float64 `json:"speed_ms"`

			// Engine is the powertrain status.
			Engine struct {
				// On reports whether the engine is running.
				On bool `json:"on"`
				// RPM is the current engine revolutions per minute; 0 when off, and always 0 for boats (no RPM model).
				RPM int `json:"rpm"`
				// Gear is the selected gear index (0 = neutral on most gearboxes); -1 for vehicles without gears (boats).
				Gear int `json:"gear"`
			} `json:"engine"`

			// LightsOn reports whether the vehicle lights are on.
			LightsOn bool `json:"lights_on"`

			// Crew counts the occupied vs total crew seats.
			Crew struct {
				Occupied int `json:"occupied"`
				Capacity int `json:"capacity"`
			} `json:"crew"`

			// FuelFraction is tank fill level, 0..1; -1.0 when the vehicle has no fuel system / fuel is not applicable.
			FuelFraction float64 `json:"fuel_fraction"`
			// CoolantFraction is the radiator coolant level, 0..1; -1.0 when the vehicle has no radiator
			// (motorbikes, boats, and other types without a cooling system).
			CoolantFraction float64 `json:"coolant_fraction"`
			// CargoItems is the total number of items stored in the vehicle's cargo (trunk + storage attachments),
			// excluding the attached parts themselves.
			CargoItems int `json:"cargo_items"`
		} `json:"state"`

		// Parts lists the vehicle's attached parts (wheels, doors, battery, radiator, spark plug, headlights, ...),
		// each with its own health 0..1. Empty for vehicles without detachable parts (boats).
		Parts []VehiclePart `json:"parts"`
	}

	// VehiclePart is a single attached vehicle part with its own condition.
	VehiclePart struct {
		// Name is the part entity class, e.g. "CarBattery", "CarRadiator", "SparkPlug", "Hatchback_02_Wheel".
		Name string `json:"name"`
		// Health is the part's normalized health, 0..1.
		Health float64 `json:"health"`
	}
)

// --------------------------------------------------------------------------------------------------------------------

// Weapons is the cumulative shots fired per weapon type, current live weapon-entity counts, and kills grouped by
// source (players vs creatures). Use for weapon-usage/balance analytics and loot-economy insight.
type Weapons struct {
	// Shots maps weapon type -> cumulative shots fired since server start.
	Shots map[string]int `json:"shots"`
	// WorldCount maps weapon type -> current live weapon-entity count.
	WorldCount map[string]int `json:"world_count"`

	// Kills breaks kills down by kill-source.
	Kills struct {
		// Players maps kill-source -> cumulative players killed by that source.
		Players map[string]int `json:"players"`
		// Creatures maps kill-source -> cumulative creatures (infected/animals) killed. Source keys are a weapon
		// type, or one of the synthetic tokens: "zombie", "animal", "transport", "grenade", "trap", "tool",
		// "player_fist", "player_melee_item_heavy|two_hands|one_hand", "none", "object", "unknown"
		// (and "eai_*" variants on EXPANSIONMODAI builds).
		Creatures map[string]int `json:"creatures"`
	} `json:"kills"`
}

// --------------------------------------------------------------------------------------------------------------------

// World is the single world snapshot - map name, uptime, in-game clock (hour, minute, time-of-day phase), population
// counts (players/infected/animals plus cumulative deaths), vehicle counts, loot & entity breakdowns, server FPS,
// weather, and food/mining/fishing stats.
// Use for a server status dashboard and performance/weather/economy monitoring.
type World struct {
	// MapName is the lowercased world/map name, e.g. "chernarusplus", "enoch".
	MapName string `json:"map_name"`
	// UptimeSeconds is server process uptime in seconds (engine tick time).
	UptimeSeconds int `json:"uptime_seconds"`

	// Time is the in-game clock.
	Time struct {
		// Hour is the in-game hour, 0..23.
		Hour int `json:"hour"`
		// Minute is the in-game minute, 0..59.
		Minute int `json:"minute"`
		// TimeOfDay is the daytime phase: "dawn", "day", "dusk", "night".
		TimeOfDay string `json:"time_of_day"`
	} `json:"time"`

	// Players counts live connections and cumulative spawn/death events.
	Players struct {
		Online int `json:"online"` // live
		Max    int `json:"max"`
		Spawns int `json:"spawns"` // cumulative
		Deaths int `json:"deaths"` // cumulative
	} `json:"players"`

	// Infected counts live zombies, corpses, and cumulative deaths.
	Infected struct {
		Alive  int `json:"alive"`  // live
		Dead   int `json:"dead"`   // corpses not yet cleaned
		Deaths int `json:"deaths"` // cumulative
	} `json:"infected"`

	// Animals counts live animals, corpses, and cumulative deaths.
	Animals struct {
		Alive  int `json:"alive"`  // live
		Dead   int `json:"dead"`   // corpses
		Deaths int `json:"deaths"` // cumulative
	} `json:"animals"`

	// Vehicles counts live vehicles plus cumulative destructions.
	Vehicles struct {
		Cars        int `json:"cars"`        // live
		Boats       int `json:"boats"`       // live
		Helicopters int `json:"helicopters"` // live
		Destroyed   struct {
			Cars        int `json:"cars"`        // cumulative
			Boats       int `json:"boats"`       // cumulative
			Helicopters int `json:"helicopters"` // cumulative
		} `json:"destroyed"`
	} `json:"vehicles"`

	// Weapons & ammo loot (live entity counts; explosives_detonated is cumulative).
	Weapons struct {
		Count               int `json:"count"` // live weapon entities
		Suppressors         int `json:"suppressors"`
		Optics              int `json:"optics"`
		AmmoBoxes           int `json:"ammo_boxes"`
		Magazines           int `json:"magazines"`
		Ammo                int `json:"ammo"`
		Explosives          int `json:"explosives"`
		ExplosivesDetonated int `json:"explosives_detonated"` // cumulative
	} `json:"weapons"`

	// Entities is a breakdown of live world items and structures.
	Entities struct {
		Items          int `json:"items"`           // all ItemBase entities
		Bottles        int `json:"bottles"`         // Bottle_Base entities
		Clothes        int `json:"clothes"`         // Clothing entities
		Boxes          int `json:"boxes"`           // non-ammo Box_Base entities
		TerritoryFlags int `json:"territory_flags"` // live territory (base) flags
		BaseBuildings  int `json:"base_buildings"`  // BaseBuildingBase entities
		Tents          int `json:"tents"`           // TentBase entities
		Stashes        int `json:"stashes"`         // UndergroundStash entities
		Gardens        int `json:"gardens"`         // GardenPlot entities
		Containers     int `json:"containers"`      // DeployableContainer_Base entities
		CarParts       int `json:"car_parts"`       // CarDoor entities
		CarWheels      int `json:"car_wheels"`      // CarWheel items
	} `json:"entities"`

	// Activity is world events - live spawns and cumulative counters.
	Activity struct {
		EffectAreas       int `json:"effect_areas"`       // live effect areas (static + dynamic contaminated zones)
		Corpses           int `json:"corpses"`            // live dead-player corpses tracked by the mission
		CrashSites        int `json:"crash_sites"`        // live crash sites (CrashBase)
		ArtilleryBarrages int `json:"artillery_barrages"` // cumulative dynamic contaminated-area spawns since start
		GiftsUnpacked     int `json:"gifts_unpacked"`     // cumulative unpacked gifts since start
		ChatMessages      int `json:"chat_messages"`      // cumulative chat messages since start
		PathGraphUpdates  int `json:"path_graph_updates"` // cumulative navmesh region updates since start
	} `json:"activity"`

	// EAI is the Expansion AI block; zero on non-EXPANSIONMODAI servers.
	EAI struct {
		Survivors int `json:"survivors"`  // live Expansion AI survivors
		NPCs      int `json:"npcs"`       // live Expansion AI NPCs
		Deaths    int `json:"deaths"`     // cumulative Expansion AI deaths
		NPCDeaths int `json:"npc_deaths"` // cumulative Expansion AI NPC deaths
	} `json:"eai"`

	// FPS is the server frame-rate block. "window" values span the last collection interval;
	// "current" is the latest ~1s average.
	FPS struct {
		Current       float64 `json:"current"`        // latest ~1s average FPS
		WindowMin     float64 `json:"window_min"`     // min 1s-average FPS in the interval
		WindowMax     float64 `json:"window_max"`     // max 1s-average FPS in the interval
		WindowAvg     float64 `json:"window_avg"`     // mean of 1s-average FPS samples
		WindowSamples int     `json:"window_samples"` // number of 1s samples in the interval
		Limit         int     `json:"limit"`          // configured FPS cap (CLI -limitFPS); 0 = unlimited
	} `json:"fps"`

	// Weather block. Temperatures are C.
	Weather struct {
		Temperature     float64 `json:"temperature"`      // temperature at map center incl. weather modifiers, C
		TemperatureBase float64 `json:"temperature_base"` // base environment temperature (no local/weather modifiers), C
		WindSpeed       float64 `json:"wind_speed"`       // wind speed, m/s
		WindDirection   float64 `json:"wind_direction"`   // wind direction, degrees 0..360
		Rain            float64 `json:"rain"`             // rain intensity, 0..1
		Snow            float64 `json:"snow"`             // snowfall intensity, 0..1
		Fog             float64 `json:"fog"`              // fog density, 0..1
		Clouds          float64 `json:"clouds"`           // overcast, 0..1
	} `json:"weather"`

	// Food is the food-item breakdown (live entity counts per category). Total counts all edibles except bottles;
	// the rest are subsets.
	Food struct {
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
	} `json:"food"`

	// Mining is the cumulative count of harvested resources since server start.
	Mining struct {
		Bushes int `json:"bushes"` // bushes mined
		Rocks  int `json:"rocks"`  // rocks mined
		Trees  int `json:"trees"`  // trees mined
	} `json:"mining"`

	// Fishing is cumulative fishing activity since server start.
	Fishing struct {
		Attempts int `json:"attempts"` // fishing actions started
		Catches  int `json:"catches"`  // successful catches
	} `json:"fishing"`
}

// --------------------------------------------------------------------------------------------------------------------

// Zombies is the live infected census - counts per canonical type and per AI mind-state (calm, disturbed, alerted,
// chase, fight, dead). Use for horde-density and AI-activity monitoring and correlating it with server FPS.
type Zombies struct {
	// ByType maps canonical infected type -> live count. Keys are derived from the aiAgentTemplate (e.g. "msoldier",
	// "soldier_heavy", "default_male", "default_female", "military", "runner", or "nbc" for contamination-resistant
	// infected). Dynamic key set.
	ByType map[string]int `json:"by_type"`
	// ByState maps mind state -> live count. Keys: "calm", "disturbed", "alerted", "chase", "fight", and "dead"
	// (killed, pending cleanup).
	ByState map[string]int `json:"by_state"`
}
