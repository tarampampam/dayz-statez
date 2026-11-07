/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/
#ifdef SERVER
/**
    \brief Updates StateZ counters on spawn/cleanup.
*/
modded class Edible_Base
{
	/**
	    \brief Increment food gauge on init.
	*/
	override void EEInit()
	{
		super.EEInit();

		// because food metrics init later, need check is storage inited
		if (!StateZ_Storage.IsInitialized())
			return;

		if (!IsInherited(Bottle_Base))
			StateZ_Storage.s_Food += 1;

		StateZ_Storage.FoodMetricChange(StateZ_GetFoodType(), true);
	}

	/**
	    \brief Decrement food gauge on delete.
	*/
	override void EEDelete(EntityAI parent)
	{
		if (StateZ_Storage.IsInitialized()) {
			if (!IsInherited(Bottle_Base))
				StateZ_Storage.s_Food -= 1;

			StateZ_Storage.FoodMetricChange(StateZ_GetFoodType(), false);
		}

		super.EEDelete(parent);
	}

	/**
	    \brief Classify edible item into StateZ_FoodTypes based on inheritance and config properties.
	*/
	StateZ_FoodTypes StateZ_GetFoodType()
	{
		// skip bottle base by exists "bottles" metric
		if (IsInherited(Bottle_Base))
			return StateZ_FoodTypes.NONE;

		if (IsFruit())
			return StateZ_FoodTypes.FRUIT;

		if (IsMushroom())
			return StateZ_FoodTypes.MUSHROOM;

		if (IsCorpse())
			return StateZ_FoodTypes.CORPSE;

		if (IsInherited(Worm))
			return StateZ_FoodTypes.WORM;

		if (IsInherited(Guts) || IsInherited(SmallGuts))
			return StateZ_FoodTypes.GUTS;

		if (IsMeat())
			return StateZ_FoodTypes.MEAT;

		if (IsInherited(HumanSteakMeat))
			return StateZ_FoodTypes.HUMAN_MEAT;

		if (GetLiquidType() == LIQUID_DISINFECTANT || GetLiquidTypeInit() == LIQUID_DISINFECTANT)
			return StateZ_FoodTypes.DISINFECTANT;

		if (ConfigGetString("stackedUnit") == "pills")
			return StateZ_FoodTypes.PILLS;

		if (IsInherited(SodaCan_ColorBase))
			return StateZ_FoodTypes.DRINK;

		if (IsInherited(Snack_ColorBase) || IsInherited(Zagorky_ColorBase))
			return StateZ_FoodTypes.SNACK;

		if (IsInherited(Candycane_Colorbase))
			return StateZ_FoodTypes.CANDY;

		if (IsKindOf("FoodCan_100g_ColorBase"))
			return StateZ_FoodTypes.CANNED_SMALL;

		if (IsKindOf("FoodCan_250g_ColorBase"))
			return StateZ_FoodTypes.CANNED_MEDIUM;

		if (ConfigGetString("soundImpactType") == "metal")
			return StateZ_FoodTypes.CANNED_BIG;

		if (ConfigGetString("soundImpactType") == "glass")
			return StateZ_FoodTypes.JAR;

		return StateZ_FoodTypes.OTHER;
	}
}

/**
    \brief Canonical food categories for StateZ.
    \details
      - Mapped from Edible_Base runtime type to stable labeled groups.
      - Order is fixed; values are static and used as map keys in StateZ_Storage.
*/
enum StateZ_FoodTypes {
	NONE = -1, // skip (not counted or tracked elsewhere)
	OTHER = 0, // fallback category
	FRUIT,
	MUSHROOM,
	CORPSE,
	WORM,
	GUTS,
	MEAT,
	HUMAN_MEAT,
	DISINFECTANT,
	PILLS,
	DRINK,
	SNACK,
	CANDY,
	CANNED_SMALL,
	CANNED_MEDIUM,
	CANNED_BIG,
	JAR,
}
#endif
