/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Global registry of active vehicle metrics.
    \details Holds one StateZ_TransportMetrics per live vehicle, regardless of
             family. Storing the metrics (not the Transport) lets the writer emit
             every vehicle without per-family casts. Updated on EEInit/EEDelete.
*/
class StateZ_TransportRegistry
{
	protected static ref array<ref StateZ_TransportMetrics> s_Registry = new array<ref StateZ_TransportMetrics>();

	/**
	    \brief Register vehicle metrics.
	    \details Ignores duplicates and null references.
	*/
	static void Register(StateZ_TransportMetrics metrics)
	{
		if (!metrics)
			return;

		if (s_Registry.Find(metrics) == -1)
			s_Registry.Insert(metrics);
	}

	/**
	    \brief Remove vehicle metrics.
	    \details Does nothing if not found or null.
	*/
	static void Unregister(StateZ_TransportMetrics metrics)
	{
		if (!metrics)
			return;

		int i = s_Registry.Find(metrics);
		if (i >= 0)
			s_Registry.Remove(i);
	}

	/**
	    \brief Returns direct reference to the live registry.
	    \warning Do NOT modify this array (Remove/Insert) while iterating!
	*/
	static array<ref StateZ_TransportMetrics> GetList()
	{
		return s_Registry;
	}
}
#endif
