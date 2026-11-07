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
		thresholds = new StateZ_ConfigDTO_Thresholds();
	}

	// Interval in seconds between data collection cycles.
	int interval = 15;

	// Delete all JSON files in the data directory on server shutdown.
	// Set to false to keep the files for inspection after the server stops.
	bool delete_on_shutdown = false;

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
		thresholds.Normalize();

		interval = (int)Math.Clamp(interval, 0, 900);
		max_players = StateZ_Helpers.GetLimitPlayers();
		fps_limit = StateZ_Helpers.GetLimitFPS();
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
