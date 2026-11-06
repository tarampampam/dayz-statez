/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief MetricZ Config Data Transfer Object
*/
class MetricZ_ConfigDTO
{
	void MetricZ_ConfigDTO()
	{
		settings = new MetricZ_ConfigDTO_BaseSettings();
		json = new MetricZ_ConfigDTO_JsonExport();
		disabled_metrics = new MetricZ_ConfigDTO_DisabledMetrics();
		thresholds = new MetricZ_ConfigDTO_Thresholds();
	}

	// Internal configuration version. **Do not modify**.
	string version = MetricZ_Constants.VERSION;

	// Base settings for metric collection.
	ref MetricZ_ConfigDTO_BaseSettings settings;

	// Settings for exporting data to local JSON files.
	ref MetricZ_ConfigDTO_JsonExport json;

	// Switches to disable specific data collection.
	ref MetricZ_ConfigDTO_DisabledMetrics disabled_metrics;

	// Collection thresholds.
	ref MetricZ_ConfigDTO_Thresholds thresholds;

	[NonSerialized()]
	int max_players = 255; // from serverDZ.cfg:maxPlayers

	[NonSerialized()]
	int fps_limit = 0; // from CLI:-limitFPS

	/**
	    \brief Normalizes configuration values within valid ranges.
	*/
	void Normalize()
	{
		settings.Normalize();
		json.Normalize();
		disabled_metrics.Normalize();
		thresholds.Normalize();

		max_players = MetricZ_Helpers.GetLimitPlayers();
		fps_limit = MetricZ_Helpers.GetLimitFPS();
	}
}

/**
    \brief Base settings for data collection.
*/
class MetricZ_ConfigDTO_BaseSettings
{
	// Overrides the Instance ID. By default, this is detected automatically from serverDZ.cfg instanceID,
	// or uses gamePort/steamQueryPort if instanceID is undefined or zero.
	// This ID must be unique across all game servers on all hosts.
	string instance_id;

	// Overrides the Host Name. By default, attempts to obtain it automatically.
	string host_name;

	// Delay in seconds before the first data collection begins.
	int init_delay_sec = 30;

	// Interval in seconds between data collection cycles.
	int collect_interval_sec = 15;

	[NonSerialized()]
	string instance_id_resolved;

	[NonSerialized()]
	string host_name_resolved;

	[NonSerialized()]
	string world_name_resolved;

	/**
	    \brief Normalizes configuration values within valid ranges.
	*/
	void Normalize()
	{
		instance_id_resolved = MetricZ_Helpers.GetInstanceID(instance_id);

		if (host_name != string.Empty)
			host_name_resolved = host_name;
		else {
			string host = GetMachineName();
			host.TrimInPlace();
			host.ToLower();
			host_name_resolved = host;
		}

		string wn;
		g_Game.GetWorldName(wn);
		wn.TrimInPlace();
		wn.ToLower();
		world_name_resolved = wn;

		init_delay_sec = (int)Math.Clamp(init_delay_sec, 0, 300);
		collect_interval_sec = (int)Math.Clamp(collect_interval_sec, 0, 900);
	}
}

/**
    \brief Settings for exporting data to local JSON files.
*/
class MetricZ_ConfigDTO_JsonExport
{
	// Enables saving data to local JSON files.
	bool enabled = true;

	// Enables atomic export. Data is written to a temporary file and copied/renamed
	// to the .json file upon completion, removing the temp file.
	// Prevents consumers from reading incomplete files.
	bool atomic = true;

	// File name override (without extension).
	// Default path: $profile:metricz/export/metricz_${instance_id}/<domain>.json
	// Override path: $profile:metricz/export/${file_name}/<domain>.json
	string file_name;

	// Delete JSON files when shutting down the server.
	bool delete_on_shutdown;

	[NonSerialized()]
	string export_dir;

	/**
	    \brief Normalizes configuration values within valid ranges.
	*/
	void Normalize() {}
}

/**
    \brief Switches to disable specific data collection.
*/
class MetricZ_ConfigDTO_DisabledMetrics
{
	// Disables RPC metrics collection.
	bool rpc_input;

	// Disables event handler metrics collection.
	bool events;

	// Disables player-related data collection.
	bool players;

	// Disables zombie per-type and mind state collection.
	bool zombies;

	// Disables animal per-type collection.
	bool animals;

	// Disables vehicle and transport data collection.
	bool transports;

	// Disables weapon per-type collection (count, shots, kills).
	bool weapons;

	// Disables hit metrics by ammo type.
	bool hits;

	// Disables kill metrics by object type.
	bool kills;

	// Disables territory flag data collection.
	bool territories;

	// Disables EffectArea data collection.
	bool areas;

	// Disables Local EffectArea collection.
	bool local_areas = true;

	// Disables player and transport coordinate collection.
	bool positions;

	// Disables height coordinate (Y) collection.
	// If `positions` is disabled, this will also be disabled forcibly.
	bool positions_height = true;

	// Disables orientation collection.
	// If `positions` is disabled, this will also be disabled forcibly.
	bool positions_yaw = true;

	/**
	    \brief Normalizes configuration values within valid ranges.
	*/
	void Normalize() {}
}

/**
    \brief Collection thresholds.
*/
class MetricZ_ConfigDTO_Thresholds
{
	// Minimum damage required to collect hit data in `EEHitBy()`.
	// Values of -1 or less disable this threshold.
	float hit_damage = 3;

	// Minimum damage from vehicles required to collect hit data in `EEHitBy()`.
	// Values of -1 or less disable this threshold.
	float hit_damage_vehicle = 15;

	/**
	    \brief Normalizes configuration values within valid ranges.
	*/
	void Normalize()
	{
		hit_damage = Math.Clamp(hit_damage, -1, 100);
		hit_damage_vehicle = Math.Clamp(hit_damage_vehicle, -15, 100);
	}
}
#endif
