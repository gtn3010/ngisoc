package cmd

import (
	"io"

	"github.com/google/go-containerregistry/pkg/crane"

	ecr "github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/github"
	"github.com/google/go-containerregistry/pkg/v1/google"
	alibabaacr "github.com/mozillazg/docker-credential-acr-helper/pkg/credhelper"

	"github.com/spf13/cobra"
)

var Root = New([]crane.Option{})

func New(options []crane.Option) *cobra.Command {
	insecure := false

	root := &cobra.Command{
		Use:   "ngisoc",
		Short: "Ngisoc is a tool to get container image attestation/provenance/sbom in container registry",
		// RunE:              func(cmd *cobra.Command, args []string) error { return cmd.Usage() },
		DisableAutoGenTag: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			options = append(options, crane.WithContext(cmd.Context()))
			if insecure {
				options = append(options, crane.Insecure)
			}
			kc := authn.NewMultiKeychain(
				authn.DefaultKeychain,
				google.Keychain,
				authn.NewKeychainFromHelper(ecr.NewECRHelper(ecr.WithLogger(io.Discard))),
				authn.NewKeychainFromHelper(credhelper.NewACRCredentialsHelper()),
				authn.NewKeychainFromHelper(alibabaacr.NewACRHelper().WithLoggerOut(io.Discard)),
				github.Keychain,
			)
			options = append(options, crane.WithAuthFromKeychain(kc))
		},
	}
	subcommands := []*cobra.Command{
		NewSbom(&options),
		NewSig(&options),
	}
	root.AddCommand(subcommands...)
	root.PersistentFlags().BoolVar(&insecure, "insecure", false, "Allow image references to be fetched without TLS")
	return root
}
