/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Collector for Zombie stats.
*/
class StateZ_CollectorZombies : StateZ_CollectorBase
{
	override string GetName()
	{
		return "zombies";
	}

	override void Flush()
	{
		string json;
		StateZ_ZombieStats.BuildJson(json);
		StateZ_JsonWriter.Write(GetName(), json);
	}
}
#endif
