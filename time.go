package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

type TimeStruct struct {
	OriginArg   string
	ParseTime   time.Time
	Timestamp10 int64
	Timestamp13 int64
	TimeFormat1 string
	TimeFormat2 string
}

// NewTimeCommand 创建一个新的时间转换命令
func NewTimeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "time",
		Aliases: []string{"t", "-t", "--t", "-time", "--time"},
		Short: "\n-------------------------------\n" +
			"| COMMAND: time               |\n" +
			"| TYPE: Time Format Convert   |\n" +
			"| INPUT:                      |\n" +
			"|   1. 10 timestamp           |\n" +
			"|   2. 13 timestamp           |\n" +
			"|   3. standard date          |\n" +
			"| EXAMPLES:                   |\n" +
			"|   1. 1727087511             |\n" +
			"|   2. 1727087511000          |\n" +
			"|   3. 2024-09-23 10:31:51    |\n" +
			"|   4. 20240923103151         |\n" +
			"-------------------------------\n",
		//Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var originArg string
			if len(args) == 1 {
				originArg = args[0]
			} else if len(args) == 2 {
				originArg = args[0] + " " + args[1]
			}

			timeFormat1 := "2006-01-02 15:04:05"
			timeFormat2 := "20060102150405"
			timeFormat1Len := len(timeFormat1)
			timeFormat2Len := len(timeFormat2)

			timeStruct := TimeStruct{}
			timeStruct.OriginArg = originArg

			l := len(originArg)
			// 统一转time.Time
			var parseTime time.Time
			if l == 10 {
				// 10位时间戳
				parseInt, _ := strconv.ParseInt(originArg, 10, 64)
				parseTime = time.Unix(parseInt, 0)
			} else if l == 13 {
				// 13位时间戳
				parseInt, _ := strconv.ParseInt(originArg, 10, 64)
				parseTime = time.UnixMilli(parseInt)
			} else if l == timeFormat1Len {
				// 格式为：20060102150405
				parseTime, _ = time.Parse(timeFormat1, originArg)
			} else if l == timeFormat2Len {
				// 格式为：2006-01-02 15:04:05
				parseTime, _ = time.Parse(timeFormat2, originArg)
			} else {
				parseTime = time.Now()
			}

			timeStruct.ParseTime = parseTime

			timeStruct.TimeFormat1 = parseTime.Format(timeFormat1)
			timeStruct.TimeFormat2 = parseTime.Format(timeFormat2)
			timeStruct.Timestamp13 = parseTime.UnixNano() / int64(time.Millisecond)
			timeStruct.Timestamp10 = parseTime.Unix()

			fmt.Printf("imput      : %s\n", timeStruct.OriginArg)
			fmt.Printf("timestamp1 : %d\n", timeStruct.Timestamp10)
			fmt.Printf("timestamp2 : %d\n", timeStruct.Timestamp13)
			fmt.Printf("format1    : %s\n", timeStruct.TimeFormat1)
			fmt.Printf("format2    : %s\n", timeStruct.TimeFormat2)

			return nil
		},
	}
	return cmd
}
