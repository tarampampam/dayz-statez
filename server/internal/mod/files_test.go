package mod_test

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"testing"

	"server/internal/mod"
	"server/internal/testutil/assert"
)

var (
	//go:embed testdata/world.json
	worldJSON []byte

	//go:embed testdata/players.json
	playersJSON []byte

	//go:embed testdata/vehicles.json
	vehiclesJSON []byte

	//go:embed testdata/areas.json
	areasJSON []byte

	//go:embed testdata/territories.json
	territoriesJSON []byte

	//go:embed testdata/zombies.json
	zombiesJSON []byte

	//go:embed testdata/animals.json
	animalsJSON []byte

	//go:embed testdata/weapons.json
	weaponsJSON []byte

	//go:embed testdata/hits.json
	hitsJSON []byte

	//go:embed testdata/events.json
	eventsJSON []byte
)

func TestDTO_Decode(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		give []byte
		want any
	}{
		"world": {worldJSON, mod.World{
			TS: 496053180, World: "chernarusplus", UptimeSeconds: 205.5, GameTimeUnix: 496053180, TimeOfDay: 1,
			MaxPlayers: 2, Infected: 48, InfectedDead: 3, Animals: 163, PathGraphUpdates: 80,
			Weapons: 901, Suppressors: 71, Optics: 468, AmmoBoxes: 365, Magazines: 763, Ammo: 1366, Explosives: 142,
			Items: 22712, Bottles: 311, Clothes: 8656, Boxes: 150, Cars: 43, Boats: 17,
			Tents: 80, Containers: 181, CarParts: 839, CarWheels: 472,
			EffectAreas: 9, CrashSites: 3,
			EAI: 4, EAINPC: 2, EAIDeaths: 1, EAINPCDeaths: 1,
			FPS:     mod.WorldFPS{Current: 6853.25, WindowMin: 5116.5, WindowMax: 7895.75, WindowAvg: 6687.5, WindowSamples: 15},
			Weather: mod.WorldWeather{Temperature: 9.48, TemperatureBase: 13.01, WindSpeed: 3.15, WindDirection: 190.41, Fog: 0.068, Clouds: 0.5723},
			Food: mod.WorldFood{
				Total: 1086, Other: 63, Fruit: 72, Disinfectant: 154, Pills: 175, Drink: 121, Snack: 141,
				CannedSmall: 75, CannedMedium: 65, CannedBig: 165, Jar: 55,
			},
		}},

		"players": {playersJSON, mod.Players{
			{
				SteamID: "76561198000000001", GUID: "nsvWKtFMDKKpPJ0=", Name: "tester",
				BloodType: "A+", Gender: "female", CharacterType: "Judy",
				IsLoaded: true, // from 1
				Health:   1.0, Blood: 0.89, Shock: 1.0, Energy: 0.5, Water: 0.6,
				TemperatureCelsius: 36.0, WeightGrams: 7000.5, Wetness: 0.1, LifetimeSeconds: 123.5,
				PosX: 13377.5, PosY: 6292.14, PosAlt: 7.17, PosAngle: -121.33,
				PingMin: 2, PingMax: 8,
				ImmunityBoosted:     false, // from 0
				IsUnconscious:       true,  // from true
				IsRestrained:        false, // from -1
				Godmode:             false, // from null
				StatPlaytimeSeconds: 5000.0, StatDistanceMeters: 693.69, StatInfectedKilled: 4,
			},
			{
				SteamID: "76561198000000002", GUID: "aBcDeFgHiJk=", Name: "boris",
				BloodType: "O-", Gender: "male", CharacterType: "Boris",
				IsLoaded: true, // from true
				Health:   0.75, Blood: 0.9, Shock: 0.95, Energy: 0.3, Water: 0.4, Toxicity: 0.2,
				TemperatureCelsius: 35.2, WeightGrams: 9500.0, Wetness: 0.6, LifetimeSeconds: 4200.0,
				PosX: 4211.0, PosY: 8765.5, PosAlt: 412.3, PosAngle: 12.0,
				PingMin: 20, PingMax: 55, Throttle: 0.1, AgentsActive: 2, BleedingSources: 1,
				ImmunityBoosted:     true,  // from 1
				IsUnconscious:       false, // from 0
				IsRestrained:        false, // from 0
				Godmode:             true,  // from 1
				StatPlaytimeSeconds: 120000.0, StatDistanceMeters: 48213.7, StatLongestHitMeters: 312.5,
				StatPlayersKilled: 3, StatInfectedKilled: 420,
			},
		}},

		"vehicles": {vehiclesJSON, mod.Vehicles{
			{
				ClassName: "Hatchback_02_Blue", VehicleType: "Car", Hash: 384863173,
				Health: 0.948, EngineOn: false /* from "" */, FuelFraction: 0.0377,
				PosX: 13376.52, PosY: 6292.08, PosAlt: 5.96, PosAngle: 86.73,
			},
			{
				ClassName: "Boat_01_Black", VehicleType: "Boat", Hash: 346351774,
				Health: 1.0, CrewOccupied: 2, SpeedMS: 12.5, EngineOn: true /* from 1 */, FuelFraction: 0.3056,
				PosX: 13346.0, PosY: 7073.0, PosAlt: 1.0, PosAngle: 132.0,
			},
			{
				ClassName: "Expansion_Mi8", VehicleType: "Helicopter", Hash: -1, // hash sentinel
				Health: 1.0, EngineOn: false, FuelFraction: -1.0, // fuel sentinel: no fuel system
				PosX: 5000.0, PosY: 9000.0, PosAlt: 120.0, PosAngle: 270.0,
			},
		}},

		"areas": {areasJSON, mod.Areas{
			{Type: "STATIC", ClassName: "ContaminatedArea_Static", PosX: 13917.0, PosY: 11178.0, Radius: 75.0},
			{Type: "DYNAMIC", ClassName: "ContaminatedArea_Dynamic", PosX: 2080.0, PosY: 3500.0, Radius: 100.0, Insiders: 2},
		}},

		"territories": {territoriesJSON, mod.Territories{
			{PosX: 7500.25, PosY: 5500.5, LifetimeFraction: 0.75, RefresherRadius: 100.0},
		}},

		"zombies": {zombiesJSON, mod.Zombies{
			ByType:  map[string]int{"nbc": 3, "military": 5, "default": 10, "default_male": 8, "default_female": 6, "msoldier": 39, "soldier_heavy": 2},
			ByState: map[string]int{"calm": 45, "disturbed": 4, "alerted": 3, "chase": 2, "fight": 1, "dead": 3},
		}},

		"animals": {animalsJSON, mod.Animals{
			ByType: map[string]int{"sheep": 32, "pig": 20, "wolf": 16, "goat": 18, "deer": 59, "cow": 16, "bear": 2},
		}},

		"weapons": {weaponsJSON, mod.Weapons{
			Shots:       map[string]int{"akm": 120, "mosin9130": 12},
			WorldCount:  map[string]int{"akm": 4, "mosin9130": 65, "glock19": 40},
			PlayerKills: map[string]int{"akm": 2, "player_fist": 1, "player_melee_item_heavy": 1},
			CreatureKills: map[string]int{
				"akm": 5, "zombie": 3, "animal": 2, "transport": 1, "grenade": 1, "trap": 1, "tool": 1,
				"player_fist": 2, "player_melee_item_heavy": 1, "player_melee_item_two_hands": 1, "player_melee_item_one_hand": 1,
				"eai_fist": 1, "eai_melee_item_heavy": 1, "eai_melee_item_two_hands": 1, "eai_melee_item_one_hand": 1,
				"none": 1, "object": 1, "unknown": 1,
			},
		}},

		"hits": {hitsJSON, mod.Hits{
			PlayerHit:   map[string]int{"762x39": 3, "556x45": 1},
			CreatureHit: map[string]int{"762x54": 2, "9x19": 4},
		}},

		"events": {eventsJSON, mod.Events{Events: map[string]int{
			"Startup": 1, "WorldCleanUP": 1, "SelectedUserChanged": 1, "ScriptLog": 1,
			"MPSessionStart": 1, "MPSessionEnd": 1, "MPSessionFail": 1, "MPSessionPlayerReady": 1,
			"MPConnectionLost": 1, "MPConnectionClose": 1, "Progress": 2,
			"NetworkManagerClient": 1, "NetworkManagerServer": 1, "DialogQueued": 1, "ChatMessage": 1, "ChatChannel": 1,
			"ClientConnected": 1, "ClientPrepare": 1, "ClientNew": 1, "ClientNewReady": 1, "ClientReSpawn": 1,
			"ClientReconnect": 1, "ClientReady": 1, "ClientDisconnected": 1, "ClientRemoved": 1,
			"ConnectivityStatsUpdated": 1, "ServerFpsStatsUpdated": 1,
			"LogoutCancel": 1, "LoginTime": 1, "Respawn": 1, "Preload": 1, "Logout": 1, "LoginStatus": 1,
			"VONState": 1, "VONStartSpeaking": 1, "VONStopSpeaking": 1, "VONStartedTransmitting": 1, "VONStoppedTransmitting": 1,
			"PartyChatStatusChanged": 1, "DLCOwnerShipFailed": 1, "SetFreeCamera": 1,
			"ConnectingStart": 1, "ConnectingAbort": 1, "PlayerDeath": 1, "NetworkInputBuffer": 1,
			"unknown_123": 1,
		}}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := reflect.New(reflect.TypeOf(tc.want)) // fresh *T of want's type

			assert.NoError(t, json.Unmarshal(tc.give, got.Interface()))
			assert.DeepEqual(t, tc.want, got.Elem().Interface())
		})
	}
}
