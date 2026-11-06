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
modded class BaseBuildingBase
{
	/**
	    \brief Increment base_buildings gauge on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		MetricZ_Storage.s_BaseBuildings += 1;
	}

	/**
	    \brief Decrement base_buildings gauge on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		MetricZ_Storage.s_BaseBuildings -= 1;

		super.EEDelete(parent);
	}
}
#endif
