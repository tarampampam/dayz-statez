/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Collector for EffectArea stats.
*/
class MetricZ_CollectorAreas : MetricZ_CollectorBase
{
	override string GetName()
	{
		return "areas";
	}

	override bool IsEnabled()
	{
		return (MetricZ_Config.IsLoaded() && !MetricZ_Config.Get().disabled_metrics.areas);
	}

	override void Flush()
	{
		string json;
		MetricZ_EntitiesWriter.BuildAreasJson(json);
		MetricZ_JsonWriter.Write(MetricZ_Config.Get().json.export_dir, GetName(), json);
	}
}
#endif
