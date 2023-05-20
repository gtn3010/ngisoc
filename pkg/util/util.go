package util

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
)

func ImageWithoutDigest(imgName string) string {
	return strings.Split(imgName, "@")[0]
}

func ImageWithoutTag(imgName string) string {
	return strings.Split(imgName, ":")[0]
}

func ImageSuffixTag(imgName string, digest string, suffix string) string {
	imgWithoutTag := ImageWithoutTag(imgName)
	sigTag := strings.Replace(digest, ":", "-", 1) + suffix
	return imgWithoutTag + ":" + sigTag
}

func Contains[T comparable](sl []T, elem T) bool {
	for i := 0; i < len(sl); i++ {
		if elem == sl[i] {
			return true
		}
	}
	return false
}

var suffixSig = map[string]string{
	"cosign":   ".sig",
	"notation": "",
}

func SuffixSigTag(suffix string) string {
	return suffixSig[suffix]
}

var suffixSbom = map[string]string{
	"cosign": ".sbom",
}

func SuffixSbomTag(suffix string) string {
	return suffixSbom[suffix]
}

func FilterArgs(filter string, allowedArgs []string) ([]string, error) {
	var kinds []string
	if filter == "" {
		kinds = allowedArgs
	} else if filter != "" && Contains[string](allowedArgs, filter) {
		kinds = append(kinds, filter)
	} else {
		return nil, fmt.Errorf("only supported filter: %v", strings.Join(allowedArgs, "/"))
	}
	return kinds, nil
}

func GetDigest(imgName string, options *[]crane.Option) (string, error) {
	manifest, err := crane.Get(imgName, *options...)
	if err != nil {
		return "", fmt.Errorf("fetching manifest %s: %w", imgName, err)
	}
	digest := manifest.Descriptor.Digest.String()
	return digest, nil
}

func ConvertToByte(r io.Reader) string {
	buf := new(strings.Builder)
	n, _ := io.Copy(buf, r)
	// fmt.Println(buf.String())
	fmt.Println(n)
	return buf.String()
}
