/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Per-Transport data collector.
    \details Stores plain fields updated during scrape. Fills snapshot on demand.
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
	int m_EngineOn;
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
			    StateZ_Config.Get().settings.collect_interval_sec * 1000,
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

		m_EngineOn = 0;
		m_FuelFraction = -1.0;

		CarScript car;
		if (Class.CastTo(car, m_Transport) && car.IsVitalFuelTank()) {
			m_EngineOn = StateZ_LabelUtils.Bool(car.EngineIsOn());
			m_FuelFraction = car.GetFluidFraction(CarFluid.FUEL);
			return;
		}

		BoatScript boat;
		if (Class.CastTo(boat, m_Transport) && boat.HasEngine()) {
			if (!boat.EngineIsOn() && speed < 0.30)
				m_SpeedMS = 0;

			m_EngineOn = StateZ_LabelUtils.Bool(boat.EngineIsOn());
			m_FuelFraction = boat.GetFluidFraction(BoatFluid.FUEL);
			return;
		}

#ifdef EXPANSIONMODVEHICLE
		ExpansionVehicleBase expVehicle;
		if (Class.CastTo(expVehicle, m_Transport)) {
			m_EngineOn = StateZ_LabelUtils.Bool(expVehicle.Expansion_EngineIsOn());
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

		m_Hash = StateZ_LabelUtils.PersistentHash(m_Transport);
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
}
#endif
