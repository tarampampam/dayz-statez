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
modded class Ammunition_Base
{
	/**
	    \brief Increment ammo gauge on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		StateZ_Storage.s_Ammo += 1;
	}

	/**
	    \brief Decrement ammo gauge on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		StateZ_Storage.s_Ammo -= 1;

		super.EEDelete(parent);
	}
}
#endif
