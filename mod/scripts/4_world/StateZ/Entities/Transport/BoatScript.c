/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Updates StateZ counters on spawn/cleanup and kills.
*/
modded class BoatScript
{
	protected bool m_StateZ_IsKilled; //!< Prevent counting kills on destroyed vehicles.
	protected ref StateZ_TransportMetrics m_StateZ; //!< Metrics for this boat.

	/**
	    \brief Initialize transport metrics for boat loaded from persistence.
	*/
	override void EEOnAfterLoad()
	{
		super.EEOnAfterLoad();

		if (!StateZ_Config.IsLoaded() || StateZ_Config.Get().disabled_metrics.transports)
			return;

		if (!m_StateZ)
			m_StateZ = new StateZ_TransportMetrics();

		// Init metrics for a persistent transport loaded from save with the actual persistence hash.
		m_StateZ.Init(this);

		// prevent double counting for save/load destroyed vehicle
		if (IsDamageDestroyed())
			m_StateZ_IsKilled = true;
	}

	/**
	    \brief Register boat in transport registry and create metrics.
	        \details No-op if transport metrics disabled.
	*/
	override void EEInit()
	{
		super.EEInit();

		if (!StateZ_Config.IsLoaded() || StateZ_Config.Get().disabled_metrics.transports)
			return;

		StateZ_Storage.s_Boats += 1;
		StateZ_TransportRegistry.Register(this);

		if (!m_StateZ)
			m_StateZ = new StateZ_TransportMetrics();

		// Scheduled init of metrics for the created transport.
		// In this state, the persistence hash is not guaranteed and must be loaded later.
		// However, if the transport was not loaded but created via debug, this is the only reliable place for integration.
		m_StateZ.InitLater(this);

		// prevent double counting for save/load destroyed vehicle
		if (IsDamageDestroyed())
			m_StateZ_IsKilled = true;
	}

	/**
	    \brief Unregister boat and decrement counts on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (StateZ_Config.IsLoaded() && !StateZ_Config.Get().disabled_metrics.transports) {
			m_StateZ = null;
			StateZ_Storage.s_Boats -= 1;
			StateZ_TransportRegistry.Unregister(this);
		}

		super.EEDelete(parent);
	}

	/**
	    \brief Increment boats destroyed counter on kill if enabled.
	*/
	override void EEKilled(Object killer)
	{
		if (StateZ_Config.IsLoaded() && !StateZ_Config.Get().disabled_metrics.transports && !m_StateZ_IsKilled) {
			StateZ_Storage.s_BoatsDestroys += 1;
			m_StateZ_IsKilled = true;
		}

		super.EEKilled(killer);
	}

	/**
	    \brief Accessor for per-boat metrics.
	    \return \p StateZ_TransportMetrics or null.
	*/
	StateZ_TransportMetrics StateZ_GetMetrics()
	{
		return m_StateZ;
	}
}
#endif
