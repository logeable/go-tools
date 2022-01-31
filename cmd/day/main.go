package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "day",
	}
)

func addCmd() *cobra.Command {
	layout := "2006/01/02 15:04:05"
	var (
		begin string
	)
	cmd := &cobra.Command{
		Use:  "add",
		Args: cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			count, unit, err := parseDuration(args[0])
			if err != nil {
				fatal(fmt.Errorf("parse duration failed: %v", err))
			}
			var year, month, day int64
			var dur time.Duration
			switch unit {
			case UnitYear:
				year += count
			case UnitMonth:
				month += count
			case UnitDay:
				day += count
			case UnitHour:
				dur = time.Hour * time.Duration(count)
			case UnitMinute:
				dur = time.Minute * time.Duration(count)
			case UnitSecond:
				dur = time.Second * time.Duration(count)
			default:
				fatal(fmt.Errorf("invalid unit: %2x(%v)", unit, unit))
			}
			beginTime, err := time.Parse(layout, begin)
			if err != nil {
				fatal(fmt.Errorf("invalid begin time: %v", err))
			}
			result := beginTime.AddDate(int(year), int(month), int(day)).Add(dur)
			fmt.Println(result.Format(layout))
		},
	}

	now := time.Now().Format(layout)
	cmd.Flags().StringVarP(&begin, "begin", "b", now, "begin time")
	return cmd
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

type Unit byte

const (
	UnitYear   Unit = 'y'
	UnitMonth  Unit = 'm'
	UnitDay    Unit = 'd'
	UnitHour   Unit = 'h'
	UnitMinute Unit = 'M'
	UnitSecond Unit = 's'
)

func parseDuration(s string) (int64, Unit, error) {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return 0, 0, errors.New("invalid length")
	}
	unit := Unit(s[len(s)-1])
	count, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
	if err != nil {
		return 0, 0, errors.New("invalid count")
	}
	return count, unit, nil
}

func main() {
	rootCmd.AddCommand(addCmd())

	err := rootCmd.Execute()
	if err != nil {
		fatal(err)
	}
}
