/*
    SPDX-License-Identifier: GPL-3.0-or-later
*/

#ifdef SERVER
/**
    \brief Shared vehicle lifecycle logic for all transport families.
    \details The engine forbids modding the common Transport base from this
             module (modded classes must live in the base class's module, and
             Transport is in 3_Game), so each 4_World vehicle script
             (CarScript/BoatScript/HelicopterScript/MotorbikeScript, +Expansion)
             keeps a thin modded shim that forwards here. All counting, registry
             and categorization lives in this one place.
*/
class StateZ_TransportHooks
{
	/**
	    \brief Spawn path: bump live count, create metrics, register, schedule identity.
	    \return New metrics instance for the shim to hold.
	*/
	static StateZ_TransportMetrics OnInit(Transport transport)
	{
		AddLive(transport, 1);

		StateZ_TransportMetrics m = new StateZ_TransportMetrics();
		// Persistence hash is not guaranteed yet on EEInit; resolved later.
		m.InitLater(transport);

		StateZ_TransportRegistry.Register(m);
		return m;
	}

	/**
	    \brief Persistence-load path: (re)init metrics with the real hash. Does not
	           count or register — EEInit already did that for persistent entities.
	    \param existing Metrics already created by EEInit, or null.
	    \return Metrics instance for the shim to hold.
	*/
	static StateZ_TransportMetrics OnAfterLoad(Transport transport, StateZ_TransportMetrics existing)
	{
		StateZ_TransportMetrics m = existing;
		if (!m)
			m = new StateZ_TransportMetrics();

		m.Init(transport);
		return m;
	}

	/**
	    \brief Delete path: drop live count and unregister the metrics.
	*/
	static void OnDelete(Transport transport, StateZ_TransportMetrics m)
	{
		// No metrics means OnInit never ran (e.g. config loaded after EEInit), so
		// nothing was counted/registered — skip to avoid a counter going negative.
		if (!m)
			return;

		AddLive(transport, -1);
		StateZ_TransportRegistry.Unregister(m);
	}

	/**
	    \brief Kill path: bump the matching destroyed counter.
	*/
	static void OnKilled(Transport transport)
	{
		AddDestroyed(transport);
	}

	/**
	    \brief Classify a vehicle into a world-count bucket.
	    \return 0 = car (default/catch-all), 1 = boat, 2 = helicopter, 3 = motorbike.
	*/
	static int Kind(Transport transport)
	{
#ifdef EXPANSIONMODVEHICLE
		// Expansion custom-sim vehicles expose their own category predicates.
		ExpansionVehicleBase ev;
		if (Class.CastTo(ev, transport)) {
			if (ev.Expansion_IsBoat())
				return 1;
			if (ev.Expansion_IsHelicopter() || ev.Expansion_IsPlane())
				return 2;
			return 0;
		}

		// Expansion vehicles built on CarScript still carry their boat/heli identity.
		if (transport.IsInherited(ExpansionBoatScript))
			return 1;
		if (transport.IsInherited(ExpansionHelicopterScript))
			return 2;
#endif

		if (transport.IsInherited(BoatScript))
			return 1;
		if (transport.IsInherited(HelicopterScript))
			return 2;
		if (transport.IsInherited(MotorbikeScript))
			return 3;

		return 0;
	}

	protected static void AddLive(Transport transport, int delta)
	{
		switch (Kind(transport)) {
		case 1:
			StateZ_Storage.s_Boats += delta;
			break;
		case 2:
			StateZ_Storage.s_Helicopters += delta;
			break;
		case 3:
			StateZ_Storage.s_Motorbikes += delta;
			break;
		default:
			StateZ_Storage.s_Cars += delta;
			break;
		}
	}

	protected static void AddDestroyed(Transport transport)
	{
		switch (Kind(transport)) {
		case 1:
			StateZ_Storage.s_BoatsDestroys += 1;
			break;
		case 2:
			StateZ_Storage.s_HelicoptersDestroys += 1;
			break;
		case 3:
			StateZ_Storage.s_MotorbikesDestroys += 1;
			break;
		default:
			StateZ_Storage.s_CarsDestroys += 1;
			break;
		}
	}
}
#endif
