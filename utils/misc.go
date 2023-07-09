package utils

func NumTo3xTB(base int64) int64 {
	base *= 1000 /** KB */
	base *= 1000 /** MB */
	base *= 1000 /** GB */
	base *= 1000 /** TB */
	base *= 3    /** Redundancy */
	return base
}
