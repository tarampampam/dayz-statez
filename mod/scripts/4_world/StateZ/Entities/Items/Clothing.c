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
modded class Clothing
{
	/**
	    \brief Increment clothes gauge on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		StateZ_Storage.s_Clothing += 1;
	}

	/**
	    \brief Decrement clothes gauge on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		StateZ_Storage.s_Clothing -= 1;

		super.EEDelete(parent);
	}
}
#endif
