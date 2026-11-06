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
modded class ItemBase
{
	/**
	    \brief Increment items gauge on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		MetricZ_Storage.s_Items += 1;

		if (IsInherited(ItemSuppressor))
			MetricZ_Storage.s_Suppressors += 1;
		else if (IsInherited(ItemOptics))
			MetricZ_Storage.s_Optics += 1;
		else if (IsInherited(CarWheel))
			MetricZ_Storage.s_CarWheels += 1;
	}

	/**
	    \brief Decrement items gauge on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		MetricZ_Storage.s_Items -= 1;

		if (IsInherited(ItemSuppressor))
			MetricZ_Storage.s_Suppressors -= 1;
		else if (IsInherited(ItemOptics))
			MetricZ_Storage.s_Optics -= 1;
		else if (IsInherited(CarWheel))
			MetricZ_Storage.s_CarWheels -= 1;

		super.EEDelete(parent);
	}

	/**
	    \brief Public helper for MetricZ: returns cached canonical label name for this item.
	    \details You can override this for set some beauty label name for your item.
	*/
	string MetricZ_GetLabelTypeName()
	{
		return MetricZ_ObjectName.GetName(this, true);
	}
}
#endif
