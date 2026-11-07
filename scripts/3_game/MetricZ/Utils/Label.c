/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Utility helpers for MetricZ.
*/
class MetricZ_LabelUtils
{
	/**
	    \brief Generate stable hash int from entity persistent ID.
	    \details Combines four persistent-ID ints and hashes them.
	             Falls back to a pseudo-unique hash for non-persistent entities.
	    \param entity EntityAI with persistence
	    \return int hash like -123456789 or 987654321
	*/
	static int PersistentHash(EntityAI entity)
	{
		if (!entity)
			return 0;

		int p1, p2, p3, p4;
		entity.GetPersistentID(p1, p2, p3, p4);

		// Non-persistent objects get a random hash;
		// stability is guaranteed only if called once per instance for labels.
		if (p1 == 0 && p2 == 0 && p3 == 0 && p4 == 0)
			return string.Format("%1_%2", entity.GetType(), Math.RandomInt(1, int.MAX)).Hash();

		return string.Format("%1_%2_%3_%4", p1, p2, p3, p4).Hash();
	}

	/**
	    \brief Convert bool to 0/1.
	    \param x Boolean
	    \return int 1 if true else 0
	*/
	static int Bool(bool x)
	{
		if (x)
			return 1;

		return 0;
	}

	/**
	    \brief Count number of set bits in bitmask.
	    \details Uses Brian Kernighan's algorithm. Each iteration clears the lowest set bit.
	    \param mask Integer bitmask
	    \return int Number of bits set to 1
	*/
	static int BitsCount(int mask)
	{
		int count = 0;
		while (mask != 0) {
			mask &= mask - 1; // clear lowest set bit
			count++;
		}

		return count;
	}
}
#endif
