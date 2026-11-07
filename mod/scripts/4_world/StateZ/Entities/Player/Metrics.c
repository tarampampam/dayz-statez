/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Per-player data collector.
    \details Stores plain fields updated during scrape. Fills snapshot on demand.
*/
class StateZ_PlayerMetrics
{
	// Spawn timestamp in seconds
	protected int m_InitTick;

	// Ping & Throttle
	protected int m_PingMinWindow = -1; //!< Min ping window.
	protected int m_PingMaxWindow = -1; //!< Max ping window.
	protected float m_ThrottleWindow = -1.0; //!< Throttle window.
	protected float m_SampleAccTime; //!< Sample accumulation time.

	// Parent player instance.
	protected PlayerBase m_Player;

	// vitals
	int m_IsLoaded;
	float m_Health;
	float m_Blood;
	float m_Shock;
	float m_Energy;
	float m_Water;
	float m_Toxicity;
	float m_TempC;
	float m_Weight;
	float m_Wetness;
	float m_LifeSeconds;

	// position
	float m_PosX;
	float m_PosY;
	float m_PosZ;
	float m_Yaw;

	// network
	int m_PingMin;
	int m_PingMax;
	float m_Throttle;

	// states
	int m_AgentsCount;
	int m_BleedingSources;
	int m_ImmunityBoosted;
	int m_Unconscious;
	int m_Restrained;
	int m_GodMode;

	// analytics
	float m_StatPlaytimeSeconds;
	float m_StatDistanceMeters;
	float m_StatLongestHitMeters;
	int m_StatPlayersKilled;
	int m_StatInfectedKilled;

	// identity (read once identity becomes available)
	protected string m_SteamId;
	protected string m_Guid;
	protected string m_Name;
	protected string m_BloodType;
	protected string m_Gender;
	protected string m_CharacterType;

	/**
	    \brief Constructor. Records spawn tick.
	*/
	void StateZ_PlayerMetrics()
	{
		m_InitTick = g_Game.GetTickTime();
	}

	/**
	    \brief One-time init with player reference.
	*/
	void Init(PlayerBase player)
	{
		if (!player || player.IsDamageDestroyed())
			return;

		m_Player = player;
	}

	/**
	    \brief Update all fields from the player state.
	*/
	void Update()
	{
		if (!m_Player)
			return;

		// resolve identity lazily
		if (m_SteamId == string.Empty)
			FetchIdentity();

		// vitals
		m_IsLoaded = StateZ_LabelUtils.Bool(m_Player.IsPlayerLoaded());
		m_Health = m_Player.GetHealth01("", "Health");
		m_Blood = m_Player.GetHealth01("", "Blood");
		m_Shock = m_Player.GetHealth01("", "Shock");
		m_Energy = m_Player.GetStatEnergy().Get() / m_Player.GetStatEnergy().GetMax();
		m_Water = m_Player.GetStatWater().Get() / m_Player.GetStatWater().GetMax();
		m_Toxicity = m_Player.GetStatToxicity().Get() / m_Player.GetStatToxicity().GetMax();
		m_TempC = m_Player.GetTemperature();
		m_Weight = m_Player.GetPlayerLoad();
		m_Wetness = m_Player.GetStatWet().Get();
		m_LifeSeconds = g_Game.GetTickTime() - m_InitTick;

		// position
		if (!StateZ_Config.Get().disabled_metrics.positions) {
			vector pos = m_Player.GetPosition();
			m_PosX = pos[0];
			m_PosZ = pos[2];

			if (!StateZ_Config.Get().disabled_metrics.positions_height)
				m_PosY = pos[1];

			if (!StateZ_Config.Get().disabled_metrics.positions_yaw)
				m_Yaw = m_Player.GetOrientation()[0];
		}

		// ping & throttle
		UpdateNetwork();

		// states
		m_AgentsCount = StateZ_LabelUtils.BitsCount(m_Player.GetAgents());
		m_BleedingSources = StateZ_LabelUtils.BitsCount(m_Player.GetBleedingBits());
		m_ImmunityBoosted = StateZ_LabelUtils.Bool(m_Player.m_ImmunityBoosted);
		m_Unconscious = StateZ_LabelUtils.Bool(m_Player.IsUnconscious());
		m_Restrained = StateZ_LabelUtils.Bool(m_Player.IsRestrained());
		m_GodMode = StateZ_LabelUtils.Bool(!m_Player.GetAllowDamage());

		// analytics
		m_StatPlaytimeSeconds = m_Player.StatGet(AnalyticsManagerServer.STAT_PLAYTIME);
		m_StatDistanceMeters = m_Player.StatGet(AnalyticsManagerServer.STAT_DISTANCE);
		m_StatLongestHitMeters = m_Player.StatGet(AnalyticsManagerServer.STAT_LONGEST_SURVIVOR_HIT);
		m_StatPlayersKilled = m_Player.StatGet(AnalyticsManagerServer.STAT_PLAYERS_KILLED);
		m_StatInfectedKilled = m_Player.StatGet(AnalyticsManagerServer.STAT_INFECTED_KILLED);
	}

	/**
	    \brief Per-frame hook from PlayerBase to sample network stats.
	    \details Called at most 1 Hz per player to accumulate min/max over scrape window.
	*/
	void SampleNetwork()
	{
		if (!m_Player)
			return;

		m_SampleAccTime += m_Player.GetDeltaT();
		if (m_SampleAccTime < 1.0) // sample ~1 Hz
			return;

		m_SampleAccTime = 0;

		PlayerIdentity idp = m_Player.GetIdentity();
		if (!idp)
			return;

		int pMin = idp.GetPingMin();
		int pMax = idp.GetPingMax();
		float throttle = idp.GetOutputThrottle();

		if (pMin > 0 && (m_PingMinWindow < 0 || pMin < m_PingMinWindow))
			m_PingMinWindow = pMin;

		if (pMax > 0 && (m_PingMaxWindow < 0 || pMax > m_PingMaxWindow))
			m_PingMaxWindow = pMax;

		if (throttle > 0 && (m_ThrottleWindow < 0 || throttle > m_ThrottleWindow))
			m_ThrottleWindow = throttle;
	}

	/**
	    \brief Fill a player snapshot with current field values.
	    \param snap Snapshot instance to populate.
	*/
	void ToSnapshot(StateZ_PlayerSnapshot snap)
	{
		if (!snap)
			return;

		snap.steam_id = m_SteamId;
		snap.guid = m_Guid;
		snap.name = m_Name;
		snap.blood_type = m_BloodType;
		snap.gender = m_Gender;
		snap.character_type = m_CharacterType;

		snap.is_loaded = m_IsLoaded;
		snap.health = m_Health;
		snap.blood = m_Blood;
		snap.shock = m_Shock;
		snap.energy = m_Energy;
		snap.water = m_Water;
		snap.toxicity = m_Toxicity;
		snap.temperature_celsius = m_TempC;
		snap.weight_grams = m_Weight;
		snap.wetness = m_Wetness;
		snap.lifetime_seconds = m_LifeSeconds;

		snap.pos_x = m_PosX;
		snap.pos_z = m_PosZ;
		snap.pos_y = m_PosY;
		snap.orientation = m_Yaw;

		snap.ping_min = m_PingMin;
		snap.ping_max = m_PingMax;
		snap.throttle = m_Throttle;

		snap.agents_active = m_AgentsCount;
		snap.bleeding_sources = m_BleedingSources;
		snap.immunity_boosted = m_ImmunityBoosted;
		snap.is_unconscious = m_Unconscious;
		snap.is_restrained = m_Restrained;
		snap.godmode = m_GodMode;

		snap.stat_playtime_seconds = m_StatPlaytimeSeconds;
		snap.stat_distance_meters = m_StatDistanceMeters;
		snap.stat_longest_hit_meters = m_StatLongestHitMeters;
		snap.stat_players_killed = m_StatPlayersKilled;
		snap.stat_infected_killed = m_StatInfectedKilled;
	}

	/**
	    \brief Read identity fields from the player (called lazily).
	*/
	protected void FetchIdentity()
	{
		if (!m_Player)
			return;

		PlayerIdentity identity = m_Player.GetIdentity();
		if (!identity)
			return;

		m_SteamId = identity.GetPlainId();
		m_Guid = identity.GetId();
		m_Name = identity.GetName();

		bool _p;
		string _t;
		m_BloodType = BloodTypes.GetBloodTypeName(m_Player.GetStatBloodType().Get(), _t, _p);

		if (m_Player.IsMale()) {
			m_Gender = "male";
			m_CharacterType = m_Player.GetType();
			m_CharacterType.Replace("SurvivorM_", "");
		} else {
			m_Gender = "female";
			m_CharacterType = m_Player.GetType();
			m_CharacterType.Replace("SurvivorF_", "");
		}

		m_CharacterType.TrimInPlace();
	}

	/**
	    \brief Flush sampled network window into fields.
	*/
	protected void UpdateNetwork()
	{
		if (!m_Player)
			return;

		int pMin = 0;
		int pMax = 0;
		float throttle = 0.0;

		if (m_PingMinWindow >= 0)
			pMin = m_PingMinWindow;
		if (m_PingMaxWindow >= 0)
			pMax = m_PingMaxWindow;
		if (m_ThrottleWindow >= 0.0)
			throttle = m_ThrottleWindow;

		// fallback to current identity values
		if (pMin <= 0 && pMax <= 0) {
			PlayerIdentity idp = m_Player.GetIdentity();
			if (idp) {
				pMin = idp.GetPingMin();
				pMax = idp.GetPingMax();
				throttle = idp.GetOutputThrottle();
			}
		}

		if (pMax < pMin)
			pMax = pMin;

		m_PingMin = pMin;
		m_PingMax = pMax;
		m_Throttle = throttle;

		m_PingMinWindow = -1;
		m_PingMaxWindow = -1;
		m_ThrottleWindow = -1.0;
	}
}
#endif
