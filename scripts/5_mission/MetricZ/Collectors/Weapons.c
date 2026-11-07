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
class MetricZ_CollectorWeapons : MetricZ_CollectorBase
{
	override string GetName()
	{
		return "weapons";
	}

	override bool IsEnabled()
	{
		return (MetricZ_Config.IsLoaded() && !MetricZ_Config.Get().disabled_metrics.weapons);
	}

	override void Flush()
	{
		string json;
		MetricZ_WeaponStats.BuildJson(json);
		MetricZ_JsonWriter.Write(MetricZ_Config.Get().json.export_dir, GetName(), json);
	}
}
#endif
