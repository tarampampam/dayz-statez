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
class StateZ_EntitiesWriter
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

		array<ref StateZ_PlayerSnapshot> snaps = new array<ref StateZ_PlayerSnapshot>();

		foreach (Man man : s_PlayersListBuffer) {
			PlayerBase player;
			if (!Class.CastTo(player, man))
				continue;

			StateZ_PlayerMetrics pm = player.StateZ_GetMetrics();
			if (!pm)
				continue;

			pm.Update();

			StateZ_PlayerSnapshot snap = new StateZ_PlayerSnapshot();
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
		array<Transport> list = StateZ_TransportRegistry.GetList();

		if (!list || list.Count() == 0) {
			json = "[]";
			return;
		}

		array<ref StateZ_TransportSnapshot> snaps = new array<ref StateZ_TransportSnapshot>();

		foreach (Transport transport : list) {
			if (!transport)
				continue;

			StateZ_TransportMetrics tm;

			CarScript car = CarScript.Cast(transport);
			if (car)
				tm = car.StateZ_GetMetrics();
			else {
				BoatScript boat = BoatScript.Cast(transport);
				if (boat)
					tm = boat.StateZ_GetMetrics();
				else {
					HelicopterScript heli = HelicopterScript.Cast(transport);
					if (heli)
						tm = heli.StateZ_GetMetrics();
#ifdef EXPANSIONMODVEHICLE
					else {
						ExpansionVehicleBase expansionVeh = ExpansionVehicleBase.Cast(transport);
						if (expansionVeh)
							tm = expansionVeh.StateZ_GetMetrics();
					}
#endif
				}
			}

			if (!tm)
				continue;

			tm.Update();

			StateZ_TransportSnapshot snap = new StateZ_TransportSnapshot();
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
		array<TerritoryFlag> list = StateZ_TerritoryRegistry.GetList();

		if (!list || list.Count() == 0) {
			json = "[]";
			return;
		}

		array<ref StateZ_TerritorySnapshot> snaps = new array<ref StateZ_TerritorySnapshot>();

		foreach (TerritoryFlag territory : list) {
			if (!territory)
				continue;

			StateZ_TerritoryMetrics fm = territory.StateZ_GetMetrics();
			if (!fm)
				continue;

			fm.Update();

			StateZ_TerritorySnapshot snap = new StateZ_TerritorySnapshot();
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
		array<EffectArea> list = StateZ_EffectAreaRegistry.GetList();

		if (!list || list.Count() == 0) {
			json = "[]";
			return;
		}

		array<ref StateZ_AreaSnapshot> snaps = new array<ref StateZ_AreaSnapshot>();

		foreach (EffectArea area : list) {
			if (!area)
				continue;

			StateZ_EffectAreaMetrics am = area.StateZ_GetMetrics();
			if (!am)
				continue;

			am.Update();

			StateZ_AreaSnapshot snap = new StateZ_AreaSnapshot();
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
