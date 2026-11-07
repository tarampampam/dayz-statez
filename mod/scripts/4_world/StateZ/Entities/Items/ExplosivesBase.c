/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Updates StateZ counters on spawn/cleanup.
*/
modded class ExplosivesBase
{
	/**
	    \brief Increment explosives gauge on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		StateZ_Storage.s_Explosives += 1;
	}

	/**
	    \brief Decrement explosives gauge on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		StateZ_Storage.s_Explosives -= 1;

		super.EEDelete(parent);
	}

	/**
	    \brief Increment detonations counter on explode.
	*/
	override protected void OnExplode()
	{
		super.OnExplode();

		StateZ_Storage.s_Exploded += 1;
	}
}
#endif
