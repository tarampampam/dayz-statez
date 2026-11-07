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
modded class CarDoor
{
	/**
	    \brief Increment car_parts count on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		StateZ_Storage.s_CarParts += 1;
	}

	/**
	    \brief Decrement car_parts count on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		StateZ_Storage.s_CarParts -= 1;

		super.EEDelete(parent);
	}
}
#endif
