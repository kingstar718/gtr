package main

import "testing"

func TestName(t *testing.T) {
	gps := GPSUtil{}
	wgs84 := gps.GCJ02_To_WGS84(38.65638297231525, 116.50661644375265)
	gaode := gps.WGS84_To_Gcj02(wgs84[0], wgs84[1])
	t.Log(wgs84)
	t.Log(gaode)
}
