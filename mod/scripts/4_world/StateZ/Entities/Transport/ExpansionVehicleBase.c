/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
#ifdef EXPANSIONMODVEHICLE
/**
    \brief Updates StateZ counters on spawn/cleanup and kills expansion vehicles.
*/
modded class ExpansionVehicleBase
{
	protected bool m_StateZ_IsKilled; //!< Prevent counting kills on destroyed vehicles.
	protected ref StateZ_TransportMetrics m_StateZ; //!< Metrics for this expansion vehicle.

	/**
	    \brief Initialize transport metrics for expansion vehicle loaded from persistence.
	*/
	override void EEOnAfterLoad()
	{
		super.EEOnAfterLoad();

		if (!StateZ_Config.IsLoaded())
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
	    \brief Increment transport based counters.
	    \details No-op if transport metrics disabled.
	*/
	override void EEInit()
	{
		super.EEInit();

		if (!StateZ_Config.IsLoaded())
			return;

		if (Expansion_IsBoat())
			StateZ_Storage.s_Boats += 1;
		else if (Expansion_IsHelicopter() || Expansion_IsPlane())
			StateZ_Storage.s_Helicopters += 1;
		else
			StateZ_Storage.s_Cars += 1;

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
	    \brief Decrement transport based counters.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (StateZ_Config.IsLoaded()) {
			m_StateZ = null;

			if (Expansion_IsBoat())
				StateZ_Storage.s_Boats -= 1;
			else if (Expansion_IsHelicopter() || Expansion_IsPlane())
				StateZ_Storage.s_Helicopters -= 1;
			else
				StateZ_Storage.s_Cars -= 1;

			StateZ_TransportRegistry.Unregister(this);
		}

		super.EEDelete(parent);
	}

	/**
	    \brief Increment transport based destroyed counter on kill if enabled.
	*/
	override void EEKilled(Object killer)
	{
		if (StateZ_Config.IsLoaded() && !m_StateZ_IsKilled) {
			if (Expansion_IsBoat())
				StateZ_Storage.s_BoatsDestroys += 1;
			else if (Expansion_IsHelicopter() || Expansion_IsPlane())
				StateZ_Storage.s_HelicoptersDestroys += 1;
			else
				StateZ_Storage.s_CarsDestroys += 1;

			m_StateZ_IsKilled = true;
		}

		super.EEKilled(killer);
	}

	/**
	    \brief Accessor for per-expansion vehicle metrics.
	    \return \p StateZ_TransportMetrics or null.
	*/
	StateZ_TransportMetrics StateZ_GetMetrics()
	{
		return m_StateZ;
	}
}
#endif
#endif
