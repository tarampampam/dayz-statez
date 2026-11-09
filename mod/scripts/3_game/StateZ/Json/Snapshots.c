/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER

// ---------------------------------------------------------------------------
// World snapshot sub-objects
// ---------------------------------------------------------------------------

class StateZ_WorldSnapshot_FPS
{
	float current;
	float window_min;
	float window_max;
	float window_avg;
	int window_samples;
	int limit;
}

class StateZ_WorldSnapshot_Weather
{
	float temperature;
	float temperature_base;
	float wind_speed;
	float wind_direction;
	float rain;
	float snow;
	float fog;
	float clouds;
}

class StateZ_WorldSnapshot_Food
{
	int total;
	int other;
	int fruit;
	int mushroom;
	int corpse;
	int worm;
	int guts;
	int meat;
	int human_meat;
	int disinfectant;
	int pills;
	int drink;
	int snack;
	int candy;
	int canned_small;
	int canned_medium;
	int canned_big;
	int jar;
}

class StateZ_WorldSnapshot_Mining
{
	int bushes;
	int rocks;
	int trees;
}

class StateZ_WorldSnapshot_Fishing
{
	int attempts;
	int catches;
}

// ---------------------------------------------------------------------------
// World snapshot — written to world.json with metadata
// ---------------------------------------------------------------------------

class StateZ_WorldSnapshot
{
	void StateZ_WorldSnapshot()
	{
		fps = new StateZ_WorldSnapshot_FPS();
		weather = new StateZ_WorldSnapshot_Weather();
		food = new StateZ_WorldSnapshot_Food();
		mining = new StateZ_WorldSnapshot_Mining();
		fishing = new StateZ_WorldSnapshot_Fishing();
	}

	// metadata
	string world;

	// timing
	float uptime_seconds;
	int hour;          // in-game hour, 0..23
	int minute;        // in-game minute, 0..59
	int time_of_day;   // 0=dawn, 1=day, 2=dusk, 3=night

	// populations
	int players_online;
	int max_players;
	int infected;
	int infected_dead;
	int animals;
	int animals_dead;
	int path_graph_updates;

	// weapons
	int weapons;
	int suppressors;
	int optics;
	int ammo_boxes;
	int magazines;
	int ammo;
	int explosives;
	int explosives_detonated;

	// items
	int items;
	int bottles;
	int clothes;
	int boxes;
	int cars;
	int boats;
	int helicopters;
	int territory_flags;
	int base_buildings;
	int tents;
	int stashes;
	int gardens;
	int containers;
	int car_parts;
	int car_wheels;

	// events
	int effect_areas;
	int corpses;
	int artillery_barrages;
	int crash_sites;
	int gifts_unpacked;

	// counters
	int chat_messages;
	int players_spawns;
	int players_deaths;
	int infected_deaths;
	int animals_deaths;
	int cars_destroyed;
	int boats_destroyed;
	int helicopters_destroyed;

#ifdef EXPANSIONMODAI
	int eai;
	int eai_npc;
	int eai_deaths;
	int eai_npc_deaths;
#endif

	// nested sub-objects
	ref StateZ_WorldSnapshot_FPS fps;
	ref StateZ_WorldSnapshot_Weather weather;
	ref StateZ_WorldSnapshot_Food food;
	ref StateZ_WorldSnapshot_Mining mining;
	ref StateZ_WorldSnapshot_Fishing fishing;
}

// ---------------------------------------------------------------------------
// Player snapshot — one entry in players.json array
// ---------------------------------------------------------------------------

class StateZ_PlayerSnapshot
{
	string steam_id;
	string guid;
	string name;
	string blood_type;
	string gender;
	string character_type;

	bool is_loaded;
	float health;
	float blood;
	float shock;
	float energy;
	float water;
	float toxicity;
	float temperature_celsius;
	float weight_grams;
	float wetness;
	float lifetime_seconds;

	float pos_x;
	float pos_y;
	float pos_alt;
	float pos_angle;

	int ping_min;
	int ping_max;
	float throttle;

	int agents_active;
	int bleeding_sources;
	bool immunity_boosted;
	bool is_unconscious;
	bool is_restrained;
	bool godmode;

	float stat_playtime_seconds;
	float stat_distance_meters;
	float stat_longest_hit_meters;
	int stat_players_killed;
	int stat_infected_killed;
}

// ---------------------------------------------------------------------------
// Transport snapshot — one entry in vehicles.json array
// ---------------------------------------------------------------------------

class StateZ_TransportPartSnapshot
{
	string name;    //!< part class name (e.g. CarBattery, CivSedanWheel, CarRadiator)
	float health;   //!< normalized part health 0..1
}

class StateZ_TransportSnapshot
{
	string class_name;
	string vehicle_type;
	int hash;

	float health;
	int crew_occupied;
	int crew_capacity;
	float speed_ms;
	int rpm;
	int gear;
	bool engine_on;
	bool lights_on;
	float fuel_fraction;
	float coolant_fraction;  //!< radiator coolant 0..1, or -1 when the vehicle has no radiator
	int cargo_items;         //!< total items stored in the vehicle's cargo (excluding parts)

	float pos_x;
	float pos_y;
	float pos_alt;
	float pos_angle;

	ref array<ref StateZ_TransportPartSnapshot> parts; //!< per-part health (attachments)
}

// ---------------------------------------------------------------------------
// Territory snapshot — one entry in territories.json array
// ---------------------------------------------------------------------------

class StateZ_TerritorySnapshot
{
	float pos_x;
	float pos_y;
	float lifetime_fraction;
	float refresher_radius;
}

// ---------------------------------------------------------------------------
// EffectArea snapshot — one entry in areas.json array
// ---------------------------------------------------------------------------

class StateZ_AreaSnapshot
{
	string type;
	string class_name;
	float pos_x;
	float pos_y;
	float radius;
	int insiders;
}

#endif
