/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Per-Transport data collector.
    \details Stores plain fields updated during collection. Fills snapshot on demand.
*/
class StateZ_TransportMetrics
{
	protected Transport m_Transport;

	// identity
	protected string m_ClassName;
	protected string m_VehicleType;
	protected int m_Hash;

	// state
	float m_Health;
	int m_Passengers;
	float m_SpeedMS;
	bool m_EngineOn;
	float m_FuelFraction;

	// position
	float m_PosX;
	float m_PosY;
	float m_PosZ;
	float m_Yaw;

	/**
	    \brief One-time init with transport reference.
	*/
	void Init(Transport transport)
	{
		if (transport) {
			m_Transport = transport;
			BuildIdentity();
		}
	}

	/**
	    \brief One-time init called after a delay.
	*/
	void InitLater(Transport transport)
	{
		if (transport) {
			m_Transport = transport;
			g_Game.GetCallQueue(CALL_CATEGORY_SYSTEM).CallLater(
			    BuildIdentity,
			    StateZ_Config.Get().interval * 1000,
			    false);
		}
	}

	/**
	    \brief Update all fields from the transport state.
	*/
	void Update()
	{
		if (!m_Transport)
			return;

		m_Health = m_Transport.GetHealth01();
		m_Passengers = GetPassengersCount();

		// position
		vector pos = m_Transport.GetPosition();
		m_PosX = pos[0];
		m_PosZ = pos[2];
		m_PosY = pos[1];
		m_Yaw = m_Transport.GetOrientation()[0];

		// speed
		float speed = GetVelocity(m_Transport).Length();
		if (speed < 0.10)
			speed = 0;
		m_SpeedMS = speed;

		m_EngineOn = false;
		m_FuelFraction = -1.0;

		CarScript car;
		if (Class.CastTo(car, m_Transport) && car.IsVitalFuelTank()) {
			m_EngineOn = car.EngineIsOn();
			m_FuelFraction = car.GetFluidFraction(CarFluid.FUEL);
			return;
		}

		BoatScript boat;
		if (Class.CastTo(boat, m_Transport) && boat.HasEngine()) {
			if (!boat.EngineIsOn() && speed < 0.30)
				m_SpeedMS = 0;

			m_EngineOn = boat.EngineIsOn();
			m_FuelFraction = boat.GetFluidFraction(BoatFluid.FUEL);
			return;
		}

#ifdef EXPANSIONMODVEHICLE
		ExpansionVehicleBase expVehicle;
		if (Class.CastTo(expVehicle, m_Transport)) {
			m_EngineOn = expVehicle.Expansion_EngineIsOn();
			m_FuelFraction = expVehicle.GetFluidFraction(CarFluid.FUEL);
			return;
		}
#endif
	}

	/**
	    \brief Fill a transport snapshot with current field values.
	    \param snap Snapshot instance to populate.
	*/
	void ToSnapshot(StateZ_TransportSnapshot snap)
	{
		if (!snap)
			return;

		snap.class_name = m_ClassName;
		snap.vehicle_type = m_VehicleType;
		snap.hash = m_Hash;

		snap.health = m_Health;
		snap.crew_occupied = m_Passengers;
		snap.speed_ms = m_SpeedMS;
		snap.engine_on = m_EngineOn;
		snap.fuel_fraction = m_FuelFraction;

		snap.pos_x = m_PosX;
		snap.pos_y = m_PosZ;
		snap.pos_alt = m_PosY;
		snap.pos_angle = m_Yaw;
	}

	/**
	    \brief Cache identity fields from the transport.
	*/
	protected void BuildIdentity()
	{
		if (!m_Transport || m_ClassName != string.Empty)
			return;

		m_ClassName = m_Transport.GetType();
		m_ClassName.TrimInPlace();

		m_VehicleType = m_Transport.GetVehicleType();
		m_VehicleType.TrimInPlace();
		m_VehicleType.Replace("VehicleType", "");

		m_Hash = PersistentHash(m_Transport);
	}

	/**
	    \brief Return number of occupied crew seats.
	*/
	private int GetPassengersCount()
	{
		if (!m_Transport)
			return 0;

		int size = m_Transport.CrewSize();
		int passengers = 0;
		for (int i = 0; i < size; ++i) {
			if (m_Transport.CrewMember(i))
				passengers++;
		}

		return passengers;
	}

	/**
	    \brief Generate a stable hash int from an entity's persistent ID.
	    \details Combines four persistent-ID ints and hashes them. The sign bit is
	             masked off so the result is always non-negative; -1 is returned
	             only when the hash is unavailable (no entity or empty ID).
	    \param entity EntityAI with persistence
	    \return int non-negative hash in [0, 2^31-1], or -1 if unavailable
	*/
	static int PersistentHash(EntityAI entity)
	{
		if (!entity)
			return -1;

		int p1, p2, p3, p4;
		entity.GetPersistentID(p1, p2, p3, p4);

		if (p1 == 0 && p2 == 0 && p3 == 0 && p4 == 0)
			return -1;

		return string.Format("%1_%2_%3_%4", p1, p2, p3, p4).Hash() & 0x7FFFFFFF;
	}
}
#endif
