/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Collector for Transport stats.
*/
class StateZ_CollectorTransports : StateZ_CollectorBase
{
	override string GetName()
	{
		return "vehicles";
	}

	override void Flush()
	{
		string json;
		StateZ_EntitiesWriter.BuildTransportJson(json);
		StateZ_JsonWriter.Write(GetName(), json);
	}
}
#endif
