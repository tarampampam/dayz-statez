/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER

class StateZ_WeaponStatsDTO
{
	void StateZ_WeaponStatsDTO()
	{
		shots = new map<string, int>();
		world_count = new map<string, int>();
		player_kills = new map<string, int>();
		creature_kills = new map<string, int>();
	}

	ref map<string, int> shots;
	ref map<string, int> world_count;
	ref map<string, int> player_kills;
	ref map<string, int> creature_kills;
}

/**
    \brief Weapon shot and live-count aggregator.
    \details
      - counts total shots and per-weapon shots
      - tracks live weapons count per canonical type
      Weapon key is canonicalized and lowercased via Weapon_Base::StateZ_GetLabelTypeName().
*/
class StateZ_WeaponStats
{
	protected static ref map<string, int> s_ShotsByWeapon = new map<string, int>(); //!< Total shots by weapon.
	protected static ref map<string, int> s_CountByType = new map<string, int>(); //!< Live count by weapon type.
	protected static ref map<string, int> s_PlayerKills = new map<string, int>(); //!< Player kills by source.
	protected static ref map<string, int> s_CreatureKills = new map<string, int>(); //!< Creature kills by source.

	/**
	    \brief Increment counters for a fired weapon.
	    \param weapon Weapon_Base instance.
	*/
	static void OnFire(Weapon_Base weapon)
	{
		if (!weapon)
			return;

		IncMap(s_ShotsByWeapon, weapon.StateZ_GetLabelTypeName());
	}

	/**
	    \brief Increment per-type live count for spawned weapon.
	    \param weapon Weapon_Base instance.
	*/
	static void OnSpawn(Weapon_Base weapon)
	{
		if (!weapon)
			return;

		IncMap(s_CountByType, weapon.StateZ_GetLabelTypeName());
	}

	/**
	    \brief Called when a Player is killed.
	    \param source Object killer instance.
	*/
	static void OnPlayerKilled(Object source)
	{
		if (!source)
			return;

		IncMap(s_PlayerKills, ResolveSourceName(source));
	}

	/**
	    \brief Called when a Creature (Zombie/Animal/eAI) is killed.
	    \param source Object killer instance.
	*/
	static void OnCreatureKilled(Object source)
	{
		if (!source)
			return;

		IncMap(s_CreatureKills, ResolveSourceName(source));
	}

	/**
	    \brief Decrement per-type live count for deleted weapon.
	    \param weapon Weapon_Base instance.
	*/
	static void OnDelete(Weapon_Base weapon)
	{
		if (!weapon)
			return;

		string type = weapon.StateZ_GetLabelTypeName();

		int v;
		if (!s_CountByType.Find(type, v))
			return;

		v = v - 1;
		if (v <= 0)
			s_CountByType.Remove(type);
		else
			s_CountByType.Set(type, v);
	}

	/**
	    \brief Serialize all weapon data to JSON.
	    \param json Output JSON string.
	*/
	static void BuildJson(out string json)
	{
		StateZ_WeaponStatsDTO dto = new StateZ_WeaponStatsDTO();
		dto.shots = s_ShotsByWeapon;
		dto.world_count = s_CountByType;
		dto.player_kills = s_PlayerKills;
		dto.creature_kills = s_CreatureKills;
		JsonSerializer().WriteToString(dto, false, json);
	}

	/**
	    \brief Determines the metric key from an Object source.
	    \param source Killer Object.
	    \return string Normalized key (e.g., "akm", "zombie", "player_melee_item_heavy").
	*/
	protected static string ResolveSourceName(Object source)
	{
		if (!source)
			return "none";

		Weapon_Base wpn;
		if (Class.CastTo(wpn, source))
			return wpn.StateZ_GetLabelTypeName();

		if (source.IsInherited(ZombieBase))
			return "zombie";

		if (source.IsInherited(AnimalBase))
			return "animal";

		if (source.IsInherited(Transport))
			return "transport";

		if (source.IsInherited(Grenade_Base))
			return "grenade";

		if (source.IsInherited(TrapBase))
			return "trap";

		if (source.IsInherited(ToolBase))
			return "tool";

		if (source.IsInherited(PlayerBase)) {
#ifdef EXPANSIONMODAI
			if (source.IsInherited(eAIBase))
				return "eai_fist";
#endif
			return "player_fist";
		}

		EntityAI eai;
		if (!Class.CastTo(eai, source))
			return "object";

		Man rootMan = eai.GetHierarchyRootPlayer();
		if (rootMan) {
			string type;
			if (eai.IsHeavyBehaviour())
				type = "heavy";
			else if (eai.IsTwoHandedBehaviour())
				type = "two_hands";
			else if (eai.IsOneHandedBehaviour())
				type = "one_hand";

#ifdef EXPANSIONMODAI
			if (rootMan.IsInherited(eAIBase))
				return string.Format("eai_melee_item_%1", type);
#endif
			return string.Format("player_melee_item_%1", type);
		}

		return "unknown";
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
