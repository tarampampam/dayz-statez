/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

class CfgMods
{
	class StateZ
	{
		type = "mod";
		dir = "statez";
		name = "StateZ - Server State Writer";
		version = "0.1.0";
		author = "tarampampam";
		credits = "WoozyMasta";
		hideName = 1;
		hidePicture = 1;
		defines[] = {"STATEZ"};
		dependencies[] = {"game", "world", "mission"};

		class defs
		{
			class gameScriptModule
			{
				files[] = {"statez/scripts/3_game"};
			};

			class worldScriptModule
			{
				files[] = { "statez/scripts/4_world" };
			};

			class missionScriptModule
			{
				files[] = { "statez/scripts/5_mission" };
			};
		};
	};
};

class CfgPatches
{
	class StateZ
	{
		requiredAddons[] = {
			"DZ_Data",
			"DZ_Scripts",
		};
	};
};
