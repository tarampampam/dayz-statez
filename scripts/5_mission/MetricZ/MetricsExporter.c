/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Background exporter for MetricZ.
    \details Owns scheduling and coordination between collectors.
             Entry points: Init() on mission start, Shutdown() on mission finish.
*/
class MetricZ_Exporter
{
	protected static ref MetricZ_Exporter s_Instance; // singleton instance

	protected int m_FlushStep; //!< Flush State Machine - step number
	protected bool s_Busy; //!< Guard prevents overlapping scrapes
	protected ref array<ref MetricZ_CollectorBase> m_Collectors; //!< Registry of all metric collectors

	/**
	    \brief Initialize MetricZ and schedule first scrape.
	    \details Loads config, initializes storage and queues first Update().
	*/
	void MetricZ_Exporter()
	{
		if (s_Instance)
			return;

		m_Collectors = new array<ref MetricZ_CollectorBase>();

		if (!MetricZ_Config.IsLoaded())
			MetricZ_Config.Get();

		MetricZ_Storage.Init();

		g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).CallLater(
		    Update,
		    MetricZ_Config.Get().settings.init_delay_sec * 1000,
		    false);

		// Register internal collectors
		RegisterCollector(new MetricZ_CollectorWorld());
		RegisterCollector(new MetricZ_CollectorPlayers());
		RegisterCollector(new MetricZ_CollectorZombies());
		RegisterCollector(new MetricZ_CollectorAnimals());
		RegisterCollector(new MetricZ_CollectorTransports());
		RegisterCollector(new MetricZ_CollectorWeapons());
		RegisterCollector(new MetricZ_CollectorHits());
		RegisterCollector(new MetricZ_CollectorTerritories());
		RegisterCollector(new MetricZ_CollectorAreas());
		RegisterCollector(new MetricZ_CollectorRPC());
		RegisterCollector(new MetricZ_CollectorEvents());

		ErrorEx("MetricZ: loaded with " + m_Collectors.Count() + " collectors", ErrorExSeverity.INFO);
	}

	/**
	    \brief Access global exporter instance.
	*/
	static MetricZ_Exporter Get()
	{
		if (!s_Instance)
			s_Instance = new MetricZ_Exporter();

		return s_Instance;
	}

	/**
	    \brief API for third party mods to register collectors.
	    \param collector MetricZ_CollectorBase collector instance.
	    \details Mods MUST call this in `MissionServer::OnInit` (after `super.OnInit()`) or later.
	*/
	static void Register(MetricZ_CollectorBase collector)
	{
		if (!s_Instance) {
			ErrorEx(
			    "MetricZ: Attempt to register collector '" + collector.GetName() + "' too early! Register in MissionServer::OnInit after super.OnInit()",
			    ErrorExSeverity.WARNING);
			return;
		}

		s_Instance.RegisterCollector(collector);
	}

	/**
	    \brief Graceful shutdown. Writes status=0 or deletes all JSON files.
	*/
	void Shutdown()
	{
		if (!MetricZ_Config.IsLoaded())
			return;

		ErrorEx("MetricZ: scrape shutting down", ErrorExSeverity.INFO);

		g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).Remove(ProcessFlushStep);
		g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).Remove(Update);

		s_Instance = null;
		s_Busy = false;

		MetricZ_ConfigDTO cfg = MetricZ_Config.Get();

		if (cfg.json.delete_on_shutdown) {
			MetricZ_JsonWriter.DeleteAll(cfg.json.export_dir);
			return;
		}

		// signal offline status
		string json = string.Format("{\"status\":0,\"ts\":%1}", MetricZ_Time.GameEpochSeconds());
		MetricZ_JsonWriter.Write(cfg.json.export_dir, "status", json);
	}

	/**
	    \brief Internal collector registration logic
	    \param collector MetricZ_CollectorBase collector instance.
	*/
	protected void RegisterCollector(MetricZ_CollectorBase collector)
	{
		if (!collector)
			return;

		if (!collector.IsEnabled()) {
#ifdef DIAG
			ErrorEx("MetricZ: Skip disabled collector: " + collector.GetName(), ErrorExSeverity.INFO);
#endif
			return;
		}

		m_Collectors.Insert(collector);

#ifdef DIAG
		ErrorEx("MetricZ: Registered collector: " + collector.GetName(), ErrorExSeverity.INFO);
#endif
	}

	/**
	    \brief Periodic scrape task. Runs all collectors one-per-frame.
	    \details
	        - Reschedules itself each call using fixed interval (minimizes drift).
	        - Skips execution if a previous scrape is still running (s_Busy).
	*/
	protected void Update()
	{
		// Schedule next tick to minimize drift
		g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).CallLater(
		    Update,
		    MetricZ_Config.Get().settings.collect_interval_sec * 1000,
		    false);

		if (!MetricZ_Config.IsLoaded())
			return;

		if (s_Busy) {
			MetricZ_Storage.s_ScrapeSkippedTotal += 1;
			ErrorEx("MetricZ: skip scrape, previous cycle still running", ErrorExSeverity.WARNING);
			return;
		}

		s_Busy = true;
		m_FlushStep = 0;

		// enter Flush State Machine
		ProcessFlushStep();
	}

	/**
	    \brief Executes ONE collector and schedules the next.
	*/
	protected void ProcessFlushStep()
	{
		if (m_FlushStep >= m_Collectors.Count()) {
			FinishFlush();
			return;
		}

		MetricZ_CollectorBase collector = m_Collectors.Get(m_FlushStep);
		if (collector && collector.IsEnabled()) {
#ifdef DIAG
			float t = g_Game.GetTickTime();
#endif
			collector.Flush();
#ifdef DIAG
			ErrorEx(
			    string.Format("MetricZ: Flush %1 in frame %2 took: %3ms", collector.GetName(), m_FlushStep, (g_Game.GetTickTime() - t) * 1000),
			    ErrorExSeverity.INFO);
#endif
		}

		m_FlushStep++;
		g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).Call(ProcessFlushStep);
	}

	/**
	    \brief Finalizes the scrape cycle.
	*/
	protected void FinishFlush()
	{
		s_Busy = false;

#ifdef DIAG
		ErrorEx("MetricZ: FinishFlush complete, step=" + m_FlushStep, ErrorExSeverity.INFO);
#endif
	}
}
#endif
