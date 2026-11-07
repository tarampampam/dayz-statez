/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER

class MetricZ_ZombieStatsDTO
{
	void MetricZ_ZombieStatsDTO()
	{
		by_type = new map<string, int>();
		by_state = new map<string, int>();
	}

	ref map<string, int> by_type;
	ref map<string, int> by_state;
}

/**
    \brief Infected live gauges.
    \details
      - tracks per-mind-state counts
      - tracks infected count by canonical zombie type
*/
class MetricZ_ZombieStats
{
	static const int MINDSTATE_DEAD = 1000; //!< Synthetic mind-state id for dead infected
	protected static ref map<int, int> s_CountByState = new map<int, int>(); //!< State -> current infected count.
	protected static ref map<int, string> s_MindStates; //!< State -> human-readable name. Built lazily.
	protected static ref map<string, int> s_TypeStorage = new map<string, int>(); //!< Type -> current infected count.

	/**
	    \brief Fill name map once.
	    \details No-op on subsequent calls.
	*/
	protected static void EnsureNames()
	{
		if (s_MindStates)
			return;

		s_MindStates = new map<int, string>();
		s_MindStates.Insert(DayZInfectedConstants.MINDSTATE_CALM, "calm");
		s_MindStates.Insert(DayZInfectedConstants.MINDSTATE_DISTURBED, "disturbed");
		s_MindStates.Insert(DayZInfectedConstants.MINDSTATE_ALERTED, "alerted");
		s_MindStates.Insert(DayZInfectedConstants.MINDSTATE_CHASE, "chase");
		s_MindStates.Insert(DayZInfectedConstants.MINDSTATE_FIGHT, "fight");
		s_MindStates.Insert(MINDSTATE_DEAD, "dead");
	}

	/**
	    \brief Apply delta to a mind-state bucket.
	    \details Ignores negative state. Never goes below zero.
	    \param state Mind state id
	    \param delta +/- amount
	*/
	protected static void AddState(int state, int delta)
	{
		if (state < 0)
			return;

		int v;
		if (s_CountByState.Find(state, v))
			s_CountByState.Set(state, Math.Max(v + delta, 0));
		else if (delta > 0)
			s_CountByState.Insert(state, delta);
	}

	/**
	    \brief Apply delta to a zombie-type bucket.
	    \details Empty type is ignored. Never goes below zero, removes when reaches 0.
	    \param type Canonical zombie type
	    \param delta +/- amount
	*/
	protected static void AddType(string type, int delta)
	{
		type.TrimInPlace();
		if (type == string.Empty || delta == 0)
			return;

		int v;
		if (s_TypeStorage.Find(type, v)) {
			v = v + delta;
			if (v <= 0)
				s_TypeStorage.Remove(type);
			else
				s_TypeStorage.Set(type, v);
		} else if (delta > 0)
			s_TypeStorage.Insert(type, delta);
	}

	/**
	    \brief Increment type bucket on spawn.
	    \param zombie Infected instance
	*/
	static void OnSpawn(ZombieBase zombie)
	{
		if (!zombie)
			return;

		AddType(zombie.MetricZ_GetLabelTypeName(), 1);
	}

	/**
	    \brief Decrement buckets on delete.
	    \param zombie Infected instance
	    \param state Last known mind state
	*/
	static void OnDelete(ZombieBase zombie, int state)
	{
		AddState(state, -1);

		if (zombie)
			AddType(zombie.MetricZ_GetLabelTypeName(), -1);
	}

	/**
	    \brief Move one unit from old to new mind state.
	    \details Ignores unchanged states.
	    \param oldState Old mind state id
	    \param newState New mind state id
	*/
	static void OnStateChange(int oldState, int newState)
	{
		if (oldState == newState)
			return;

		AddState(oldState, -1);
		AddState(newState, 1);
	}

	/**
	    \brief Move from live mind state bucket to "dead".
	    \param lastState Last known live mind state
	*/
	static void OnKilled(int lastState)
	{
		if (lastState >= 0)
			AddState(lastState, -1);

		AddState(MINDSTATE_DEAD, 1);
	}

	/**
	    \brief Serialize all zombie data to JSON.
	    \param json Output JSON string.
	*/
	static void BuildJson(out string json)
	{
		EnsureNames();

		MetricZ_ZombieStatsDTO dto = new MetricZ_ZombieStatsDTO();

		foreach (int id, int stateCount : s_CountByState) {
			string stateName = "unknown";
			s_MindStates.Find(id, stateName);
			dto.by_state.Insert(stateName, stateCount);
		}

		foreach (string type, int typeCount : s_TypeStorage)
			dto.by_type.Insert(type, typeCount);

		JsonSerializer().WriteToString(dto, false, json);
	}
}
#endif
