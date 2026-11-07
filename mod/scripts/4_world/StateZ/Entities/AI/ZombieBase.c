/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Updates StateZ counters on spawn/cleanup and kills.
*/
modded class ZombieBase
{
	// prevent counting kills on already killed body
	protected bool m_StateZ_IsKilled;
	// prevent counting hits on death bodies except last hit
	protected bool m_StateZ_IsLastHit;
	// last reported mind state
	protected int m_StateZ_State = -1;
	// cache for canonical zombie type names: original class -> canonical type
	protected static ref map<string, string> s_StateZ_TypeNames = new map<string, string>();

	/**
	    \brief Increment infected count on entity init.
	*/
	override void EEInit()
	{
		super.EEInit();

		if (!StateZ_Config.IsLoaded())
			return;

		StateZ_Storage.s_Infected += 1;

		if (!StateZ_Config.Get().disabled_metrics.zombies)
			StateZ_ZombieStats.OnSpawn(this);
	}

	/**
	    \brief Decrement infected count and update mind-state stats on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (StateZ_Config.IsLoaded()) {
			StateZ_Storage.s_Infected -= 1;

			if (m_StateZ_IsKilled)
				StateZ_Storage.s_InfectedCorpses -= 1;

			if (!StateZ_Config.Get().disabled_metrics.zombies)
				StateZ_ZombieStats.OnDelete(this, m_StateZ_State);
		}

		super.EEDelete(parent);
	}

	/**
	    \brief Increment infected deaths counter on kill and reset cached state.
	*/
	override void EEKilled(Object killer)
	{
		if (StateZ_Config.IsLoaded() && !m_StateZ_IsKilled) {
			m_StateZ_IsKilled = true;

			if (killer != this)
				StateZ_Storage.s_InfectedDeaths += 1;

			StateZ_Storage.s_InfectedCorpses += 1;

			if (!StateZ_Config.Get().disabled_metrics.zombies) {
				StateZ_ZombieStats.OnKilled(m_StateZ_State);
				m_StateZ_State = StateZ_ZombieStats.MINDSTATE_DEAD;
			} else
				m_StateZ_State = -1;

			if (!StateZ_Config.Get().disabled_metrics.weapons && killer != this)
				StateZ_WeaponStats.OnCreatureKilled(killer);
		}

		super.EEKilled(killer);
	}

	/**
	    \brief Capture hit statistics for infected.
	*/
	override void EEHitBy(TotalDamageResult damageResult, int damageType, EntityAI source, int component, string dmgZone, string ammo, vector modelPos, float speedCoef)
	{
		super.EEHitBy(damageResult, damageType, source, component, dmgZone, ammo, modelPos, speedCoef);

		if (!StateZ_Config.IsLoaded())
			return;

		if (!StateZ_Config.Get().disabled_metrics.hits && source != this && !m_StateZ_IsLastHit) {
			if (IsDamageDestroyed())
				m_StateZ_IsLastHit = true;

			if (damageResult) {
				float damage = damageResult.GetDamage(dmgZone, "");
				if (damage < StateZ_Config.Get().thresholds.hit_damage)
					return;

				if (source && source.IsTransport() && damage < StateZ_Config.Get().thresholds.hit_damage_vehicle)
					return;
			}

			StateZ_HitStats.OnCreatureHit(ammo);
		}
	}

	/**
	    \brief Track mind state transitions for infected.
	    \details When state changes, decrement old and increment new in StateZ_ZombieStats.
	    \return \p bool Base return value.
	*/
	override bool HandleMindStateChange(int pCurrentCommandID, DayZInfectedInputController pInputController, float pDt)
	{
		bool ret = super.HandleMindStateChange(pCurrentCommandID, pInputController, pDt);

		if (m_MindState != m_StateZ_State) {
			StateZ_ZombieStats.OnStateChange(m_StateZ_State, m_MindState);
			m_StateZ_State = m_MindState;
		}

		return ret;
	}

	/**
	    \brief Public helper for StateZ: returns cached canonical zombie type.
	    \details You can override this for set some beauty type name for your zombie.
	*/
	string StateZ_GetTypeName()
	{
		return StateZ_GetZombieTypeFromConfig();
	}

	/**
	    \brief Resolve canonical zombie type using ResistContaminatedEffect and aiAgentTemplate.
	    \return \p string Canonical zombie type, or empty string on failure.
	*/
	protected string StateZ_GetZombieTypeFromConfig()
	{
		// contamination has highest priority
		if (ResistContaminatedEffect())
			return "nbc";

		string type = GetType();

		// cache hit
		string cached;
		if (s_StateZ_TypeNames && s_StateZ_TypeNames.Find(type, cached))
			return cached;

		string result;

		// read aiAgentTemplate from config
		string cfgPath = string.Format("%1 %2 aiAgentTemplate", CFG_VEHICLESPATH, type);
		string tpl;

		if (g_Game.ConfigIsExisting(cfgPath)) {
			g_Game.ConfigGetText(cfgPath, tpl);
			tpl.TrimInPlace();

			if (tpl != string.Empty) {
				if (tpl == "Infected" || tpl == "InfectedFemale" || tpl == "InfectedMale") {
					if (IsZombieMilitary())
						result = "military";
					else {
						if (tpl == "InfectedFemale")
							result = "default_female";
						else if (tpl == "InfectedMale")
							result = "default_male";
						else
							result = "default";
					}
				} else {
					// non-base templates: strip "Infected" prefix if present
					string prefix = "Infected";
					int plen = prefix.Length();

					if (tpl.Length() >= plen && tpl.Substring(0, plen) == prefix)
						result = tpl.Substring(plen, tpl.Length() - plen);
					else {
						result = tpl; // custom template without prefix
					}
				}
			}
		}

		// fallback path if template missing/empty
		if (result == string.Empty) {
			if (IsZombieMilitary())
				result = "military";
			else
				result = type;
		}

		result.TrimInPlace();
		result.ToLower();

		if (result == string.Empty) {
			result = type;
			result.TrimInPlace();
			result.ToLower();
		}

		// cache store
		if (!s_StateZ_TypeNames)
			s_StateZ_TypeNames = new map<string, string>();
		s_StateZ_TypeNames.Set(type, result);

		return result;
	}
}
#endif
