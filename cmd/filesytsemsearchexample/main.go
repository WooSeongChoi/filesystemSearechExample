package main

import (
	"filesystemSearchExample/internal/search"
	"flag"
	"fmt"
	"github.com/spf13/afero"
	"math"
)

type ParsedFlag struct {
	RootDir string
	Unit    string
}

var Unit2Scale = map[string]float64{
	"TB": math.Pow(1024, 4),
	"GB": math.Pow(1024, 3),
}

func ParseFlags() *ParsedFlag {
	rootDir := flag.String("rootDir", "", "Search root Directory")
	unit := flag.String("unit", "GB", "Search unit")
	flag.Parse()

	if *rootDir == "" {
		panic("You must specify -rootDir")
	}

	return &ParsedFlag{*rootDir, *unit}
}

func main() {
	parsed := ParseFlags()
	fs := afero.NewOsFs()

	dirInfos, err := afero.ReadDir(fs, parsed.RootDir)
	if err != nil {
		panic(err)
	}

	rootPaths := make([]string, 0, len(dirInfos))
	for _, dirInfo := range dirInfos {
		rootPaths = append(rootPaths, fmt.Sprintf(
			"%s/%s",
			parsed.RootDir,
			dirInfo.Name()))
	}

	summaryMap := search.GetInspectionSummary(fs, &rootPaths)
	PrintSummary(summaryMap, Unit2Scale[parsed.Unit], parsed.Unit)

	var Sum uint64
	for _, val := range *summaryMap {
		Sum += val
	}
	fmt.Println(float64(Sum) / Unit2Scale[parsed.Unit])

}

func PrintSummary(summaryMap *map[string]uint64, scale float64, unit string) {
	for key, val := range *summaryMap {
		fmt.Printf("%s: %.2f %s\n", key, float64(val)/scale, unit)
	}
}
