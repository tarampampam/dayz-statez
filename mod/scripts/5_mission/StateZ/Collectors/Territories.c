/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Collector for Territory stats.
*/
class StateZ_CollectorTerritories : StateZ_CollectorBase
{
	override string GetName()
	{
		return "territories";
	}

	override void Flush()
	{
		string json;
		StateZ_EntitiesWriter.BuildTerritoryJson(json);
		StateZ_JsonWriter.Write(StateZ_Config.Get().json.export_dir, GetName(), json);
	}
}
#endif
