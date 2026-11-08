/*
    SPDX-License-Identifier: GPL-3.0-or-later
*/

#ifdef SERVER
/**
    \brief Numeric helpers for snapshot values.
    \details Rounding quantizes float fields before serialization: it trims
             engine float jitter so unchanged snapshots hash identically
             (fewer redundant writes via write-if-changed). Note `float` is
             32-bit, so the scaled value must stay within its exact-integer
             range (~1.6e7): Round2 is safe for map-scale magnitudes, Round4
             only for small (0..1) values.
*/
class StateZ_Num
{
	/**
	    \brief Round to 2 decimals. For positions/distances/angles (map scale).
	*/
	static float Round2(float v)
	{
		return Math.Round(v * 100.0) / 100.0;
	}

	/**
	    \brief Round to 4 decimals. For normalized 0..1 values only.
	*/
	static float Round4(float v)
	{
		return Math.Round(v * 10000.0) / 10000.0;
	}
}
#endif
