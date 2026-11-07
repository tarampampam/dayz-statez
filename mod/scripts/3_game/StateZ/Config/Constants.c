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
	static const string NAME = "StateZ";
	static const string VERSION = "dev";
	static const string COMMIT_SHA = "fffffff";
	static const string BUILD_DATE = "1970-01-01T00:00:00+00:00";

	// Working directory
	static const string WORK_DIR = "$profile:statez/";
	static const string CONFIG_FILE = WORK_DIR + "config.json";

	// File export directory
	static const string EXPORT_DIR = WORK_DIR + "export/";
}
#endif
