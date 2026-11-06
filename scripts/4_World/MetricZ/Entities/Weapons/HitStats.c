/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER

class MetricZ_HitStatsDTO
{
	void MetricZ_HitStatsDTO()
	{
		player_hit = new map<string, int>();
		creature_hit = new map<string, int>();
	}

	ref map<string, int> player_hit;
	ref map<string, int> creature_hit;
}

/**
    \brief Damage counters aggregator.
    \details Tracks hits by ammo type for players and creatures.
*/
class MetricZ_HitStats
{
	protected static ref map<string, int> s_PlayerHit = new map<string, int>(); //!< Player hits by ammo type.
	protected static ref map<string, int> s_CreatureHit = new map<string, int>(); //!< Creature hits by ammo type.

	/**
	    \brief Register a hit on a Player.
	    \param ammo Ammo type.
	*/
	static void OnPlayerHit(string ammo)
	{
		if (ammo == string.Empty)
			return;

		IncMap(s_PlayerHit, ammo);
	}

	/**
	    \brief Register a hit on a Creature.
	    \param ammo Ammo type.
	*/
	static void OnCreatureHit(string ammo)
	{
		if (ammo == string.Empty)
			return;

		IncMap(s_CreatureHit, ammo);
	}

	/**
	    \brief Serialize all hit data to JSON.
	    \param json Output JSON string.
	*/
	static void BuildJson(out string json)
	{
		MetricZ_HitStatsDTO dto = new MetricZ_HitStatsDTO();
		dto.player_hit = s_PlayerHit;
		dto.creature_hit = s_CreatureHit;
		JsonSerializer().WriteToString(dto, false, json);
	}

	/**
	    \brief Helper to increment value in a map.
	*/
	protected static void IncMap(map<string, int> store, string key)
	{
		int value;
		if (store.Find(key, value)) {
			store.Set(key, value + 1);
			return;
		}

		store.Insert(key, 1);
	}
}
#endif
