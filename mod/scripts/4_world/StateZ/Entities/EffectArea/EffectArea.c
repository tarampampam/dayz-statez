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
modded class EffectArea
{
	protected bool m_StateZ_AreaInit; //!< True if the area has been initialized.
	protected int m_StateZ_InsidersCount; //!< Count of players currently inside the zone.
	protected ref StateZ_EffectAreaMetrics m_StateZ; //!< Metrics for this area.

	/**
	    \brief Standard entity initialization.
	    \details
	        StateZ initialization is skipped here. For dynamic zones, EEInit triggers
	        while the artillery projectile is still airborne. We must wait for the actual zone creation.
	*/
	override void EEInit()
	{
		super.EEInit();

		// StateZ_Init();
	}

	/**
	    \brief Finalizes zone setup and registers metrics collector.
	*/
	override void InitZone()
	{
		super.InitZone();

		StateZ_Init();
	}

	/**
	    \brief Unregister metrics and decrement count on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		StateZ_Delete();

		super.EEDelete(parent);
	}

	/**
	    \brief Increase insiders players count.
	*/
	override void OnPlayerEnterServer(PlayerBase player, EffectTrigger trigger)
	{
		super.OnPlayerEnterServer(player, trigger);

		if (player)
			m_StateZ_InsidersCount++;
	}

	/**
	    \brief Decrease insiders players count.
	*/
	override void OnPlayerExitServer(PlayerBase player, EffectTrigger trigger)
	{
		super.OnPlayerExitServer(player, trigger);

		if (player)
			m_StateZ_InsidersCount--;
	}

	/**
	    \brief Get string representation of zone type.
	    \return \p string Enum eZoneType as string.
	*/
	string StateZ_GetType()
	{
		return EnumTools.EnumToString(eZoneType, m_Type);
	}

	/**
	    \brief Get area radius.
	    \return \p int Radius.
	*/
	float StateZ_GetRadius()
	{
		return m_Radius;
	}

	/**
	    \brief Get count of players currently inside the zone.
	    \return \p int
	*/
	int StateZ_GetInsidersCount()
	{
		return m_StateZ_InsidersCount;
	}

	/**
	    \brief Accessor for per-area metrics.
	    \return \p StateZ_EffectAreaMetrics or null.
	*/
	StateZ_EffectAreaMetrics StateZ_GetMetrics()
	{
		return m_StateZ;
	}

	/**
	    \brief Initialize metrics tracking for this zone.
	    \details Updates global count, registers in StateZ_EffectAreaRegistry, creates metric instance.
	*/
	protected void StateZ_Init()
	{
		if (!StateZ_Config.IsLoaded() || m_StateZ_AreaInit)
			return;

		m_StateZ_AreaInit = true;
		StateZ_Storage.s_EffectAreas += 1;

		StateZ_EffectAreaRegistry.Register(this);

		if (!m_StateZ) {
			m_StateZ = new StateZ_EffectAreaMetrics();
			m_StateZ.Init(this);
		}
	}

	/**
	    \brief Cleanup metrics on deletion.
	    \details Updates global count, unregisters from registry, destroys metric instance.
	*/
	protected void StateZ_Delete()
	{
		if (!StateZ_Config.IsLoaded() || !m_StateZ_AreaInit)
			return;

		m_StateZ = null;
		StateZ_EffectAreaRegistry.Unregister(this);

		StateZ_Storage.s_EffectAreas -= 1;
	}
}
#endif
