package handlers

import (
	"maps"

	"server/internal/mod"
)

type Animals struct {
	// ByType maps canonical animal type -> live count, e.g. "wolf", "deer", "cow", "pig", "sheep", "goat", "bear".
	// Keys are derived from the animal's steak/pelt config, lowercased. Dynamic key set.
	ByType map[string]int `json:"by_type"`
}

func (dst *Animals) FillFrom(src *mod.Animals) {
	if src == nil {
		dst.ByType = map[string]int{}

		return
	}

	dst.ByType = make(map[string]int, len(src.ByType))

	maps.Copy(dst.ByType, src.ByType)
}

// --------------------------------------------------------------------------------------------------------------------

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

func (dst *Areas) FillFrom(src *mod.Areas) {
	if src == nil {
		*dst = nil

		return
	}

	*dst = make(Areas, len(*src))

	for i := range *src {
		v, a := &(*src)[i], &(*dst)[i]

		a.Type = v.Type
		a.Pos.X = v.PosX
		a.Pos.Y = v.PosY
		a.Radius = v.Radius
		a.Insiders = v.Insiders
	}
}

// --------------------------------------------------------------------------------------------------------------------

type Events struct {
	// Events maps event name -> cumulative occurrences since server start. Keys are human-readable engine event names
	// (e.g. "ClientConnected", "ClientDisconnected", "PlayerDeath", "ChatMessage", "MPSessionStart"), or "unknown_<id>"
	// for event ids without a known name.
	Events map[string]int `json:"events"`
}

func (dst *Events) FillFrom(src *mod.Events) {
	if src == nil {
		dst.Events = map[string]int{}

		return
	}

	dst.Events = make(map[string]int, len(src.Events))

	maps.Copy(dst.Events, src.Events)
}

// --------------------------------------------------------------------------------------------------------------------

// Hits is the DTO for HitsFileName file.
type Hits struct {
	// PlayerHit maps ammo type -> cumulative qualifying hits on players.
	PlayerHit map[string]int `json:"player_hit"`
	// CreatureHit maps ammo type -> cumulative qualifying hits on creatures.
	CreatureHit map[string]int `json:"creature_hit"`
}

func (dst *Hits) FillFrom(src *mod.Hits) {
	if src == nil {
		dst.PlayerHit = map[string]int{}
		dst.CreatureHit = map[string]int{}

		return
	}

	dst.PlayerHit = make(map[string]int, len(src.PlayerHit))
	dst.CreatureHit = make(map[string]int, len(src.CreatureHit))

	maps.Copy(dst.PlayerHit, src.PlayerHit)
	maps.Copy(dst.CreatureHit, src.CreatureHit)
}

// --------------------------------------------------------------------------------------------------------------------

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

func (dst *Players) FillFrom(src *mod.Players) {
	if src == nil {
		*dst = nil

		return
	}

	*dst = make(Players, len(*src))

	for i := range *src {
		v, p := &(*src)[i], &(*dst)[i]

		p.ID.SteamID = v.SteamID
		p.ID.GUID = v.GUID
		p.ID.Name = v.Name

		p.Character.Type = v.CharacterType
		p.Character.Gender = v.Gender
		p.Character.BloodType = v.BloodType

		p.Pos.X = v.PosX
		p.Pos.Y = v.PosY
		p.Pos.Alt = v.PosAlt
		p.Pos.Angle = v.PosAngle

		p.Vitals.Health = v.Health
		p.Vitals.Blood = v.Blood
		p.Vitals.Shock = v.Shock
		p.Vitals.Energy = v.Energy
		p.Vitals.Water = v.Water
		p.Vitals.Toxicity = v.Toxicity
		p.Vitals.Wetness = v.Wetness
		p.Vitals.TemperatureCelsius = v.TemperatureCelsius

		p.State.IsLoaded = bool(v.IsLoaded)
		p.State.IsUnconscious = bool(v.IsUnconscious)
		p.State.IsRestrained = bool(v.IsRestrained)
		p.State.Godmode = bool(v.Godmode)
		p.State.ImmunityBoosted = bool(v.ImmunityBoosted)
		p.State.AgentsActive = v.AgentsActive
		p.State.BleedingSources = v.BleedingSources

		p.Network.Ping.Min = v.PingMin
		p.Network.Ping.Max = v.PingMax
		p.Network.Throttle = v.Throttle

		p.Stats.PlaytimeSeconds = v.StatPlaytimeSeconds
		p.Stats.DistanceMeters = v.StatDistanceMeters
		p.Stats.LongestHitMeters = v.StatLongestHitMeters
		p.Stats.PlayersKilled = v.StatPlayersKilled
		p.Stats.InfectedKilled = v.StatInfectedKilled

		p.WeightGrams = v.WeightGrams
		p.LifetimeSeconds = v.LifetimeSeconds
	}
}

// --------------------------------------------------------------------------------------------------------------------

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

func (dst *Territories) FillFrom(src *mod.Territories) {
	if src == nil {
		*dst = nil

		return
	}

	*dst = make(Territories, len(*src))

	for i := range *src {
		v, t := &(*src)[i], &(*dst)[i]

		t.Pos.X = v.PosX
		t.Pos.Y = v.PosY

		t.Refresher.LifetimeFraction = v.LifetimeFraction
		t.Refresher.Radius = v.RefresherRadius
	}
}

// --------------------------------------------------------------------------------------------------------------------

type Vehicles []struct {
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
		// Health is normalized global health, 0..1.
		Health float64 `json:"health"`
		// CrewOccupied is the number of occupied crew seats.
		CrewOccupied int `json:"crew_occupied"`
		// SpeedMS is the current speed in m/s (small values snapped to 0).
		SpeedMS float64 `json:"speed_ms"`
		// EngineOn reports whether the engine is running.
		EngineOn bool `json:"engine_on"`
		// FuelFraction is tank fill level, 0..1; -1.0 when the vehicle has no fuel system / fuel is not applicable.
		FuelFraction float64 `json:"fuel_fraction"`
	} `json:"state"`
}

func (dst *Vehicles) FillFrom(src *mod.Vehicles) {
	if src == nil {
		*dst = nil

		return
	}

	*dst = make(Vehicles, len(*src))

	for i := range *src {
		v, c := &(*src)[i], &(*dst)[i]

		c.Hash = v.Hash

		c.Class.Name = v.ClassName
		c.Class.Type = v.VehicleType

		c.Pos.X = v.PosX
		c.Pos.Y = v.PosY
		c.Pos.Alt = v.PosAlt
		c.Pos.Angle = v.PosAngle

		c.State.Health = v.Health
		c.State.CrewOccupied = v.CrewOccupied
		c.State.SpeedMS = v.SpeedMS
		c.State.EngineOn = bool(v.EngineOn)
		c.State.FuelFraction = v.FuelFraction
	}
}

// --------------------------------------------------------------------------------------------------------------------

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

func (dst *Weapons) FillFrom(src *mod.Weapons) {
	if src == nil {
		dst.Shots = map[string]int{}
		dst.WorldCount = map[string]int{}
		dst.Kills.Players = map[string]int{}
		dst.Kills.Creatures = map[string]int{}

		return
	}

	dst.Shots = make(map[string]int, len(src.Shots))
	dst.WorldCount = make(map[string]int, len(src.WorldCount))
	dst.Kills.Players = make(map[string]int, len(src.PlayerKills))
	dst.Kills.Creatures = make(map[string]int, len(src.CreatureKills))

	maps.Copy(dst.Shots, src.Shots)
	maps.Copy(dst.WorldCount, src.WorldCount)
	maps.Copy(dst.Kills.Players, src.PlayerKills)
	maps.Copy(dst.Kills.Creatures, src.CreatureKills)
}

// --------------------------------------------------------------------------------------------------------------------

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

func (dst *World) FillFrom(src *mod.World) {
	if src == nil {
		*dst = World{}

		return
	}

	dst.MapName = src.World
	dst.UptimeSeconds = int(src.UptimeSeconds)

	// in-game clock derived from the artificial in-game epoch; its intra-day part (% 86400) is exact
	sod := ((src.GameTimeUnix % 86400) + 86400) % 86400
	dst.Time.Hour = sod / 3600
	dst.Time.Minute = (sod % 3600) / 60

	switch src.TimeOfDay { // engine daytime phase (0..3) -> human-readable string
	case 0:
		dst.Time.TimeOfDay = "dawn"
	case 1:
		dst.Time.TimeOfDay = "day"
	case 2:
		dst.Time.TimeOfDay = "dusk"
	case 3:
		dst.Time.TimeOfDay = "night"
	}

	dst.Players.Online = src.PlayersOnline
	dst.Players.Max = src.MaxPlayers
	dst.Players.Spawns = src.PlayersSpawns
	dst.Players.Deaths = src.PlayersDeaths

	dst.Infected.Alive = src.Infected
	dst.Infected.Dead = src.InfectedDead
	dst.Infected.Deaths = src.InfectedDeaths

	dst.Animals.Alive = src.Animals
	dst.Animals.Dead = src.AnimalsDead
	dst.Animals.Deaths = src.AnimalsDeaths

	dst.Vehicles.Cars = src.Cars
	dst.Vehicles.Boats = src.Boats
	dst.Vehicles.Helicopters = src.Helicopters
	dst.Vehicles.Destroyed.Cars = src.CarsDestroyed
	dst.Vehicles.Destroyed.Boats = src.BoatsDestroyed
	dst.Vehicles.Destroyed.Helicopters = src.HelicoptersDestroyed

	dst.Weapons.Count = src.Weapons
	dst.Weapons.Suppressors = src.Suppressors
	dst.Weapons.Optics = src.Optics
	dst.Weapons.AmmoBoxes = src.AmmoBoxes
	dst.Weapons.Magazines = src.Magazines
	dst.Weapons.Ammo = src.Ammo
	dst.Weapons.Explosives = src.Explosives
	dst.Weapons.ExplosivesDetonated = src.ExplosivesDetonated

	dst.Entities.Items = src.Items
	dst.Entities.Bottles = src.Bottles
	dst.Entities.Clothes = src.Clothes
	dst.Entities.Boxes = src.Boxes
	dst.Entities.TerritoryFlags = src.TerritoryFlags
	dst.Entities.BaseBuildings = src.BaseBuildings
	dst.Entities.Tents = src.Tents
	dst.Entities.Stashes = src.Stashes
	dst.Entities.Gardens = src.Gardens
	dst.Entities.Containers = src.Containers
	dst.Entities.CarParts = src.CarParts
	dst.Entities.CarWheels = src.CarWheels

	dst.Activity.EffectAreas = src.EffectAreas
	dst.Activity.Corpses = src.Corpses
	dst.Activity.CrashSites = src.CrashSites
	dst.Activity.ArtilleryBarrages = src.ArtilleryBarrages
	dst.Activity.GiftsUnpacked = src.GiftsUnpacked
	dst.Activity.ChatMessages = src.ChatMessages
	dst.Activity.PathGraphUpdates = src.PathGraphUpdates

	dst.EAI.Survivors = src.EAI
	dst.EAI.NPCs = src.EAINPC
	dst.EAI.Deaths = src.EAIDeaths
	dst.EAI.NPCDeaths = src.EAINPCDeaths

	dst.FPS.Current = src.FPS.Current
	dst.FPS.WindowMin = src.FPS.WindowMin
	dst.FPS.WindowMax = src.FPS.WindowMax
	dst.FPS.WindowAvg = src.FPS.WindowAvg
	dst.FPS.WindowSamples = src.FPS.WindowSamples
	dst.FPS.Limit = src.FPS.Limit

	dst.Weather.Temperature = src.Weather.Temperature
	dst.Weather.TemperatureBase = src.Weather.TemperatureBase
	dst.Weather.WindSpeed = src.Weather.WindSpeed
	dst.Weather.WindDirection = src.Weather.WindDirection
	dst.Weather.Rain = src.Weather.Rain
	dst.Weather.Snow = src.Weather.Snow
	dst.Weather.Fog = src.Weather.Fog
	dst.Weather.Clouds = src.Weather.Clouds

	dst.Food.Total = src.Food.Total
	dst.Food.Other = src.Food.Other
	dst.Food.Fruit = src.Food.Fruit
	dst.Food.Mushroom = src.Food.Mushroom
	dst.Food.Corpse = src.Food.Corpse
	dst.Food.Worm = src.Food.Worm
	dst.Food.Guts = src.Food.Guts
	dst.Food.Meat = src.Food.Meat
	dst.Food.HumanMeat = src.Food.HumanMeat
	dst.Food.Disinfectant = src.Food.Disinfectant
	dst.Food.Pills = src.Food.Pills
	dst.Food.Drink = src.Food.Drink
	dst.Food.Snack = src.Food.Snack
	dst.Food.Candy = src.Food.Candy
	dst.Food.CannedSmall = src.Food.CannedSmall
	dst.Food.CannedMedium = src.Food.CannedMedium
	dst.Food.CannedBig = src.Food.CannedBig
	dst.Food.Jar = src.Food.Jar

	dst.Mining.Bushes = src.Mining.Bushes
	dst.Mining.Rocks = src.Mining.Rocks
	dst.Mining.Trees = src.Mining.Trees

	dst.Fishing.Attempts = src.Fishing.Attempts
	dst.Fishing.Catches = src.Fishing.Catches
}

// --------------------------------------------------------------------------------------------------------------------

type Zombies struct {
	// ByType maps canonical infected type -> live count. Keys are derived from the aiAgentTemplate (e.g. "msoldier",
	// "soldier_heavy", "default_male", "default_female", "military", "runner", or "nbc" for contamination-resistant
	// infected). Dynamic key set.
	ByType map[string]int `json:"by_type"`
	// ByState maps mind state -> live count. Keys: "calm", "disturbed", "alerted", "chase", "fight", and "dead"
	// (killed, pending cleanup).
	ByState map[string]int `json:"by_state"`
}

func (dst *Zombies) FillFrom(src *mod.Zombies) {
	if src == nil {
		dst.ByType = map[string]int{}
		dst.ByState = map[string]int{}

		return
	}

	dst.ByType = make(map[string]int, len(src.ByType))
	dst.ByState = make(map[string]int, len(src.ByState))

	maps.Copy(dst.ByType, src.ByType)
	maps.Copy(dst.ByState, src.ByState)
}
