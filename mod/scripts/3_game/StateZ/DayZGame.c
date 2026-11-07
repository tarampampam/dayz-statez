/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    Server-side instrumentation for Event metrics.
    Counts all events by EventType.
*/
modded class DayZGame
{
	// number of calls to UpdatePathgraphRegionByObject() on server
	static int s_StateZ_PathGraphUpdates;

	/**
	    \brief Pump frame monitor each frame.
	    \param doSim True if the simulation is running, false otherwise
	    \param timeslice Delta time of the last frame in seconds.
	*/
	override void OnPostUpdate(bool doSim, float timeslice)
	{
		if (doSim)
			StateZ_FrameMonitor.OnUpdate(timeslice);

		super.OnPostUpdate(doSim, timeslice);
	}

	/**
	    \brief Hook engine events to count occurrences.
	    \details Increments per-EventType counter, then calls base.
	*/
	override void OnEvent(EventType eventTypeId, Param params)
	{
		// count all events on server
		if (StateZ_Config.IsLoaded())
			StateZ_EventStats.Inc(eventTypeId);

		super.OnEvent(eventTypeId, params);
	}

	/**
	    \brief Hook CGame event for calculate path graph updates by placing objects.
	    \details Increments static counter for reuse in StateZ_Storage.
	*/
	override void UpdatePathgraphRegionByObject(Object object)
	{
		if (object)
			s_StateZ_PathGraphUpdates++;

		super.UpdatePathgraphRegionByObject(object);
	}
}
#endif
