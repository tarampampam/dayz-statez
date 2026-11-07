/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Base class for MetricZ collectors.
    \details Provides common functionality for all collectors.
*/
class MetricZ_CollectorBase
{
	/**
	    \brief Domain name used as the JSON file base name.
	    \return string Name of this collector's domain.
	*/
	string GetName()
	{
		return "unknown";
	}

	/**
	    \brief Check if this collector is enabled.
	    \return bool True if the collector is enabled.
	*/
	bool IsEnabled()
	{
		return MetricZ_Config.IsLoaded();
	}

	/**
	    \brief Main flush method called by the exporter each scrape cycle.
	*/
	void Flush()
	{
		// override in subclasses
	}
}
#endif
