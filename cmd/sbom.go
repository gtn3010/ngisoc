package cmd

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/spf13/cobra"

	"github.com/gtn3010/ngisoc/pkg/sbom"
	"github.com/gtn3010/ngisoc/pkg/util"
)

func NewSbom(options *[]crane.Option) *cobra.Command {
	platform := &util.PlatformValue{}
	filter := ""
	sbomCmd := &cobra.Command{
		Use:   "sbom",
		Short: "Get sbom of container image uploaded by buildkit/cosign",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var supportedKinds = []string{"cosign", "buildkit"}
			imgName := args[0]
			kinds, err := util.FilterArgs(filter, supportedKinds)
			if err != nil {
				return cmd.Usage()
			}

			// manifest, err := crane.Get(imgName, *options...)
			// if err != nil {
			// 	return fmt.Errorf("fetching manifest %s: %w", imgName, err)
			// }

			// digest := manifest.Descriptor.Digest.String()
			if len(kinds) > 1 {
				output, err := sbom.GetBuildkit(imgName, platform, options)
				if err != nil {
					fmt.Println(output)
				}
			} else {
				output, err := sbom.GetBuildkit(imgName, platform, options)
				if err == nil {
					fmt.Println(output)
				}

			}
			return nil
		},
	}
	sbomCmd.PersistentFlags().StringVar(&filter, "filter", "", "Filter output by kind: cosign/buildkit. Default no filter means getting all types.")
	sbomCmd.PersistentFlags().Var(platform, "platform", "Specifies the platform in the form os/arch[/variant][:osversion] (e.g. linux/amd64). Support with buildkit only.")
	return sbomCmd
}
