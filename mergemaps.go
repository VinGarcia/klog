package klog

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
