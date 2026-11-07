/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Collector for World stats.
*/
class StateZ_CollectorWorld : StateZ_CollectorBase
{
	override string GetName()
	{
		return "world";
	}

	override void Flush()
	{
		if (!StateZ_Config.IsLoaded())
			return;

		StateZ_Storage.Update();

		StateZ_WorldSnapshot snap = new StateZ_WorldSnapshot();
		StateZ_Storage.FillSnapshot(snap);

		string json;
		JsonSerializer().WriteToString(snap, false, json);
		StateZ_JsonWriter.Write(StateZ_Config.Get().json.export_dir, GetName(), json);
	}
}
#endif
