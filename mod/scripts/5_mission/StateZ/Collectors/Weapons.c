/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Collector for Weapon stats.
*/
class StateZ_CollectorWeapons : StateZ_CollectorBase
{
	override string GetName()
	{
		return "weapons";
	}

	override void Flush()
	{
		string json;
		StateZ_WeaponStats.BuildJson(json);
		StateZ_JsonWriter.Write(GetName(), json);
	}
}
#endif
