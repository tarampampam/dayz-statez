/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER

class StateZ_EventStatsDTO
{
	void StateZ_EventStatsDTO()
	{
		events = new map<string, int>();
	}

	ref map<string, int> events;
}

/**
    \brief Event counters aggregator
*/
class StateZ_EventStats
{
	protected static ref map<EventType, string> s_EventNames; //!< Map of EventType to name
	protected static ref map<EventType, int> s_EventsRegistry = new map<EventType, int>(); //!< Map of EventType to count

	/**
	    \brief Increment counter for an EventType.
	    \param eventTypeId Engine event id.
	*/
	static void Inc(EventType eventTypeId)
	{
		int v;
		if (s_EventsRegistry.Find(eventTypeId, v))
			s_EventsRegistry.Set(eventTypeId, v + 1);
		else
			s_EventsRegistry.Insert(eventTypeId, 1);
	}

	/**
	    \brief Serialize all event data to JSON.
	    \param json Output JSON string.
	*/
	static void BuildJson(out string json)
	{
		if (s_EventsRegistry.Count() == 0) {
			json = "{}";
			return;
		}

		MakeNameMap();

		StateZ_EventStatsDTO dto = new StateZ_EventStatsDTO();

		for (int i = 0; i < s_EventsRegistry.Count(); ++i) {
			EventType id = s_EventsRegistry.GetKey(i);
			int val = s_EventsRegistry.GetElement(i);
			dto.events.Insert(EventName(id), val);
		}

		JsonSerializer().WriteToString(dto, false, json);
	}

	/**
	    \brief Build the EventType->name map once.
	    \details No-op if the map already exists.
	*/
	protected static void MakeNameMap()
	{
		if (s_EventNames)
			return;

		s_EventNames = new map<EventType, string>();

		// system
		s_EventNames.Insert(StartupEventTypeID, "Startup");
		s_EventNames.Insert(WorldCleaupEventTypeID, "WorldCleanUP");
		s_EventNames.Insert(SelectedUserChangedEventTypeID, "SelectedUserChanged");
		s_EventNames.Insert(ScriptLogEventTypeID, "ScriptLog");

		// MP session
		s_EventNames.Insert(MPSessionStartEventTypeID, "MPSessionStart");
		s_EventNames.Insert(MPSessionEndEventTypeID, "MPSessionEnd");
		s_EventNames.Insert(MPSessionFailEventTypeID, "MPSessionFail");
		s_EventNames.Insert(MPSessionPlayerReadyEventTypeID, "MPSessionPlayerReady");
		s_EventNames.Insert(MPConnectionLostEventTypeID, "MPConnectionLost");
		s_EventNames.Insert(MPConnectionCloseEventTypeID, "MPConnectionClose");

		// misc
		s_EventNames.Insert(ProgressEventTypeID, "Progress");
		s_EventNames.Insert(NetworkManagerClientEventTypeID, "NetworkManagerClient");
		s_EventNames.Insert(NetworkManagerServerEventTypeID, "NetworkManagerServer");
		s_EventNames.Insert(DialogQueuedEventTypeID, "DialogQueued");

		// chat
		s_EventNames.Insert(ChatMessageEventTypeID, "ChatMessage");
		s_EventNames.Insert(ChatChannelEventTypeID, "ChatChannel");

		// client lifecycle
		s_EventNames.Insert(ClientConnectedEventTypeID, "ClientConnected");
		s_EventNames.Insert(ClientPrepareEventTypeID, "ClientPrepare");
		s_EventNames.Insert(ClientNewEventTypeID, "ClientNew");
		s_EventNames.Insert(ClientNewReadyEventTypeID, "ClientNewReady");
		s_EventNames.Insert(ClientRespawnEventTypeID, "ClientReSpawn");
		s_EventNames.Insert(ClientReconnectEventTypeID, "ClientReconnect");
		s_EventNames.Insert(ClientReadyEventTypeID, "ClientReady");
		s_EventNames.Insert(ClientDisconnectedEventTypeID, "ClientDisconnected");
		s_EventNames.Insert(ClientRemovedEventTypeID, "ClientRemoved");

		// connectivity and server stats
		s_EventNames.Insert(ConnectivityStatsUpdatedEventTypeID, "ConnectivityStatsUpdated");
		s_EventNames.Insert(ServerFpsStatsUpdatedEventTypeID, "ServerFpsStatsUpdated");

		// login/logout flow
		s_EventNames.Insert(LogoutCancelEventTypeID, "LogoutCancel");
		s_EventNames.Insert(LoginTimeEventTypeID, "LoginTime");
		s_EventNames.Insert(RespawnEventTypeID, "Respawn");
		s_EventNames.Insert(PreloadEventTypeID, "Preload");
		s_EventNames.Insert(LogoutEventTypeID, "Logout");
		s_EventNames.Insert(LoginStatusEventTypeID, "LoginStatus");

		// VON
		s_EventNames.Insert(VONStateEventTypeID, "VONState");
		s_EventNames.Insert(VONStartSpeakingEventTypeID, "VONStartSpeaking");
		s_EventNames.Insert(VONStopSpeakingEventTypeID, "VONStopSpeaking");
		s_EventNames.Insert(VONUserStartedTransmittingAudioEventTypeID, "VONStartedTransmitting");
		s_EventNames.Insert(VONUserStoppedTransmittingAudioEventTypeID, "VONStoppedTransmitting");

		// other
		s_EventNames.Insert(PartyChatStatusChangedEventTypeID, "PartyChatStatusChanged");
		s_EventNames.Insert(DLCOwnerShipFailedEventTypeID, "DLCOwnerShipFailed");
		s_EventNames.Insert(SetFreeCameraEventTypeID, "SetFreeCamera");
		s_EventNames.Insert(ConnectingStartEventTypeID, "ConnectingStart");
		s_EventNames.Insert(ConnectingAbortEventTypeID, "ConnectingAbort");
		s_EventNames.Insert(PlayerDeathEventTypeID, "PlayerDeath");
		s_EventNames.Insert(NetworkInputBufferEventTypeID, "NetworkInputBuffer");
	}

	/**
	    \brief Resolve a human-readable name for an EventType.
	    \param id EventType id.
	    \return string Name or "unknown_<id>".
	*/
	protected static string EventName(EventType id)
	{
		if (!s_EventNames)
			return string.Format("unknown_%1", id);

		string name;
		if (s_EventNames.Find(id, name))
			return name;

		return string.Format("unknown_%1", id);
	}
}
#endif
