package cmd

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/spf13/cobra"

	"github.com/gtn3010/ngisoc/pkg/signature"
	"github.com/gtn3010/ngisoc/pkg/util"
)

func NewSig(options *[]crane.Option) *cobra.Command {
	var filter string
	sigCmd := &cobra.Command{
		Use:   "signature",
		Short: "Get Container Image Signature signed by cosign/aws signer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var supportedKinds = []string{"cosign", "notation"}
			imgName := args[0]
			kinds, err := util.FilterArgs(filter, supportedKinds)
			if err != nil {
				return cmd.Usage()
			}

			digest, err := util.GetDigest(imgName, options)
			if err != nil {
				return err
			}

			if len(kinds) > 1 {
				output, err := signature.GetAll(imgName, digest, kinds, options)
				if err != nil {
					return fmt.Errorf("command failed with %w", err)
				} else {
					fmt.Println("{\"Cosign\": " + output["cosign"] + ", \"Notation\": " + output["notation"] + "}")
				}
			} else {
				tag := util.ImageSuffixTag(imgName, digest, util.SuffixSigTag(kinds[0]))
				sig, err := crane.Get(tag, *options...)
				if err != nil {
					return err
				} else {
					fmt.Println(string(sig.Manifest))
				}
			}
			return nil
		},
	}
	sigCmd.PersistentFlags().StringVar(&filter, "filter", "", "Filter output by kind: cosign/notation. Default no filter means getting all types.")
	return sigCmd
}
