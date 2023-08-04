package desc

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/pefish/pmon2/app"
	"github.com/pefish/pmon2/app/model"
	"github.com/pefish/pmon2/app/output"
	"github.com/spf13/cobra"
	"strconv"
)

var Cmd = &cobra.Command{
	Use:     "desc",
	Aliases: []string{"show"},
	Short:   "print the process detail message",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			app.Log.Fatalf("The command need process name or id, Example: sudo pmon2 show test")
			return
		}

		cmdRun(args)
	},
}

func cmdRun(args []string) {
	val := args[0]

	var process model.Process
	err := app.Db().Find(&process, "name = ? or id = ?", val, val).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			app.Log.Fatalf("pmon2 run err: %v", err)
		}

		// not found
		app.Log.Errorf("process %s not exist", val)
		return
	}

	rel := [][]string{
		{"status", process.Status},
		{"id", strconv.Itoa(int(process.ID))},
		{"name", process.Name},
		{"pid", strconv.Itoa(process.Pid)},
		{"process", process.ProcessFile},
		{"args", process.Args},
		{"user", process.Username},
		{"log", process.Log},
		{"no-autorestart", process.NoAutoRestartStr()},
		{"created_at", process.CreatedAt.Format("2006-01-02 15:04:05")},
		{"updated_at", process.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	output.DescTable(rel)
}
