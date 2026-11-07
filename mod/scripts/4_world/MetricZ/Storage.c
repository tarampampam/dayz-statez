/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Global storage and updater for world-level data.
    \details Owns plain-field data, event-driven counters, and flush logic.
*/
class MetricZ_Storage
{
	protected static bool s_Initialized; //!< Initialization flag.

	// Core
	static int s_Status;
	static int s_ScrapeSkippedTotal;

	// FPS
	static float s_FPS;
	static float s_FPSMin;
	static float s_FPSMax;
	static float s_FPSAvg;
	static int s_FPSSamples;
	static int s_FPSLimit;

	// Time
	static float s_ServerUptimeSec;
	static int s_TimeUnixSec;
	static int s_TimeOfDay;

	// World totals
	static int s_PlayersOnline;
	static int s_MaxPlayers;
	static int s_Infected;
	static int s_InfectedCorpses;
	static int s_Animals;
	static int s_AnimalsCorpses;
	static int s_PathGraphUpdates;

	// Weapon related
	static int s_Weapons;
	static int s_Suppressors;
	static int s_Optics;
	static int s_AmmoBoxes;
	static int s_Magazines;
	static int s_Ammo;
	static int s_Explosives;
	static int s_Exploded;

	// Items
	static int s_Items;
	static int s_Bottles;
	static int s_Clothing;
	static int s_Boxes;
	static int s_Cars;
	static int s_Boats;
	static int s_Helicopters;
	static int s_TerritoryFlags;
	static int s_BaseBuildings;
	static int s_Tents;
	static int s_Stashes;
	static int s_Gardens;
	static int s_Containers;
	static int s_CarParts;
	static int s_CarWheels;

	// Food
	static int s_Food;
	static int s_FoodOther;
	static int s_FoodFruit;
	static int s_FoodMushroom;
	static int s_FoodCorpse;
	static int s_FoodWorm;
	static int s_FoodGuts;
	static int s_FoodMeat;
	static int s_FoodHumanMeat;
	static int s_FoodDisinfectant;
	static int s_FoodPills;
	static int s_FoodDrink;
	static int s_FoodSnack;
	static int s_FoodCandy;
	static int s_FoodCannedSmall;
	static int s_FoodCannedMedium;
	static int s_FoodCannedBig;
	static int s_FoodCannedJar;

	// Mining
	static int s_MinedBushes;
	static int s_MinedRocks;
	static int s_MinedTrees;

	// Fishing
	static int s_FishingAttempts;
	static int s_FishingCatches;

	// Events
	static int s_EffectAreas;
	static int s_Corpses;
	static int s_Artillery;
	static int s_CrashSites;
	static int s_GiftsUnpacked;

	// Weather
	static float s_Temperature;
	static float s_TemperatureBase;
	static float s_WindSpeed;
	static float s_WindDirection;
	static float s_Rain;
	static float s_Snow;
	static float s_Fog;
	static float s_Clouds;

	// Counters
	static int s_ChatMessages;
	static int s_PlayersSpawns;
	static int s_PlayersDeaths;
	static int s_InfectedDeaths;
	static int s_AnimalsDeaths;
	static int s_CarsDestroys;
	static int s_BoatsDestroys;
	static int s_HelicoptersDestroys;

#ifdef EXPANSIONMODAI
	static int s_ExpansionAI;
	static int s_ExpansionAINPC;
	static int s_ExpansionAIDeaths;
	static int s_ExpansionAINPCDeaths;
#endif

	/**
	    \brief Check storage is initialized
	*/
	static bool IsInitialized()
	{
		return s_Initialized;
	}

	/**
	    \brief One-time initialization.
	*/
	static void Init()
	{
		if (s_Initialized || !MetricZ_Config.IsLoaded())
			return;

		s_Status = 1;
		s_MaxPlayers = MetricZ_Config.Get().max_players;
		s_FPSLimit = MetricZ_Config.Get().fps_limit;

		s_Initialized = true;
	}

	/**
	    \brief Refresh gauge values from game state.
	*/
	static void Update()
	{
		Init();

		// FPS
		s_FPS = MetricZ_FrameMonitor.GetFPS();

		float fpsMin, fpsMax, fpsAvg;
		int fpsSamples;
		MetricZ_FrameMonitor.SnapshotWindow(fpsMin, fpsMax, fpsAvg, fpsSamples);
		s_FPSMin = fpsMin;
		s_FPSMax = fpsMax;
		s_FPSAvg = fpsAvg;
		s_FPSSamples = fpsSamples;

		// Time
		s_ServerUptimeSec = g_Game.GetTickTime();
		s_TimeUnixSec = MetricZ_Time.GameEpochSeconds();

		int phase;
		switch (g_Game.GetMission().GetWorldData().GetDaytime()) {
		case WorldDataDaytime.DAWN:
			phase = 0;
			break;
		case WorldDataDaytime.DAY:
			phase = 1;
			break;
		case WorldDataDaytime.DUSK:
			phase = 2;
			break;
		case WorldDataDaytime.NIGHT:
		default:
			phase = 3;
			break;
		}
		s_TimeOfDay = phase;

		// Players
		array<Man> mans = new array<Man>();
		g_Game.GetPlayers(mans);
		s_PlayersOnline = mans.Count();

		// Path Graph
		s_PathGraphUpdates = g_Game.s_MetricZ_PathGraphUpdates;

		// World temperature
		s_Temperature = MetricZ_WorldTemperature.Get();
		s_TemperatureBase = g_Game.GetMission().GetWorldData().GetBaseEnvTemperature();

		// Weather
		Weather w = g_Game.GetWeather();
		if (w) {
			s_WindSpeed = w.GetWindSpeed();

			float windDirDeg = w.WindDirectionToAngle(w.GetWind()) * Math.RAD2DEG;
			if (windDirDeg < 0)
				windDirDeg += 360;
			s_WindDirection = windDirDeg;

			s_Rain = w.GetRain().GetActual();
			s_Snow = w.GetSnowfall().GetActual();
			s_Fog = w.GetFog().GetActual();
			s_Clouds = w.GetOvercast().GetActual();
		}
	}

	/**
	    \brief Fill a WorldSnapshot from current field values.
	    \param snap Snapshot instance to populate.
	*/
	static void FillSnapshot(MetricZ_WorldSnapshot snap)
	{
		if (!snap)
			return;

		snap.ts = MetricZ_Time.GameEpochSeconds();
		snap.world = MetricZ_Config.Get().settings.world_name_resolved;
		snap.instance_id = MetricZ_Config.Get().settings.instance_id_resolved;
		snap.host = MetricZ_Config.Get().settings.host_name_resolved;

		snap.uptime_seconds = s_ServerUptimeSec;
		snap.game_time_unix = s_TimeUnixSec;
		snap.time_of_day = s_TimeOfDay;

		snap.fps.current = s_FPS;
		snap.fps.window_min = s_FPSMin;
		snap.fps.window_max = s_FPSMax;
		snap.fps.window_avg = s_FPSAvg;
		snap.fps.window_samples = s_FPSSamples;
		snap.fps.limit = s_FPSLimit;

		snap.players_online = s_PlayersOnline;
		snap.max_players = s_MaxPlayers;
		snap.infected = s_Infected;
		snap.infected_dead = s_InfectedCorpses;
		snap.animals = s_Animals;
		snap.animals_dead = s_AnimalsCorpses;
		snap.path_graph_updates = s_PathGraphUpdates;

		snap.weapons = s_Weapons;
		snap.suppressors = s_Suppressors;
		snap.optics = s_Optics;
		snap.ammo_boxes = s_AmmoBoxes;
		snap.magazines = s_Magazines;
		snap.ammo = s_Ammo;
		snap.explosives = s_Explosives;
		snap.explosives_detonated = s_Exploded;

		snap.items = s_Items;
		snap.bottles = s_Bottles;
		snap.clothes = s_Clothing;
		snap.boxes = s_Boxes;
		snap.cars = s_Cars;
		snap.boats = s_Boats;
		snap.helicopters = s_Helicopters;
		snap.territory_flags = s_TerritoryFlags;
		snap.base_buildings = s_BaseBuildings;
		snap.tents = s_Tents;
		snap.stashes = s_Stashes;
		snap.gardens = s_Gardens;
		snap.containers = s_Containers;
		snap.car_parts = s_CarParts;
		snap.car_wheels = s_CarWheels;

		snap.effect_areas = s_EffectAreas;
		snap.corpses = s_Corpses;
		snap.artillery_barrages = s_Artillery;
		snap.crash_sites = s_CrashSites;
		snap.gifts_unpacked = s_GiftsUnpacked;

		snap.chat_messages = s_ChatMessages;
		snap.players_spawns = s_PlayersSpawns;
		snap.players_deaths = s_PlayersDeaths;
		snap.infected_deaths = s_InfectedDeaths;
		snap.animals_deaths = s_AnimalsDeaths;
		snap.cars_destroyed = s_CarsDestroys;
		snap.boats_destroyed = s_BoatsDestroys;
		snap.helicopters_destroyed = s_HelicoptersDestroys;

#ifdef EXPANSIONMODAI
		snap.eai = s_ExpansionAI;
		snap.eai_npc = s_ExpansionAINPC;
		snap.eai_deaths = s_ExpansionAIDeaths;
		snap.eai_npc_deaths = s_ExpansionAINPCDeaths;
#endif

		snap.weather.temperature = s_Temperature;
		snap.weather.temperature_base = s_TemperatureBase;
		snap.weather.wind_speed = s_WindSpeed;
		snap.weather.wind_direction = s_WindDirection;
		snap.weather.rain = s_Rain;
		snap.weather.snow = s_Snow;
		snap.weather.fog = s_Fog;
		snap.weather.clouds = s_Clouds;

		snap.food.total = s_Food;
		snap.food.other = s_FoodOther;
		snap.food.fruit = s_FoodFruit;
		snap.food.mushroom = s_FoodMushroom;
		snap.food.corpse = s_FoodCorpse;
		snap.food.worm = s_FoodWorm;
		snap.food.guts = s_FoodGuts;
		snap.food.meat = s_FoodMeat;
		snap.food.human_meat = s_FoodHumanMeat;
		snap.food.disinfectant = s_FoodDisinfectant;
		snap.food.pills = s_FoodPills;
		snap.food.drink = s_FoodDrink;
		snap.food.snack = s_FoodSnack;
		snap.food.candy = s_FoodCandy;
		snap.food.canned_small = s_FoodCannedSmall;
		snap.food.canned_medium = s_FoodCannedMedium;
		snap.food.canned_big = s_FoodCannedBig;
		snap.food.jar = s_FoodCannedJar;

		snap.mining.bushes = s_MinedBushes;
		snap.mining.rocks = s_MinedRocks;
		snap.mining.trees = s_MinedTrees;

		snap.fishing.attempts = s_FishingAttempts;
		snap.fishing.catches = s_FishingCatches;
	}

	/**
	    \brief Adjust per-type food sub-counter on spawn/delete.
	    \details Does not touch s_Food total — caller handles that.
	    \param foodType Food type enum value.
	    \param increase true on spawn, false on delete.
	*/
	static void FoodMetricChange(MetricZ_FoodTypes foodType, bool increase)
	{
		int delta = 1;
		if (!increase)
			delta = -1;

		switch (foodType) {
		case MetricZ_FoodTypes.OTHER:
			s_FoodOther += delta;
			break;
		case MetricZ_FoodTypes.FRUIT:
			s_FoodFruit += delta;
			break;
		case MetricZ_FoodTypes.MUSHROOM:
			s_FoodMushroom += delta;
			break;
		case MetricZ_FoodTypes.CORPSE:
			s_FoodCorpse += delta;
			break;
		case MetricZ_FoodTypes.WORM:
			s_FoodWorm += delta;
			break;
		case MetricZ_FoodTypes.GUTS:
			s_FoodGuts += delta;
			break;
		case MetricZ_FoodTypes.MEAT:
			s_FoodMeat += delta;
			break;
		case MetricZ_FoodTypes.HUMAN_MEAT:
			s_FoodHumanMeat += delta;
			break;
		case MetricZ_FoodTypes.DISINFECTANT:
			s_FoodDisinfectant += delta;
			break;
		case MetricZ_FoodTypes.PILLS:
			s_FoodPills += delta;
			break;
		case MetricZ_FoodTypes.DRINK:
			s_FoodDrink += delta;
			break;
		case MetricZ_FoodTypes.SNACK:
			s_FoodSnack += delta;
			break;
		case MetricZ_FoodTypes.CANDY:
			s_FoodCandy += delta;
			break;
		case MetricZ_FoodTypes.CANNED_SMALL:
			s_FoodCannedSmall += delta;
			break;
		case MetricZ_FoodTypes.CANNED_MEDIUM:
			s_FoodCannedMedium += delta;
			break;
		case MetricZ_FoodTypes.CANNED_BIG:
			s_FoodCannedBig += delta;
			break;
		case MetricZ_FoodTypes.JAR:
			s_FoodCannedJar += delta;
			break;
		default:
			break;
		}
	}
}
#endif
