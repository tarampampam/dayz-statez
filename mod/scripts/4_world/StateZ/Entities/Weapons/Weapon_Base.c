/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Updates StateZ counters on spawn/cleanup and weapon shots.
*/
modded class Weapon_Base
{
	// cache for weapon label names: original type -> canonical name.
	protected static ref map<string, string> s_StateZ_LabelNames = new map<string, string>();

	/**
	    \brief Increment weapons gauge on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		if (!StateZ_Config.IsLoaded())
			return;

		if (!StateZ_Config.Get().disabled_metrics.weapons)
			StateZ_WeaponStats.OnSpawn(this);

		StateZ_Storage.s_Weapons += 1;
	}

	/**
	    \brief Decrement weapons gauge on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (StateZ_Config.IsLoaded()) {
			StateZ_Storage.s_Weapons -= 1;

			if (!StateZ_Config.Get().disabled_metrics.weapons)
				StateZ_WeaponStats.OnDelete(this);
		}

		super.EEDelete(parent);
	}

	/**
	    \brief Count weapon shot for stats.
	    \details Calls base then forwards to StateZ_WeaponStats unless disabled.
	*/
	override void OnFire(int muzzle_index)
	{
		super.OnFire(muzzle_index);

		if (!StateZ_Config.IsLoaded())
			return;

		if (!StateZ_Config.Get().disabled_metrics.weapons)
			StateZ_WeaponStats.OnFire(this);
	}

	/**
	    \brief Public helper for StateZ: returns cached canonical label name for this weapon.
	    \details You can override this for set some beauty label name for your weapon.
	*/
	override string StateZ_GetLabelTypeName()
	{
		return StateZ_GetWeaponNameByType();
	}

	/**
	    \brief Compute canonical weapon name from type and cache the result.
	    \details
	      - remove any `sawedoff` token (case-insensitive, anywhere)
	      - collapse `__` -> `_`, then trim leading/trailing `_`
	      - cut suffix after last `_`
	      - lowercase; fallback to original lowercase if empty
	*/
	protected string StateZ_GetWeaponNameByType()
	{
		string type = GetType();

		// cache hit
		string cached;
		if (s_StateZ_LabelNames && s_StateZ_LabelNames.Find(type, cached))
			return cached;

		string result = StateZ_ObjectName.GetName(this, true);
		result.Replace("sawedoff", "");

		// cache store
		if (!s_StateZ_LabelNames)
			s_StateZ_LabelNames = new map<string, string>();
		s_StateZ_LabelNames.Set(type, result);

		return result;
	}
}
#endif
