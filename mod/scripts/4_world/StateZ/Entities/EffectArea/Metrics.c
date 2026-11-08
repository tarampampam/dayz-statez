/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Per-EffectArea data collector.
    \details Stores plain fields updated during collection. Fills snapshot on demand.
*/
class StateZ_EffectAreaMetrics
{
	protected EffectArea m_Area;

	string m_Type;
	string m_ClassName;
	float m_PosX;
	float m_PosZ;
	float m_Radius;
	int m_Insiders;

	/**
	    \brief One-time init with area reference.
	    \param area EffectArea instance to initialize metrics for.
	*/
	void Init(EffectArea area)
	{
		if (!area || m_Area)
			return;

		m_Area = area;
		m_Type = area.StateZ_GetType();
		m_ClassName = area.ClassName();
		m_Radius = area.StateZ_GetRadius();

		vector pos = area.m_Position;
		if (pos == vector.Zero)
			pos = area.GetPosition();

		m_PosX = pos[0];
		m_PosZ = pos[2];
	}

	/**
	    \brief Update all fields from the EffectArea state.
	*/
	void Update()
	{
		if (!m_Area)
			return;

		m_Insiders = m_Area.StateZ_GetInsidersCount();
	}

	/**
	    \brief Fill an area snapshot with current field values.
	    \param snap Snapshot instance to populate.
	*/
	void ToSnapshot(StateZ_AreaSnapshot snap)
	{
		if (!snap)
			return;

		snap.type = m_Type;
		snap.class_name = m_ClassName;
		snap.pos_x = StateZ_Num.Round2(m_PosX);
		snap.pos_y = StateZ_Num.Round2(m_PosZ);
		snap.radius = StateZ_Num.Round2(m_Radius);
		snap.insiders = m_Insiders;
	}
}
#endif
