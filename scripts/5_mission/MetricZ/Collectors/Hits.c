/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Collector for Hit stats.
*/
class MetricZ_CollectorHits : MetricZ_CollectorBase
{
	override string GetName()
	{
		return "hits";
	}

	override bool IsEnabled()
	{
		return (MetricZ_Config.IsLoaded() && !MetricZ_Config.Get().disabled_metrics.hits);
	}

	override void Flush()
	{
		string json;
		MetricZ_HitStats.BuildJson(json);
		MetricZ_JsonWriter.Write(MetricZ_Config.Get().json.export_dir, GetName(), json);
	}
}
#endif
