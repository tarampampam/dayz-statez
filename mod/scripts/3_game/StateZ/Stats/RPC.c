/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER

class StateZ_RpcStatsDTO
{
	void StateZ_RpcStatsDTO()
	{
		rpc_input = new map<string, int>();
	}

	ref map<string, int> rpc_input;
}

/**
    \brief RPC counters aggregator
*/
class StateZ_RpcStats
{
	protected static ref map<int, int> s_InputRPCsRegistry = new map<int, int>(); //!< Map of RPC type to count

	/**
	    \brief Increment counter for an RPC type.
	    \param rpc_type Engine RPC id.
	*/
	static void Inc(int rpc_type)
	{
		int v;
		if (s_InputRPCsRegistry.Find(rpc_type, v))
			s_InputRPCsRegistry.Set(rpc_type, v + 1);
		else
			s_InputRPCsRegistry.Insert(rpc_type, 1);
	}

	/**
	    \brief Serialize all RPC data to JSON.
	    \param json Output JSON string.
	*/
	static void BuildJson(out string json)
	{
		if (s_InputRPCsRegistry.Count() == 0) {
			json = "{}";
			return;
		}

		StateZ_RpcStatsDTO dto = new StateZ_RpcStatsDTO();

		for (int i = 0; i < s_InputRPCsRegistry.Count(); ++i) {
			int id = s_InputRPCsRegistry.GetKey(i);
			int val = s_InputRPCsRegistry.GetElement(i);
			dto.rpc_input.Insert(id.ToString(), val);
		}

		JsonSerializer().WriteToString(dto, false, json);
	}
}
#endif
