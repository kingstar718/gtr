package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func isCoordinateFormat(input string) bool {
	coordPatterns := []string{
		`^-?\d+\.?\d*[,|]\s*-?\d+\.?\d*$`,
		`^-?\d+\.?\d*\s+-?\d+\.?\d*$`,
	}

	for _, pattern := range coordPatterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return true
		}
	}
	return false
}

func handleCoordinateConvert(input string) error {
	var longitude, latitude float64
	var originCoordinate = coordinateGcj02

	parts, err := coordinateStringToFloat(input)
	if err != nil {
		return err
	}
	longitude = parts[0]
	latitude = parts[1]

	gpsUtil := GPSUtil{}
	var bd09 []float64
	var wgs84 []float64
	var gcj02 []float64

	if originCoordinate == coordinateGcj02 {
		bd09 = gpsUtil.gcj02_To_Bd09(latitude, longitude)
		wgs84 = gpsUtil.GCJ02_To_WGS84(latitude, longitude)
		gcj02 = []float64{latitude, longitude}
	} else if originCoordinate == coordinateWgs84 {
		bd09 = gpsUtil.WGS84_To_bd09(latitude, longitude)
		gcj02 = gpsUtil.WGS84_To_Gcj02(latitude, longitude)
		wgs84 = []float64{latitude, longitude}
	} else {
		gcj02 = gpsUtil.bd09_To_Gcj02(latitude, longitude)
		wgs84 = gpsUtil.bd09_To_WGS84(latitude, longitude)
		bd09 = []float64{latitude, longitude}
	}

	fmt.Println("--------|------------------------------------")
	fmt.Printf("input  : %v,%v\n", longitude, latitude)
	fmt.Printf("coord  : %s\n", originCoordinate)
	fmt.Println("--------|------------------------------------")
	fmt.Printf("gcj02-6: %v,%v\n", gpsUtil.retain6(gcj02[1]), gpsUtil.retain6(gcj02[0]))
	fmt.Printf("gcj02  : %v,%v\n", gcj02[1], gcj02[0])
	fmt.Println("--------|------------------------------------")
	fmt.Printf("wgs84-6: %v,%v\n", gpsUtil.retain6(wgs84[1]), gpsUtil.retain6(wgs84[0]))
	fmt.Printf("wgs84  : %v,%v\n", wgs84[1], wgs84[0])
	fmt.Println("--------|------------------------------------")
	fmt.Printf("bd09-6 : %v,%v\n", gpsUtil.retain6(bd09[1]), gpsUtil.retain6(bd09[0]))
	fmt.Printf("bd09   : %v,%v\n", bd09[1], bd09[0])
	fmt.Println("--------|------------------------------------")
	return nil
}

// GPSUtil provides utilities for GPS coordinate system conversions.
// It supports conversions between WGS84, GCJ02 (Chinese coordinate system),
// and BD09 (Baidu coordinate system).
type GPSUtil struct {
}

const (
	pi  = 3.1415926535897932384626             // 圆周率
	xPi = 3.14159265358979324 * 3000.0 / 180.0 // 圆周率对应的经纬度偏移
	a   = 6378245.0                            // 长半轴
	ee  = 0.00669342162296594323               // 扁率
)

func (receiver *GPSUtil) transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*pi) + 20.0*math.Sin(2.0*x*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*pi) + 40.0*math.Sin(y/3.0*pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*pi) + 320*math.Sin(y*pi/30.0)) * 2.0 / 3.0
	return ret
}

func (receiver *GPSUtil) transformlng(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*pi) + 20.0*math.Sin(2.0*x*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*pi) + 40.0*math.Sin(x/3.0*pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*pi) + 300.0*math.Sin(x/30.0*pi)) * 2.0 / 3.0
	return ret
}

func (receiver *GPSUtil) outOfChina(lat, lng float64) bool {
	if lng < 72.004 || lng > 137.8347 {
		return true
	}
	if lat < 0.8293 || lat > 55.8271 {
		return true
	}
	return false
}

func (receiver *GPSUtil) transform(lat, lng float64) []float64 {
	if receiver.outOfChina(lat, lng) {
		return []float64{lat, lng}
	}
	dLat := receiver.transformLat(lng-105.0, lat-35.0)
	dlng := receiver.transformlng(lng-105.0, lat-35.0)
	radLat := lat / 180.0 * pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	SqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * SqrtMagic) * pi)
	dlng = (dlng * 180.0) / (a / SqrtMagic * math.Cos(radLat) * pi)
	mgLat := lat + dLat
	mglng := lng + dlng
	return []float64{mgLat, mglng}
}

// WGS84_To_Gcj02 converts WGS84 coordinates to GCJ02 (Chinese coordinate system).
func (receiver *GPSUtil) WGS84_To_Gcj02(lat, lng float64) []float64 {
	if receiver.outOfChina(lat, lng) {
		return []float64{lat, lng}
	}
	dLat := receiver.transformLat(lng-105.0, lat-35.0)
	dlng := receiver.transformlng(lng-105.0, lat-35.0)
	radLat := lat / 180.0 * pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	SqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * SqrtMagic) * pi)
	dlng = (dlng * 180.0) / (a / SqrtMagic * math.Cos(radLat) * pi)
	mgLat := lat + dLat
	mglng := lng + dlng
	return []float64{mgLat, mglng}
}

// GCJ02_To_WGS84 converts GCJ02 coordinates to WGS84.
func (receiver *GPSUtil) GCJ02_To_WGS84(lat, lng float64) []float64 {
	gps := receiver.transform(lat, lng)
	lngtitude := lng*2 - gps[1]
	latitude := lat*2 - gps[0]
	return []float64{latitude, lngtitude}
}

// gcj02_To_Bd09 converts GCJ02 coordinates to BD09 (Baidu coordinate system).
func (receiver *GPSUtil) gcj02_To_Bd09(lat, lng float64) []float64 {
	x := lng
	y := lat
	z := math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*xPi)
	theta := math.Atan2(y, x) + 0.000003*math.Cos(x*xPi)
	templng := z*math.Cos(theta) + 0.0065
	tempLat := z*math.Sin(theta) + 0.006
	gps := []float64{tempLat, templng}
	return gps
}

// bd09_To_Gcj02 converts BD09 coordinates to GCJ02.
func (receiver *GPSUtil) bd09_To_Gcj02(lat, lng float64) []float64 {
	x := lng - 0.0065
	y := lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*xPi)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*xPi)
	templng := z * math.Cos(theta)
	tempLat := z * math.Sin(theta)
	gps := []float64{tempLat, templng}
	return gps
}

// WGS84_To_bd09 converts WGS84 coordinates to BD09.
func (receiver *GPSUtil) WGS84_To_bd09(lat, lng float64) []float64 {
	gcj02 := receiver.WGS84_To_Gcj02(lat, lng)
	bd09 := receiver.gcj02_To_Bd09(gcj02[0], gcj02[1])
	return bd09
}

// bd09_To_WGS84 converts BD09 coordinates to WGS84.
func (receiver *GPSUtil) bd09_To_WGS84(lat, lng float64) []float64 {
	gcj02 := receiver.bd09_To_Gcj02(lat, lng)
	WGS84 := receiver.GCJ02_To_WGS84(gcj02[0], gcj02[1])
	//保留小数点后六位
	//WGS84[0] = receiver.retain6(WGS84[0])
	//WGS84[1] = receiver.retain6(WGS84[1])
	return WGS84
}

// retain6 rounds a float64 number to 6 decimal places.
func (receiver *GPSUtil) retain6(num float64) float64 {
	value, _ := strconv.ParseFloat(strconv.FormatFloat(num, 'f', 6, 64), 64)
	return value
}

const (
	coordinateBd    = "bd"
	coordinateBd09  = "bd09"
	coordinateGcj   = "gcj"
	coordinateGcj02 = "gcj02"
	coordinateWgs   = "wgs"
	coordinateWgs84 = "wgs84"
	coordinateGd    = "gd"
	coordinateGg    = "gg"
)

var coordinateMap = map[string]string{
	coordinateBd:   coordinateBd09,
	coordinateBd09: coordinateBd09,

	coordinateGcj:   coordinateGcj02,
	coordinateGcj02: coordinateGcj02,
	coordinateGd:    coordinateGcj02,

	coordinateWgs:   coordinateWgs84,
	coordinateWgs84: coordinateWgs84,
	coordinateGg:    coordinateWgs84,
}

// NewCoordinateCommand creates a new coordinate conversion command for the CLI.
func NewCoordinateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "coordinate",
		Aliases: []string{"coor", "-coor", "--coor", "c", "-c", "--c"},
		Short: "\n-----------------------------------\n" +
			"| COMMAND: coordinate             |\n" +
			"| TYPE: Coordinate Convert        |\n" +
			"| INPUT:                          |\n" +
			"|   1. [type] longitude,latitude  |\n" +
			"|   2. [type] longitude|latitude  |\n" +
			"|   3. [type] longitude latitude  |\n" +
			"| EXAMPLES:                       |\n" +
			"|   1. 113.901495,22.499501       |\n" +
			"|   2. 113.901495 22.499501       |\n" +
			"|   3. gcj 113.901495,22.499501   |\n" +
			"|   4. wgs 113.901495 22.499501   |\n" +
			"-----------------------------------\n",
		//Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var longitude float64
			var latitude float64
			originCoordinate := coordinateGcj02
			var argLen = len(args)

			// 1. 经度,纬度
			args0 := args[0]

			if argLen == 1 {
				floats, err := coordinateStringToFloat(args0)
				if err != nil {
					return err
				}
				longitude = floats[0]
				latitude = floats[1]
			} else if argLen == 2 {
				value, existed := coordinateMap[args0]
				// 3. 坐标系 经度,纬度
				if existed {
					originCoordinate = value
					floats, err := coordinateStringToFloat(args[1])
					if err != nil {
						return err
					}
					longitude = floats[0]
					latitude = floats[1]
				} else {
					// 2. 经度 纬度
					floats, err := coordinateStringToFloat(args0 + "|" + args[1])
					if err != nil {
						return err
					}
					longitude = floats[0]
					latitude = floats[1]
				}
				// 4. 坐标系 经度 纬度
			} else if argLen == 3 {
				value, existed := coordinateMap[args0]
				floats, err := coordinateStringToFloat(args[1] + "|" + args[2])
				if err != nil {
					return err
				}
				longitude = floats[0]
				latitude = floats[1]
				if existed {
					originCoordinate = value
				}
			}

			s := coordinateMap[originCoordinate]
			gpsUtil := GPSUtil{}
			var bd09 []float64
			var wgs84 []float64
			var gcj02 []float64

			if s == coordinateGcj02 {
				bd09 = gpsUtil.gcj02_To_Bd09(latitude, longitude)
				wgs84 = gpsUtil.GCJ02_To_WGS84(latitude, longitude)
				gcj02 = []float64{latitude, longitude}
			} else if s == coordinateWgs84 {
				bd09 = gpsUtil.WGS84_To_bd09(latitude, longitude)
				gcj02 = gpsUtil.WGS84_To_Gcj02(latitude, longitude)
				wgs84 = []float64{latitude, longitude}
			} else {
				gcj02 = gpsUtil.bd09_To_Gcj02(latitude, longitude)
				wgs84 = gpsUtil.bd09_To_WGS84(latitude, longitude)
				bd09 = []float64{latitude, longitude}
			}

			fmt.Println("--------|------------------------------------")
			fmt.Printf("input  : %v,%v\n", longitude, latitude)
			fmt.Printf("coord  : %s\n", originCoordinate)
			fmt.Println("--------|------------------------------------")
			fmt.Printf("gcj02-6: %v,%v\n", gpsUtil.retain6(gcj02[1]), gpsUtil.retain6(gcj02[0]))
			fmt.Printf("gcj02  : %v,%v\n", gcj02[1], gcj02[0])
			fmt.Println("--------|------------------------------------")
			fmt.Printf("wgs84-6: %v,%v\n", gpsUtil.retain6(wgs84[1]), gpsUtil.retain6(wgs84[0]))
			fmt.Printf("wgs84  : %v,%v\n", wgs84[1], wgs84[0])
			fmt.Println("--------|------------------------------------")
			fmt.Printf("bd09-6 : %v,%v\n", gpsUtil.retain6(bd09[1]), gpsUtil.retain6(bd09[0]))
			fmt.Printf("bd09   : %v,%v\n", bd09[1], bd09[0])
			fmt.Println("--------|------------------------------------")
			return nil
		},
	}
	return cmd
}

func coordinateStringToFloat(s string) ([]float64, error) {
	var longitude float64
	var latitude float64
	var err error

	if strings.Contains(s, "|") {
		split := strings.Split(s, "|")
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid coordinate format: expected 2 parts separated by |, got %d", len(split))
		}
		longitude, err = strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid longitude: %w", err)
		}
		latitude, err = strconv.ParseFloat(strings.TrimSpace(split[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid latitude: %w", err)
		}
		return []float64{longitude, latitude}, nil
	}

	if strings.Contains(s, ",") {
		split := strings.Split(s, ",")
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid coordinate format: expected 2 parts separated by comma, got %d", len(split))
		}
		longitude, err = strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid longitude: %w", err)
		}
		latitude, err = strconv.ParseFloat(strings.TrimSpace(split[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid latitude: %w", err)
		}
		return []float64{longitude, latitude}, nil
	}

	return nil, fmt.Errorf("invalid coordinate format: must contain comma or pipe separator")
}
