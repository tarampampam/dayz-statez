/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
modded class ActionMineTree
{
	/**
	    \brief Count mined tree on action completion.
	*/
	override void OnFinishProgress(ActionData action_data)
	{
		MetricZ_Storage.s_MinedTrees += 1;

		super.OnFinishProgress(action_data);
	}
}
#endif
