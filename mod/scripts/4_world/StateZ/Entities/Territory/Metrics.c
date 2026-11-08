/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Per-TerritoryFlag data collector.
    \details Stores plain fields updated during collection. Fills snapshot on demand.
*/
class StateZ_TerritoryMetrics
{
	protected TerritoryFlag m_TerritoryFlag;

	float m_LifetimeFraction;
	float m_PosX;
	float m_PosZ;
	float m_RefresherRadius;

	/**
	    \brief One-time init with territory reference.
	    \param territory TerritoryFlag instance to initialize metrics for.
	*/
	void Init(TerritoryFlag territory)
	{
		if (!territory || m_TerritoryFlag)
			return;

		m_TerritoryFlag = territory;
		m_RefresherRadius = GameConstants.REFRESHER_RADIUS;

		vector pos = territory.GetPosition();
		m_PosX = pos[0];
		m_PosZ = pos[2];
	}

	/**
	    \brief Update all fields from the territory state.
	*/
	void Update()
	{
		if (!m_TerritoryFlag)
			return;

		m_LifetimeFraction = m_TerritoryFlag.GetRefresherTime01();
	}

	/**
	    \brief Fill a territory snapshot with current field values.
	    \param snap Snapshot instance to populate.
	*/
	void ToSnapshot(StateZ_TerritorySnapshot snap)
	{
		if (!snap)
			return;

		snap.pos_x = StateZ_Num.Round2(m_PosX);
		snap.pos_y = StateZ_Num.Round2(m_PosZ);
		snap.lifetime_fraction = StateZ_Num.Round4(m_LifetimeFraction);
		snap.refresher_radius = StateZ_Num.Round2(m_RefresherRadius);
	}
}
#endif
