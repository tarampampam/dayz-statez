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
modded class Box_Base
{
	/**
	    \brief Increment boxes or ammo_boxes count on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		if (GetType().IndexOf("AmmoBox") != -1)
			StateZ_Storage.s_AmmoBoxes += 1;
		else
			StateZ_Storage.s_Boxes += 1;
	}

	/**
	    \brief Decrement boxes or ammo_boxes count on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (GetType().IndexOf("AmmoBox") != -1)
			StateZ_Storage.s_AmmoBoxes -= 1;
		else
			StateZ_Storage.s_Boxes -= 1;

		super.EEDelete(parent);
	}
}
#endif
