/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Mission hooks for StateZ lifecycle and periodic scraping.
*/
modded class MissionServer : MissionBase
{
	protected ref StateZ_Scheduler m_StateZ;

	/**
	    \brief Initialize StateZ and schedule first collection.
	*/
	override void OnInit()
	{
		m_StateZ = StateZ_Scheduler.Get();

		super.OnInit();
	}

	/**
	    \brief Called on mission finish to shutdown StateZ.
	*/
	override void OnMissionFinish()
	{
		if (m_StateZ)
			m_StateZ.Shutdown();

		super.OnMissionFinish();
	}

	/**
	    \brief Count new player spawns.
	    \return \p PlayerBase Created player
	*/
	override PlayerBase OnClientNewEvent(PlayerIdentity identity, vector pos, ParamsReadContext ctx)
	{
		super.OnClientNewEvent(identity, pos, ctx);

		if (m_player)
			StateZ_Storage.s_PlayersSpawns += 1;

		return m_player;
	}

	/**
	    \brief Track number of corpses after server update.
	    \details Calls base, then mirrors current count into StateZ_Storage.s_Corpses.
	*/
	override void UpdateCorpseStatesServer()
	{
		super.UpdateCorpseStatesServer();

		if (StateZ_Storage.IsInitialized())
			StateZ_Storage.s_Corpses = m_DeadPlayersArray.Count();
	}

	/**
	    \brief Count player chat messages.
	*/
	override void OnEvent(EventType eventTypeId, Param params)
	{
		if (eventTypeId == ChatMessageEventTypeID)
			StateZ_Storage.s_ChatMessages += 1;

		super.OnEvent(eventTypeId, params);
	}
}
#endif
