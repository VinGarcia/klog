package klog

// MergeMaps will merge into the baseMap
// all the maps from the `maps` slice.
func MergeMaps(baseMap *Body, maps ...Body) {
	if *baseMap == nil {
		*baseMap = Body{}
	}

	for _, m := range maps {
		for k, v := range m {
			(*baseMap)[k] = v
		}
	}
}
