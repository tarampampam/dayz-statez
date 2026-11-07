/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief StateZ Config Data Transfer Object
*/
class StateZ_ConfigDTO
{
	void StateZ_ConfigDTO()
	{
		settings = new StateZ_ConfigDTO_BaseSettings();
		thresholds = new StateZ_ConfigDTO_Thresholds();
	}

	// Internal configuration version. **Do not modify**.
	string version = StateZ_Constants.VERSION;

	// Base settings for metric collection.
	ref StateZ_ConfigDTO_BaseSettings settings;

	// Collection thresholds.
	ref StateZ_ConfigDTO_Thresholds thresholds;

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
		thresholds.Normalize();

		max_players = StateZ_Helpers.GetLimitPlayers();
		fps_limit = StateZ_Helpers.GetLimitFPS();
	}
}

/**
    \brief Base settings for data collection.
*/
class StateZ_ConfigDTO_BaseSettings
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
		instance_id_resolved = StateZ_Helpers.GetInstanceID(instance_id);

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
    \brief Collection thresholds.
*/
class StateZ_ConfigDTO_Thresholds
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
