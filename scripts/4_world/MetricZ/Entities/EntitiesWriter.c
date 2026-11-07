/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Helpers for building JSON arrays from entity metrics.
*/
class MetricZ_EntitiesWriter
{
	static ref array<Man> s_PlayersListBuffer = new array<Man>();

	/**
	    \brief Build a JSON array of all online players.
	    \param json Output JSON string (empty array if no players).
	*/
	static void BuildPlayersJson(out string json)
	{
		s_PlayersListBuffer.Clear();
		g_Game.GetPlayers(s_PlayersListBuffer);

		array<ref MetricZ_PlayerSnapshot> snaps = new array<ref MetricZ_PlayerSnapshot>();

		foreach (Man man : s_PlayersListBuffer) {
			PlayerBase player;
			if (!Class.CastTo(player, man))
				continue;

			MetricZ_PlayerMetrics pm = player.MetricZ_GetMetrics();
			if (!pm)
				continue;

			pm.Update();

			MetricZ_PlayerSnapshot snap = new MetricZ_PlayerSnapshot();
			pm.ToSnapshot(snap);
			snaps.Insert(snap);
		}

		if (snaps.Count() == 0) {
			json = "[]";
			return;
		}

		JsonSerializer().WriteToString(snaps, false, json);
	}

	/**
	    \brief Build a JSON array of all tracked transports.
	    \param json Output JSON string (empty array if none).
	*/
	static void BuildTransportJson(out string json)
	{
		array<Transport> list = MetricZ_TransportRegistry.GetList();

		if (!list || list.Count() == 0) {
			json = "[]";
			return;
		}

		array<ref MetricZ_TransportSnapshot> snaps = new array<ref MetricZ_TransportSnapshot>();

		foreach (Transport transport : list) {
			if (!transport)
				continue;

			MetricZ_TransportMetrics tm;

			CarScript car = CarScript.Cast(transport);
			if (car)
				tm = car.MetricZ_GetMetrics();
			else {
				BoatScript boat = BoatScript.Cast(transport);
				if (boat)
					tm = boat.MetricZ_GetMetrics();
				else {
					HelicopterScript heli = HelicopterScript.Cast(transport);
					if (heli)
						tm = heli.MetricZ_GetMetrics();
#ifdef EXPANSIONMODVEHICLE
					else {
						ExpansionVehicleBase expansionVeh = ExpansionVehicleBase.Cast(transport);
						if (expansionVeh)
							tm = expansionVeh.MetricZ_GetMetrics();
					}
#endif
				}
			}

			if (!tm)
				continue;

			tm.Update();

			MetricZ_TransportSnapshot snap = new MetricZ_TransportSnapshot();
			tm.ToSnapshot(snap);
			snaps.Insert(snap);
		}

		if (snaps.Count() == 0) {
			json = "[]";
			return;
		}

		JsonSerializer().WriteToString(snaps, false, json);
	}

	/**
	    \brief Build a JSON array of all territory flags.
	    \param json Output JSON string (empty array if none).
	*/
	static void BuildTerritoryJson(out string json)
	{
		array<TerritoryFlag> list = MetricZ_TerritoryRegistry.GetList();

		if (!list || list.Count() == 0) {
			json = "[]";
			return;
		}

		array<ref MetricZ_TerritorySnapshot> snaps = new array<ref MetricZ_TerritorySnapshot>();

		foreach (TerritoryFlag territory : list) {
			if (!territory)
				continue;

			MetricZ_TerritoryMetrics fm = territory.MetricZ_GetMetrics();
			if (!fm)
				continue;

			fm.Update();

			MetricZ_TerritorySnapshot snap = new MetricZ_TerritorySnapshot();
			fm.ToSnapshot(snap);
			snaps.Insert(snap);
		}

		if (snaps.Count() == 0) {
			json = "[]";
			return;
		}

		JsonSerializer().WriteToString(snaps, false, json);
	}

	/**
	    \brief Build a JSON array of all effect areas.
	    \param json Output JSON string (empty array if none).
	*/
	static void BuildAreasJson(out string json)
	{
		array<EffectArea> list = MetricZ_EffectAreaRegistry.GetList();

		if (!list || list.Count() == 0) {
			json = "[]";
			return;
		}

		array<ref MetricZ_AreaSnapshot> snaps = new array<ref MetricZ_AreaSnapshot>();

		foreach (EffectArea area : list) {
			if (!area)
				continue;

			MetricZ_EffectAreaMetrics am = area.MetricZ_GetMetrics();
			if (!am)
				continue;

			am.Update();

			MetricZ_AreaSnapshot snap = new MetricZ_AreaSnapshot();
			am.ToSnapshot(snap);
			snaps.Insert(snap);
		}

		if (snaps.Count() == 0) {
			json = "[]";
			return;
		}

		JsonSerializer().WriteToString(snaps, false, json);
	}
}
#endif
