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
modded class CarScript
{
	protected bool m_StateZ_IsKilled; //!< Prevent counting kills on destroyed vehicles.
	protected ref StateZ_TransportMetrics m_StateZ; //!< Metrics for this car.

	/**
	    \brief Initialize transport metrics for car loaded from persistence.
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
	    \brief Register car in transport registry and create metrics.
	    \details No-op if transport metrics disabled.
	*/
	override void EEInit()
	{
		super.EEInit();

		if (!StateZ_Config.IsLoaded() || StateZ_Config.Get().disabled_metrics.transports)
			return;

#ifdef EXPANSIONMODVEHICLE
		if (IsInherited(ExpansionBoatScript))
			StateZ_Storage.s_Boats += 1;
		else if (IsInherited(ExpansionHelicopterScript))
			StateZ_Storage.s_Helicopters += 1;
		else
			StateZ_Storage.s_Cars += 1;
#else
		StateZ_Storage.s_Cars += 1;
#endif

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
	    \brief Unregister car and decrement counts on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (StateZ_Config.IsLoaded() && !StateZ_Config.Get().disabled_metrics.transports) {
			m_StateZ = null;

#ifdef EXPANSIONMODVEHICLE
			if (IsInherited(ExpansionBoatScript))
				StateZ_Storage.s_Boats -= 1;
			else if (IsInherited(ExpansionHelicopterScript))
				StateZ_Storage.s_Helicopters -= 1;
			else
				StateZ_Storage.s_Cars -= 1;
#else
			StateZ_Storage.s_Cars -= 1;
#endif

			StateZ_TransportRegistry.Unregister(this);
		}

		super.EEDelete(parent);
	}

	/**
	    \brief Increment cars destroyed counter on kill if enabled.
	*/
	override void EEKilled(Object killer)
	{
		if (StateZ_Config.IsLoaded() && !StateZ_Config.Get().disabled_metrics.transports && !m_StateZ_IsKilled) {
#ifdef EXPANSIONMODVEHICLE
			if (IsInherited(ExpansionBoatScript))
				StateZ_Storage.s_BoatsDestroys += 1;
			else if (IsInherited(ExpansionHelicopterScript))
				StateZ_Storage.s_HelicoptersDestroys += 1;
			else
				StateZ_Storage.s_CarsDestroys += 1;
#else
			StateZ_Storage.s_CarsDestroys += 1;
#endif
			m_StateZ_IsKilled = true;
		}

		super.EEKilled(killer);
	}

	/**
	    \brief Accessor for per-car metrics.
	    \return \p StateZ_TransportMetrics or null.
	*/
	StateZ_TransportMetrics StateZ_GetMetrics()
	{
		return m_StateZ;
	}
}
#endif
