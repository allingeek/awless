package commands

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:                "sync",
	Short:              "Manage your local infrastructure",
	PersistentPreRun:   applyHooks(initAwlessEnvHook, initCloudServicesHook, initSyncerHook, checkStatsHook),
	PersistentPostRunE: saveHistoryHook,

	RunE: func(cmd *cobra.Command, args []string) error {
		region := database.MustGetDefaultRegion()

		graphs, err := sync.DefaultSyncer.Sync(aws.InfraService, aws.AccessService)
		if err != nil {
			return err
		}

		if verboseFlag {
			printWithTabs := func(res *graph.Resource, distance int) {
				var tabs bytes.Buffer
				for i := 0; i < distance; i++ {
					tabs.WriteByte('\t')
				}
				fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), res.Type(), res.Id())
			}

			root := graph.InitResource(region, graph.Region)

			graphs[aws.InfraService.Name()].VisitChildren(root, printWithTabs)
			graphs[aws.AccessService.Name()].VisitChildren(root, printWithTabs)
		}

		return nil
	},
}