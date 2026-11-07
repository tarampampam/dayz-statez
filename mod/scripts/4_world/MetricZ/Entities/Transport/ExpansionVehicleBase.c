/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
#ifdef EXPANSIONMODVEHICLE
/**
    \brief Updates MetricZ counters on spawn/cleanup and kills expansion vehicles.
*/
modded class ExpansionVehicleBase
{
	protected bool m_MetricZ_IsKilled; //!< Prevent counting kills on destroyed vehicles.
	protected ref MetricZ_TransportMetrics m_MetricZ; //!< Metrics for this expansion vehicle.

	/**
	    \brief Initialize transport metrics for expansion vehicle loaded from persistence.
	*/
	override void EEOnAfterLoad()
	{
		super.EEOnAfterLoad();

		if (!MetricZ_Config.IsLoaded() || MetricZ_Config.Get().disabled_metrics.transports)
			return;

		if (!m_MetricZ)
			m_MetricZ = new MetricZ_TransportMetrics();

		// Init metrics for a persistent transport loaded from save with the actual persistence hash.
		m_MetricZ.Init(this);

		// prevent double counting for save/load destroyed vehicle
		if (IsDamageDestroyed())
			m_MetricZ_IsKilled = true;
	}


	/**
	    \brief Increment transport based counters.
	    \details No-op if transport metrics disabled.
	*/
	override void EEInit()
	{
		super.EEInit();

		if (!MetricZ_Config.IsLoaded() || MetricZ_Config.Get().disabled_metrics.transports)
			return;

		if (Expansion_IsBoat())
			MetricZ_Storage.s_Boats += 1;
		else if (Expansion_IsHelicopter() || Expansion_IsPlane())
			MetricZ_Storage.s_Helicopters += 1;
		else
			MetricZ_Storage.s_Cars += 1;

		MetricZ_TransportRegistry.Register(this);

		if (!m_MetricZ)
			m_MetricZ = new MetricZ_TransportMetrics();

		// Scheduled init of metrics for the created transport.
		// In this state, the persistence hash is not guaranteed and must be loaded later.
		// However, if the transport was not loaded but created via debug, this is the only reliable place for integration.
		m_MetricZ.InitLater(this);

		// prevent double counting for save/load destroyed vehicle
		if (IsDamageDestroyed())
			m_MetricZ_IsKilled = true;
	}

	/**
	    \brief Decrement transport based counters.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (MetricZ_Config.IsLoaded() && !MetricZ_Config.Get().disabled_metrics.transports) {
			m_MetricZ = null;

			if (Expansion_IsBoat())
				MetricZ_Storage.s_Boats -= 1;
			else if (Expansion_IsHelicopter() || Expansion_IsPlane())
				MetricZ_Storage.s_Helicopters -= 1;
			else
				MetricZ_Storage.s_Cars -= 1;

			MetricZ_TransportRegistry.Unregister(this);
		}

		super.EEDelete(parent);
	}

	/**
	    \brief Increment transport based destroyed counter on kill if enabled.
	*/
	override void EEKilled(Object killer)
	{
		if (MetricZ_Config.IsLoaded() && !MetricZ_Config.Get().disabled_metrics.transports && !m_MetricZ_IsKilled) {
			if (Expansion_IsBoat())
				MetricZ_Storage.s_BoatsDestroys += 1;
			else if (Expansion_IsHelicopter() || Expansion_IsPlane())
				MetricZ_Storage.s_HelicoptersDestroys += 1;
			else
				MetricZ_Storage.s_CarsDestroys += 1;

			m_MetricZ_IsKilled = true;
		}

		super.EEKilled(killer);
	}

	/**
	    \brief Accessor for per-expansion vehicle metrics.
	    \return \p MetricZ_TransportMetrics or null.
	*/
	MetricZ_TransportMetrics MetricZ_GetMetrics()
	{
		return m_MetricZ;
	}
}
#endif
#endif
