package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "timetrack",
	Short: "timetrack parses a file and looks for lines like `t=3m` and aggregates the total time",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		duration, err := aggregateDurationsInFile(filePath)
		if err != nil {
			fmt.Println("error occurred", err)
			os.Exit(1)
		}

		fmt.Println("Total Duration:", duration.String())
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func aggregateDurationsInFile(filePath string) (time.Duration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}

	fileStr := string(fileBytes)

	return aggregateDurationsInString(fileStr)
}

func aggregateDurationsInString(str string) (time.Duration, error) {
	matches := getParams(`t(?P<duration>=[\dm|s|h]*)`, str)

	var totalDuration time.Duration

	for _, mm := range matches {
		for _, m := range mm {
			m = strings.Replace(m, "=", "", 1)
			duration, err := time.ParseDuration(m)
			if err != nil {
				return 0, err
			}

			totalDuration = totalDuration + duration
		}
	}

	return totalDuration, nil
}

func getParams(regEx, s string) map[string][]string {

	var compRegEx = regexp.MustCompile(regEx)
	matches := compRegEx.FindAllStringSubmatch(s, -1)

	paramsMap := make(map[string][]string)
	for i, name := range compRegEx.SubexpNames() {
		for _, match := range matches {
			if i > 0 && i <= len(match) {
				paramsMap[name] = append(paramsMap[name], match[i])
			}
		}
	}

	return paramsMap
}
