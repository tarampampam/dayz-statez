/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Configuration Manager for StateZ.
    \details Handles the lifecycle of the configuration: loading from JSON,
             creating default files, and ensuring the output directory exists.
             Implements the Singleton pattern.
*/
class StateZ_Config
{
	protected static ref StateZ_Config s_Instance; //!< singleton instance
	protected static ref StateZ_ConfigDTO s_Config; //!< singleton state
	protected static bool s_Loaded; //!< configuration loaded flag

	/**
	    \brief Retrieves the singleton configuration instance.
	    \details Lazy-loads the configuration file on the first call.
	    \return Global `StateZ_ConfigDTO` instance.
	*/
	static StateZ_ConfigDTO Get()
	{
		if (!s_Instance)
			s_Instance = new StateZ_Config();

		if (!s_Config && !s_Loaded) {
			s_Config = new StateZ_ConfigDTO();
			s_Instance.Load();
		}

		return s_Config;
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
		InitJsonExport();
	}

	/**
	    \brief Handles JSON file operations (Load/Save/Upgrade).
	    \details - If the file exists, it attempts to load and validate it.
	             - If the file is missing, it creates a new default configuration.
	*/
	protected void LoadConfigFile()
	{
		if (!FileExist(StateZ_Constants.WORK_DIR))
			MakeDirectory(StateZ_Constants.WORK_DIR);

		string cfgFile = StateZ_Constants.CONFIG_FILE;
		string error;
		if (FileExist(cfgFile)) {
			if (!JsonFileLoader<StateZ_ConfigDTO>.LoadFile(cfgFile, s_Config, error)) {
				ErrorEx(error);
				return;
			}

			SetLoaded();
			return;
		}

		if (!JsonFileLoader<StateZ_ConfigDTO>.SaveFile(cfgFile, s_Config, error)) {
			ErrorEx(error);
			return;
		}

		ErrorEx(
		    string.Format("StateZ: saved new config file: %1", cfgFile),
		    ErrorExSeverity.INFO);
		SetLoaded();
	}

	/**
	    \brief Ensures the JSON output directory exists.
	*/
	protected void InitJsonExport()
	{
		if (!FileExist(StateZ_Constants.DATA_DIR))
			MakeDirectory(StateZ_Constants.DATA_DIR);
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
		                 "StateZ: loaded %1 (%2) build %3",
		                 StateZ_Constants.VERSION,
		                 StateZ_Constants.COMMIT_SHA,
		                 StateZ_Constants.BUILD_DATE);

		ErrorEx(msg, ErrorExSeverity.INFO);
	}

	/**
	    \brief Dumps the current configuration to the script log
	*/
	protected void DebugConfig()
	{
		string json;
		JsonSerializer().WriteToString(s_Config, true, json);
		ErrorEx("StateZ: configuration trace:\n" + json, ErrorExSeverity.INFO);
	}
}
#endif
