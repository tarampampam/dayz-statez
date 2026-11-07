/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Update StateZ counters on spawn/cleanup and kills, and register a per-player metrics instance.
*/
modded class PlayerBase
{
	protected bool m_StateZ_IsKilled; //!< Prevent counting kills on already killed body.
	protected bool m_StateZ_IsLastHit; //!< Prevent counting hits on death bodies except last hit.

	protected ref StateZ_PlayerMetrics m_StateZ;

	/**
	    \brief Register per-player metrics collector on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		if (!StateZ_Config.IsLoaded())
			return;

#ifdef EXPANSIONMODAI
		if (IsInherited(eAIBase)) {
			if (IsInherited(eAINPCBase))
				StateZ_Storage.s_ExpansionAINPC += 1;
			else
				StateZ_Storage.s_ExpansionAI += 1;

			return;
		}
#endif

		if (!m_StateZ) {
			m_StateZ = new StateZ_PlayerMetrics();
			m_StateZ.Init(this);
		}
	}

	/**
	    \brief Drop metrics reference on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (StateZ_Config.IsLoaded()) {
#ifdef EXPANSIONMODAI
			if (IsInherited(eAIBase)) {
				if (IsInherited(eAINPCBase))
					StateZ_Storage.s_ExpansionAINPC -= 1;
				else
					StateZ_Storage.s_ExpansionAI -= 1;
			}
#endif

			m_StateZ = null;
		}

		super.EEDelete(parent);
	}

	/**
	    \brief Increment players_deaths counter on kill.
	*/
	override void EEKilled(Object killer)
	{
		if (StateZ_Config.IsLoaded() && !m_StateZ_IsKilled) {
			m_StateZ_IsKilled = true;

#ifdef EXPANSIONMODAI
			if (IsInherited(eAIBase)) {
				if (killer != this)
					StateZ_WeaponStats.OnCreatureKilled(killer);

				if (IsInherited(eAINPCBase))
					StateZ_Storage.s_ExpansionAINPCDeaths += 1;
				else
					StateZ_Storage.s_ExpansionAIDeaths += 1;

				super.EEKilled(killer);

				return;
			}
#endif

			StateZ_Storage.s_PlayersDeaths += 1;

			if (killer != this)
				StateZ_WeaponStats.OnPlayerKilled(killer);
		}

		super.EEKilled(killer);
	}

	/**
	    \brief Capture hit statistics for metric counters.
	*/
	override void EEHitBy(TotalDamageResult damageResult, int damageType, EntityAI source, int component, string dmgZone, string ammo, vector modelPos, float speedCoef)
	{
		super.EEHitBy(damageResult, damageType, source, component, dmgZone, ammo, modelPos, speedCoef);

		if (!StateZ_Config.IsLoaded())
			return;

		if (source != this && !m_StateZ_IsLastHit) {
			if (IsDamageDestroyed())
				m_StateZ_IsLastHit = true;

			if (damageResult) {
				float damage = damageResult.GetDamage(dmgZone, "");
				if (damage < StateZ_Config.Get().thresholds.hit_damage)
					return;

				if (source && source.IsTransport() && damage < StateZ_Config.Get().thresholds.hit_damage_vehicle)
					return;
			}

#ifdef EXPANSIONMODAI
			if (IsInherited(eAIBase)) {
				StateZ_HitStats.OnCreatureHit(ammo);
				return;
			}
#endif

			StateZ_HitStats.OnPlayerHit(ammo);
		}
	}

	/**
	    \brief Accessor for per-player metrics.
	    \return \p StateZ_PlayerMetrics or null if not initialized.
	*/
	StateZ_PlayerMetrics StateZ_GetMetrics()
	{
		return m_StateZ;
	}

	/**
	    \brief Update player network stats.
	*/
	override void EOnPostFrame(IEntity other, int extra)
	{
		super.EOnPostFrame(other, extra);

		if (!StateZ_Config.IsLoaded())
			return;

#ifdef EXPANSIONMODAI
		if (IsInherited(eAIBase))
			return;
#endif

		if (m_StateZ)
			m_StateZ.SampleNetwork();
	}
}
#endif
