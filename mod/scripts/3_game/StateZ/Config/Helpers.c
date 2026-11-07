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
class StateZ_Helpers
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
}
#endif
