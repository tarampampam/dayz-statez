/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief MetricZ Constants
*/
class MetricZ_Constants
{
	// Build information
	static const string NAME = "MetricZ";
	static const string VERSION = "dev";
	static const string COMMIT_SHA = "fffffff";
	static const string BUILD_DATE = "1970-01-01T00:00:00+00:00";

	// Working directory
	static const string WORK_DIR = "$profile:metricz/";
	static const string CONFIG_FILE = WORK_DIR + "config.json";

	// File export directory
	static const string EXPORT_DIR = WORK_DIR + "export/";
}
#endif
