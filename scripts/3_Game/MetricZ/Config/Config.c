/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Configuration Manager for MetricZ.
    \details Handles the lifecycle of the configuration: loading from JSON, version migration,
             creating default files, and initializing export paths.
             Implements the Singleton pattern.
*/
class MetricZ_Config
{
	protected static ref MetricZ_Config s_Instance; //!< singleton instance
	protected static ref MetricZ_ConfigDTO s_Config; //!< singleton state
	protected static bool s_Loaded; //!< configuration loaded flag

	/**
	    \brief Retrieves the singleton configuration instance.
	    \details Lazy-loads the configuration file on the first call.
	    \return Global `MetricZ_ConfigDTO` instance.
	*/
	static MetricZ_ConfigDTO Get()
	{
		if (!s_Instance)
			s_Instance = new MetricZ_Config();

		if (!s_Config && !s_Loaded) {
			s_Config = new MetricZ_ConfigDTO();
			s_Instance.Load();
		}

		return s_Config;
	}

	/**
	    \brief Drop instance and force reload on next
	*/
	static void Reset()
	{
		if (s_Config)
			s_Config = null;

		s_Loaded = false;
		ErrorEx("MetricZ: configuration reset", ErrorExSeverity.INFO);
	}

	/**
	    \brief Checks if the configuration is successfully loaded.
	*/
	static bool IsLoaded()
	{
		return (s_Instance && s_Config && s_Loaded);
	}

	/**
	    \brief Orchestrates the loading process.
	*/
	protected void Load()
	{
		LoadConfigFile();

		if (s_Config.json.enabled)
			InitJsonExport();
	}

	/**
	    \brief Handles JSON file operations (Load/Save/Upgrade).
	    \details - If the file exists, it attempts to load and validate it.
	             - If the version mismatches, it performs an upgrade and saves the file.
	             - If the file is missing, it creates a new default configuration.
	*/
	protected void LoadConfigFile()
	{
		if (!FileExist(MetricZ_Constants.WORK_DIR))
			MakeDirectory(MetricZ_Constants.WORK_DIR);

		string cfgFile = MetricZ_Constants.CONFIG_FILE;
		string error;
		if (FileExist(cfgFile)) {
			if (!JsonFileLoader<MetricZ_ConfigDTO>.LoadFile(cfgFile, s_Config, error)) {
				ErrorEx(error);
				return;
			}

			s_Config.Normalize();

			if (s_Config.version != MetricZ_Constants.VERSION) {
				s_Config.version = MetricZ_Constants.VERSION;

				if (!JsonFileLoader<MetricZ_ConfigDTO>.SaveFile(cfgFile, s_Config, error)) {
					ErrorEx(error);
					return;
				}

				ErrorEx(
				    string.Format("MetricZ: saved upgraded config file: %1", cfgFile),
				    ErrorExSeverity.INFO);
			}

			SetLoaded();
			return;
		}

		s_Config.version = MetricZ_Constants.VERSION;

		if (!JsonFileLoader<MetricZ_ConfigDTO>.SaveFile(cfgFile, s_Config, error)) {
			ErrorEx(error);
			return;
		}

		ErrorEx(
		    string.Format("MetricZ: saved new config file: %1", cfgFile),
		    ErrorExSeverity.INFO);
		SetLoaded();
	}

	/**
	    \brief Determines the export directory for JSON files.
	*/
	protected void InitJsonExport()
	{
		string exportBase = MetricZ_Constants.EXPORT_DIR;
		if (!FileExist(exportBase))
			MakeDirectory(exportBase);

		string dirName = s_Config.json.file_name;
		if (dirName == string.Empty)
			dirName = string.Format("metricz_%1", s_Config.settings.instance_id_resolved);

		string exportDir = string.Format("%1%2/", exportBase, dirName);
		if (!FileExist(exportDir))
			MakeDirectory(exportDir);

		s_Config.json.export_dir = exportDir;
	}

	/**
	    \brief Finalizes the loading process.
	*/
	protected void SetLoaded()
	{
		if (s_Config)
			s_Config.Normalize();

		s_Loaded = true;

#ifdef DIAG
		DebugConfig();
#endif

		string msg = string.Format(
		                 "MetricZ: loaded %1 (%2) build %3, instance_id=%4",
		                 MetricZ_Constants.VERSION,
		                 MetricZ_Constants.COMMIT_SHA,
		                 MetricZ_Constants.BUILD_DATE,
		                 s_Config.settings.instance_id_resolved);

		ErrorEx(msg, ErrorExSeverity.INFO);
	}

	/**
	    \brief Dumps the current configuration to the script log
	*/
	protected void DebugConfig()
	{
		string json;
		JsonSerializer().WriteToString(s_Config, true, json);
		ErrorEx("MetricZ: configuration trace:\n" + json, ErrorExSeverity.INFO);
	}
}
#endif
