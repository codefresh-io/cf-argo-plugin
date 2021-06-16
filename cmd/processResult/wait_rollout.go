package processResult

import (
	"fmt"
	"github.com/spf13/cobra"
)

var waitRolloutArgsOptions struct {
	PipelineId string
	BuildId    string
}

var WaitRolloutCmd = &cobra.Command{
	Use:   "wait-rollout [app]",
	Short: "Wait rollout",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			fmt.Println("Wrong amount of arguments")
			return nil
		}
		return GetWaitRolloutHandler().Handle(args[0])
	},
}

func init() {
	f := WaitRolloutCmd.Flags()
	f.StringVar(&waitRolloutArgsOptions.PipelineId, "pipeline-id", "", "Pipeline id where argo sync was executed")
	f.StringVar(&waitRolloutArgsOptions.BuildId, "build-id", "", "Build id where argo sync was executed")
}
