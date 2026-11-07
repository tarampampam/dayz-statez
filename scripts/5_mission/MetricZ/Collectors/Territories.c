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
class MetricZ_CollectorTerritories : MetricZ_CollectorBase
{
	override string GetName()
	{
		return "territories";
	}

	override bool IsEnabled()
	{
		return (MetricZ_Config.IsLoaded() && !MetricZ_Config.Get().disabled_metrics.territories);
	}

	override void Flush()
	{
		string json;
		MetricZ_EntitiesWriter.BuildTerritoryJson(json);
		MetricZ_JsonWriter.Write(MetricZ_Config.Get().json.export_dir, GetName(), json);
	}
}
#endif
