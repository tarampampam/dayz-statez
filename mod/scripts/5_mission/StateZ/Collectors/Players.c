/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Collector for Player stats.
*/
class StateZ_CollectorPlayers : StateZ_CollectorBase
{
	override string GetName()
	{
		return "players";
	}

	override bool IsEnabled()
	{
		return (StateZ_Config.IsLoaded() && !StateZ_Config.Get().disabled_metrics.players);
	}

	override void Flush()
	{
		string json;
		StateZ_EntitiesWriter.BuildPlayersJson(json);
		StateZ_JsonWriter.Write(StateZ_Config.Get().json.export_dir, GetName(), json);
	}
}
#endif
