/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Helpers for configuration
*/
class MetricZ_Helpers
{
	/**
	    \brief Get Players limit from server config
	*/
	static int GetLimitPlayers()
	{
		return g_Game.ServerConfigGetInt("maxPlayers");
	}

	/**
	    \brief Get FPS limit from CLI param
	*/
	static int GetLimitFPS()
	{
		string limit;
		if (GetCLIParam("limitFPS", limit)) {
			int fps = limit.ToInt();
			if (fps > 0)
				return fps;
		}

		return 0;
	}

	/**
	    \brief Get instance ID if input empty string, use game config or port number
	*/
	static string GetInstanceID(string id)
	{
#ifdef DIAG
		ErrorEx("MetricZ: requested instance ID with '" + id + "' input", ErrorExSeverity.WARNING);
#endif

		// get override from input
		if (id != string.Empty)
			return id;

		// get configured instance id in serverDZ.cfg
		id = g_Game.ServerConfigGetInt("instanceId").ToString();
		if (id != "0")
			return id;

		// use game port as instance id
		GetCLIParam("port", id);
		if (id != "0") {
			ErrorEx(
			    "MetricZ: instanceId is 0. Used game port " + id + " as instanceId.",
			    ErrorExSeverity.INFO);

			return id;
		}

		// use steam query port as instance id
		id = g_Game.ServerConfigGetInt("steamQueryPort").ToString();
		if (id != "0") {
			ErrorEx(
			    "MetricZ: instanceId and game port is 0. Used steam query port " + id + " as instanceId.",
			    ErrorExSeverity.INFO);

			return id;
		}

		// fallback to 0
		ErrorEx(
		    "MetricZ: instanceId is 0. Set unique 'instanceId' in serverDZ.cfg or override 'instance_id' in metricz.json",
		    ErrorExSeverity.WARNING);

		return "0";
	}

}
#endif
