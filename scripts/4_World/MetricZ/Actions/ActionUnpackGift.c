/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
modded class ActionUnpackGift
{
	/**
	    \brief Count opened gift on action completion.
	*/
	override void OnFinishProgress(ActionData action_data)
	{
		MetricZ_Storage.s_GiftsUnpacked += 1;

		super.OnFinishProgress(action_data);
	}
}
#endif
