package start

import (
	"fmt"
	"github.com/pefish/pmon2/app"
	"github.com/pefish/pmon2/app/model"
	"github.com/pefish/pmon2/app/output"
	"github.com/pefish/pmon2/app/svc/process"
	"github.com/spf13/cobra"
	"os"
)

var flag model.ExecFlags

var Cmd = &cobra.Command{
	Use:   "start",
	Short: "start some process by id or name",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args, flag.Json())
	},
}

func init() {
	Cmd.Flags().StringVarP(&flag.LogDir, "log_dir", "d", "", "the process stdout log dir")
	Cmd.Flags().StringVarP(&flag.Log, "log", "l", "", "the process stdout log")
}

func cmdRun(args []string, flags string) {
	if len(args) == 0 {
		app.Log.Fatal("please input start process id or name")
	}

	val := args[0]
	var m model.Process
	if err := app.Db().First(&m, "id = ? or name = ?", val, val).Error; err != nil {
		app.Log.Fatal(fmt.Sprintf("the process %s not exist", val))
	}

	// checkout process state
	if process.IsRunning(m.Pid) {
		if m.Status != model.StatusRunning {
			m.Status = model.StatusRunning
			app.Db().Save(&m)
		}
		output.TableOne(m.RenderTable())
		return
	}

	rel, err := process.TryStart(m, flags)
	if err != nil {
		if len(os.Getenv("PMON2_DEBUG")) > 0 {
			app.Log.Fatalf("%+v", err)
		} else {
			app.Log.Fatal(err.Error())
		}
	}

	output.TableOne(rel)
}
