package main

func getAverageLightlevel(lightdata []LightData) int {
	if len(lightdata) == 0 {
		return 0
	}
	total := 0
	for _, lightdatum := range lightdata {
		total += lightdatum.LightLevel
	}
	return total / len(lightdata)
}

func getIntensityFromLightlevel(lightlevel int) int {
	if lightlevel > 100 {
		lightlevel = 100
	}
	return 100 - lightlevel
}
