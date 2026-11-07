/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER

class StateZ_AnimalStatsDTO
{
	void StateZ_AnimalStatsDTO()
	{
		by_type = new map<string, int>();
	}

	ref map<string, int> by_type;
}

/**
    \brief Animals-by-type gauges.
    \details Tracks live animals count per canonical type. Updated on EEInit/EEDelete.
*/
class StateZ_AnimalStats
{
	protected static ref map<string, int> s_CountByType = new map<string, int>(); // animal -> count

	/**
	    \brief Increment per-type count for spawned animal.
	*/
	static void OnSpawn(AnimalBase animal)
	{
		if (!animal)
			return;

		string type = animal.StateZ_GetLabelTypeName();

		int v;
		if (s_CountByType.Find(type, v))
			s_CountByType.Set(type, v + 1);
		else
			s_CountByType.Insert(type, 1);
	}

	/**
	    \brief Decrement per-type count for deleted animal.
	*/
	static void OnDelete(AnimalBase animal)
	{
		if (!animal)
			return;

		string type = animal.StateZ_GetLabelTypeName();

		int v;
		if (!s_CountByType.Find(type, v))
			return;

		v = v - 1;
		if (v <= 0)
			s_CountByType.Remove(type);
		else
			s_CountByType.Set(type, v);
	}

	/**
	    \brief Serialize all animal data to JSON.
	    \param json Output JSON string.
	*/
	static void BuildJson(out string json)
	{
		StateZ_AnimalStatsDTO dto = new StateZ_AnimalStatsDTO();
		dto.by_type = s_CountByType;
		JsonSerializer().WriteToString(dto, false, json);
	}
}
#endif
