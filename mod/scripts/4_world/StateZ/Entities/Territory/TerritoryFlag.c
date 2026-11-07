/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Updates StateZ counters on spawn/cleanup.
*/
modded class TerritoryFlag
{
	protected ref StateZ_TerritoryMetrics m_StateZ;

	/**
	    \brief Track territory flag and register metrics collector.
	    \details Increments count and registers in registry.
	*/
	override void EEInit()
	{
		super.EEInit();

		if (!StateZ_Config.IsLoaded())
			return;

		StateZ_Storage.s_TerritoryFlags += 1;

		StateZ_TerritoryRegistry.Register(this);
		if (!m_StateZ) {
			m_StateZ = new StateZ_TerritoryMetrics();
			m_StateZ.Init(this);
		}
	}

	/**
	    \brief Unregister metrics and decrement count on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (StateZ_Config.IsLoaded()) {
			m_StateZ = null;
			StateZ_TerritoryRegistry.Unregister(this);
		}

		StateZ_Storage.s_TerritoryFlags -= 1;

		super.EEDelete(parent);
	}

	/**
	    \brief Accessor for per-territory metrics.
	    \return \p StateZ_TerritoryMetrics or null.
	*/
	StateZ_TerritoryMetrics StateZ_GetMetrics()
	{
		return m_StateZ;
	}
}
#endif
