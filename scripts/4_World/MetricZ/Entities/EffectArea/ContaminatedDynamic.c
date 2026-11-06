/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Updates MetricZ artillery strike counter.
*/
modded class ContaminatedArea_Dynamic
{
	override void EEOnCECreate()
	{
		super.EEOnCECreate();

		if (m_DecayState == eAreaDecayStage.INIT)
			MetricZ_Storage.s_Artillery += 1;
	}
}
#endif
