/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief StateZ Constants
*/
class StateZ_Constants
{
	// Build information
	static const string VERSION = "dev";
	static const string COMMIT_SHA = "fffffff";
	static const string BUILD_DATE = "1970-01-01T00:00:00+00:00";

	// Working directory
	static const string WORK_DIR = "$profile:statez/";
	static const string CONFIG_FILE = WORK_DIR + "config.json";

	// JSON output directory. Static path under the server profile.
	// To store JSON outside the DayZ server root, replace this directory
	// (<profiles>/statez/data) with a symlink to the desired location.
	static const string DATA_DIR = WORK_DIR + "data/";
}
#endif
