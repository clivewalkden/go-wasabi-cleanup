package reporting

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/viper"
	"math"
)

type Result struct {
	Name        string
	Kept        int
	KeptSize    string
	Deleted     int
	DeletedSize string
}

type Report struct {
	Result []Result
	DryRun bool
}

func Output(report Report) {
	if viper.GetBool("verbose") == false {
		fmt.Print("\033[H\033[2J") //clear screen
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	heading := color.New(color.FgBlue, color.Bold).PrintlnFunc()
	if report.DryRun {
		heading("Potential Results: (DryRun Mode)")
	} else {
		heading("Results:")
	}
	fmt.Println("")

	tbl := table.New("Name", "Deleted Files", "Deleted Size", "Deleted (%)", "Remaining Files", "Remaining Size")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt).WithPadding(5)

	for _, element := range report.Result {
		tbl.AddRow(element.Name, element.Deleted, element.DeletedSize, deletedPerc(element), element.Kept, element.KeptSize)
	}
	tbl.Print()
	fmt.Println("")
	fmt.Println("")
}

func deletedPerc(result Result) (delta float64) {
	originalNumber := float64(result.Kept + result.Deleted)
	decrease := originalNumber - float64(result.Kept)
	delta = (decrease / originalNumber) * 100
	return math.Round(delta*100) / 100
}
