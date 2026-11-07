/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Collector for Animal stats.
*/
class StateZ_CollectorAnimals : StateZ_CollectorBase
{
	override string GetName()
	{
		return "animals";
	}

	override void Flush()
	{
		string json;
		StateZ_AnimalStats.BuildJson(json);
		StateZ_JsonWriter.Write(GetName(), json);
	}
}
#endif
