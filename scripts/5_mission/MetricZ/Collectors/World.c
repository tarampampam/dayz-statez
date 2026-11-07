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
class MetricZ_CollectorWorld : MetricZ_CollectorBase
{
	override string GetName()
	{
		return "world";
	}

	override bool IsEnabled()
	{
		return MetricZ_Config.IsLoaded();
	}

	override void Flush()
	{
		if (!MetricZ_Config.IsLoaded())
			return;

		MetricZ_Storage.Update();

		MetricZ_WorldSnapshot snap = new MetricZ_WorldSnapshot();
		MetricZ_Storage.FillSnapshot(snap);

		string json;
		JsonSerializer().WriteToString(snap, false, json);
		MetricZ_JsonWriter.Write(MetricZ_Config.Get().json.export_dir, GetName(), json);
	}
}
#endif
