/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Background scheduler for StateZ.
    \details Owns scheduling and coordination between collectors.
             Entry points: Init() on mission start, Shutdown() on mission finish.
*/
class StateZ_Scheduler
{
	protected static ref StateZ_Scheduler s_Instance; // singleton instance

	protected int m_FlushStep; //!< Flush State Machine - step number
	protected bool s_Busy; //!< Guard prevents overlapping collections
	protected ref array<ref StateZ_CollectorBase> m_Collectors; //!< Registry of all metric collectors

	/**
	    \brief Initialize StateZ and schedule first collection.
	    \details Loads config, initializes storage and queues first Update().
	*/
	void StateZ_Scheduler()
	{
		if (s_Instance)
			return;

		m_Collectors = new array<ref StateZ_CollectorBase>();

		if (!StateZ_Config.IsLoaded())
			StateZ_Config.Get();

		StateZ_Storage.Init();

		g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).CallLater(
		    Update,
		    StateZ_Config.Get().settings.init_delay_sec * 1000,
		    false);

		// Register internal collectors
		RegisterCollector(new StateZ_CollectorWorld());
		RegisterCollector(new StateZ_CollectorPlayers());
		RegisterCollector(new StateZ_CollectorZombies());
		RegisterCollector(new StateZ_CollectorAnimals());
		RegisterCollector(new StateZ_CollectorTransports());
		RegisterCollector(new StateZ_CollectorWeapons());
		RegisterCollector(new StateZ_CollectorHits());
		RegisterCollector(new StateZ_CollectorTerritories());
		RegisterCollector(new StateZ_CollectorAreas());
		RegisterCollector(new StateZ_CollectorEvents());

		ErrorEx("StateZ: loaded with " + m_Collectors.Count() + " collectors", ErrorExSeverity.INFO);
	}

	/**
	    \brief Access global scheduler instance.
	*/
	static StateZ_Scheduler Get()
	{
		if (!s_Instance)
			s_Instance = new StateZ_Scheduler();

		return s_Instance;
	}

	/**
	    \brief API for third party mods to register collectors.
	    \param collector StateZ_CollectorBase collector instance.
	    \details Mods MUST call this in `MissionServer::OnInit` (after `super.OnInit()`) or later.
	*/
	static void Register(StateZ_CollectorBase collector)
	{
		if (!s_Instance) {
			ErrorEx(
			    "StateZ: Attempt to register collector '" + collector.GetName() + "' too early! Register in MissionServer::OnInit after super.OnInit()",
			    ErrorExSeverity.WARNING);
			return;
		}

		s_Instance.RegisterCollector(collector);
	}

	/**
	    \brief Graceful shutdown. Deletes all JSON files.
	*/
	void Shutdown()
	{
		if (!StateZ_Config.IsLoaded())
			return;

		ErrorEx("StateZ: shutting down", ErrorExSeverity.INFO);

		g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).Remove(ProcessFlushStep);
		g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).Remove(Update);

		s_Instance = null;
		s_Busy = false;

		StateZ_JsonWriter.DeleteAll();
	}

	/**
	    \brief Internal collector registration logic
	    \param collector StateZ_CollectorBase collector instance.
	*/
	protected void RegisterCollector(StateZ_CollectorBase collector)
	{
		if (!collector)
			return;

		if (!collector.IsEnabled()) {
#ifdef DIAG
			ErrorEx("StateZ: Skip disabled collector: " + collector.GetName(), ErrorExSeverity.INFO);
#endif
			return;
		}

		m_Collectors.Insert(collector);

#ifdef DIAG
		ErrorEx("StateZ: Registered collector: " + collector.GetName(), ErrorExSeverity.INFO);
#endif
	}

	/**
	    \brief Periodic collection task. Runs all collectors one-per-frame.
	    \details
	        - Reschedules itself each call using fixed interval (minimizes drift).
	        - Skips execution if a previous collection is still running (s_Busy).
	*/
	protected void Update()
	{
		// Schedule next tick to minimize drift
		g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).CallLater(
		    Update,
		    StateZ_Config.Get().settings.collect_interval_sec * 1000,
		    false);

		if (!StateZ_Config.IsLoaded())
			return;

		if (s_Busy) {
			ErrorEx("StateZ: skip collection, previous cycle still running", ErrorExSeverity.WARNING);
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

		StateZ_CollectorBase collector = m_Collectors.Get(m_FlushStep);
		if (collector && collector.IsEnabled()) {
#ifdef DIAG
			float t = g_Game.GetTickTime();
#endif
			collector.Flush();
#ifdef DIAG
			ErrorEx(
			    string.Format("StateZ: Flush %1 in frame %2 took: %3ms", collector.GetName(), m_FlushStep, (g_Game.GetTickTime() - t) * 1000),
			    ErrorExSeverity.INFO);
#endif
		}

		m_FlushStep++;
		g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).Call(ProcessFlushStep);
	}

	/**
	    \brief Finalizes the collection cycle.
	*/
	protected void FinishFlush()
	{
		s_Busy = false;

#ifdef DIAG
		ErrorEx("StateZ: FinishFlush complete, step=" + m_FlushStep, ErrorExSeverity.INFO);
#endif
	}
}
#endif
