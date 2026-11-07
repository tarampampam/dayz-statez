/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Updates MetricZ counters on spawn/cleanup.
*/
modded class MagazineStorage
{
	/**
	    \brief Increment magazines gauge on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		MetricZ_Storage.s_Magazines += 1;
	}

	/**
	    \brief Decrement magazines gauge on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		MetricZ_Storage.s_Magazines -= 1;

		super.EEDelete(parent);
	}
}
#endif
