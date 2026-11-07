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
modded class AnimalBase
{
	// prevent counting kills on already killed body
	protected bool m_StateZ_IsKilled;
	// prevent counting hits on death bodies except last hit
	protected bool m_StateZ_IsLastHit;
	// cache for canonical animal type names: original -> canonical
	protected static ref map<string, string> s_StateZ_TypeNames = new map<string, string>();

	/**
	    \brief Increment animal count on entity init.
	*/
	override void EEInit()
	{
		super.EEInit();

		if (!StateZ_Config.IsLoaded())
			return;

		if (!StateZ_Config.Get().disabled_metrics.animals)
			StateZ_AnimalStats.OnSpawn(this);

		StateZ_Storage.s_Animals += 1;
	}

	/**
	    \brief Decrement animal count on entity delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (StateZ_Config.IsLoaded()) {
			StateZ_Storage.s_Animals -= 1;

			if (m_StateZ_IsKilled)
				StateZ_Storage.s_AnimalsCorpses -= 1;

			if (!StateZ_Config.Get().disabled_metrics.animals)
				StateZ_AnimalStats.OnDelete(this);
		}

		super.EEDelete(parent);
	}

	/**
	    \brief Increment animal deaths counter on kill.
	*/
	override void EEKilled(Object killer)
	{
		if (StateZ_Config.IsLoaded() && !m_StateZ_IsKilled) {
			m_StateZ_IsKilled = true;

			if (killer != this)
				StateZ_Storage.s_AnimalsDeaths += 1;

			StateZ_Storage.s_AnimalsCorpses += 1;

			if (!StateZ_Config.Get().disabled_metrics.weapons && killer != this)
				StateZ_WeaponStats.OnCreatureKilled(killer);
		}

		super.EEKilled(killer);
	}

	/**
	    \brief Capture hit statistics for animals.
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
	    \brief Public helper for StateZ: returns cached canonical type name for this animal.
	    \details You can override this for set some beauty type name for your animal.
	*/
	string StateZ_GetTypeName()
	{
		return StateZ_GetAnimalNameFromConfig();
	}

	/**
	    \brief Resolve canonical animal name using config (Steaks -> Pelt -> type) and cache it.
	    \details
	      - Try Skinning.ObtainedSteaks.item and strip:
	          "SteakMeat", "BreastMeat", "LegMeat"
	      - If empty, try Skinning.ObtainedPelt.item and strip "Pelt"
	      - If still empty, use GetType()
	      - Result is lowercased and cached per animal type
	*/
	protected string StateZ_GetAnimalNameFromConfig()
	{
		string type = GetType();

		// cache hit
		string cached;
		if (s_StateZ_TypeNames && s_StateZ_TypeNames.Find(type, cached))
			return cached;

		string cfgSteaks = string.Format("%1 %2 Skinning ObtainedSteaks item", CFG_VEHICLESPATH, type);
		string cfgPelts = string.Format("%1 %2 Skinning ObtainedPelt item", CFG_VEHICLESPATH, type);
		string result;

		// 1) Steaks
		if (g_Game.ConfigIsExisting(cfgSteaks)) {
			string steakItem;
			g_Game.ConfigGetText(cfgSteaks, steakItem);
			steakItem.TrimInPlace();

			if (steakItem != string.Empty) {
				// try supported meat suffixes
				string nameCandidate;

				nameCandidate = StateZ_GetAnimalNameByItem(steakItem, "SteakMeat");
				if (nameCandidate == string.Empty)
					nameCandidate = StateZ_GetAnimalNameByItem(steakItem, "BreastMeat");
				if (nameCandidate == string.Empty)
					nameCandidate = StateZ_GetAnimalNameByItem(steakItem, "LegMeat");

				if (nameCandidate != string.Empty)
					result = nameCandidate;
			}
		}

		// 2) Pelt fallback
		if (result == string.Empty && g_Game.ConfigIsExisting(cfgPelts)) {
			string peltItem;
			g_Game.ConfigGetText(cfgPelts, peltItem);
			peltItem.TrimInPlace();

			if (peltItem != string.Empty) {
				string peltName = StateZ_GetAnimalNameByItem(peltItem, "Pelt");
				if (peltName != string.Empty)
					result = peltName;
			}
		}

		// 3) Last fallback
		if (result == string.Empty)
			result = type;

		result.ToLower();

		// cache store
		if (!s_StateZ_TypeNames)
			s_StateZ_TypeNames = new map<string, string>();
		s_StateZ_TypeNames.Set(type, result);

		return result;
	}

	/**
	    \brief Extract base animal name from item class name using suffix.
	    \param item Item classname from config (e.g. "CowSteakMeat")
	    \param suffix Expected suffix (e.g. "SteakMeat", "BreastMeat", "LegMeat", "Pelt")
	    \return Base name without suffix, or empty string on mismatch.
	*/
	protected string StateZ_GetAnimalNameByItem(string item, string suffix)
	{
		item.TrimInPlace();
		if (item == string.Empty)
			return "";

		int lenItem = item.Length();
		int lenSuf = suffix.Length();

		if (lenSuf == 0 || lenSuf > lenItem)
			return "";

		int pos = item.LastIndexOf(suffix);
		if (pos == -1)
			return "";

		// suffix must be at the end
		if (pos + lenSuf != lenItem)
			return "";

		return item.Substring(0, pos);
	}
}
#endif
