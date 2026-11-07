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
class MetricZ_CollectorZombies : MetricZ_CollectorBase
{
	override string GetName()
	{
		return "zombies";
	}

	override bool IsEnabled()
	{
		return (MetricZ_Config.IsLoaded() && !MetricZ_Config.Get().disabled_metrics.zombies);
	}

	override void Flush()
	{
		string json;
		MetricZ_ZombieStats.BuildJson(json);
		MetricZ_JsonWriter.Write(MetricZ_Config.Get().json.export_dir, GetName(), json);
	}
}
#endif
