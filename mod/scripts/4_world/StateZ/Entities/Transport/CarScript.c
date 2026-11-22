/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Thin StateZ lifecycle shim for cars. Logic lives in StateZ_TransportHooks.
*/
modded class CarScript
{
	protected bool m_StateZ_IsKilled; //!< Prevent counting kills on destroyed vehicles.
	protected ref StateZ_TransportMetrics m_StateZ; //!< Metrics for this vehicle.

	override void EEOnAfterLoad()
	{
		super.EEOnAfterLoad();

		if (!StateZ_Config.IsLoaded())
			return;

		m_StateZ = StateZ_TransportHooks.OnAfterLoad(this, m_StateZ);

		// prevent double counting for save/load destroyed vehicle
		if (IsDamageDestroyed())
			m_StateZ_IsKilled = true;
	}

	override void EEInit()
	{
		super.EEInit();

		if (!StateZ_Config.IsLoaded())
			return;

		m_StateZ = StateZ_TransportHooks.OnInit(this);

		// prevent double counting for save/load destroyed vehicle
		if (IsDamageDestroyed())
			m_StateZ_IsKilled = true;
	}

	override void EEDelete(EntityAI parent)
	{
		if (StateZ_Config.IsLoaded()) {
			StateZ_TransportHooks.OnDelete(this, m_StateZ);
			m_StateZ = null;
		}

		super.EEDelete(parent);
	}

	override void EEKilled(Object killer)
	{
		if (StateZ_Config.IsLoaded() && !m_StateZ_IsKilled) {
			StateZ_TransportHooks.OnKilled(this);
			m_StateZ_IsKilled = true;
		}

		super.EEKilled(killer);
	}
}
#endif
