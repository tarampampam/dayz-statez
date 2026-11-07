/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
modded class ActionMineBush
{
	/**
	    \brief Count mined bush on action completion.
	*/
	override void OnFinishProgress(ActionData action_data)
	{
		StateZ_Storage.s_MinedBushes += 1;

		super.OnFinishProgress(action_data);
	}
}
#endif
