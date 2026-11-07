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
modded class ItemBase
{
	/**
	    \brief Increment items count on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		StateZ_Storage.s_Items += 1;

		if (IsInherited(ItemSuppressor))
			StateZ_Storage.s_Suppressors += 1;
		else if (IsInherited(ItemOptics))
			StateZ_Storage.s_Optics += 1;
		else if (IsInherited(CarWheel))
			StateZ_Storage.s_CarWheels += 1;
	}

	/**
	    \brief Decrement items count on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		StateZ_Storage.s_Items -= 1;

		if (IsInherited(ItemSuppressor))
			StateZ_Storage.s_Suppressors -= 1;
		else if (IsInherited(ItemOptics))
			StateZ_Storage.s_Optics -= 1;
		else if (IsInherited(CarWheel))
			StateZ_Storage.s_CarWheels -= 1;

		super.EEDelete(parent);
	}

	/**
	    \brief Public helper for StateZ: returns cached canonical type name for this item.
	    \details You can override this for set some beauty type name for your item.
	*/
	string StateZ_GetTypeName()
	{
		return StateZ_ObjectName.GetName(this, true);
	}
}
#endif
