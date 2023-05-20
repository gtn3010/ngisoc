package signature

import (
	"sync"

	"github.com/google/go-containerregistry/pkg/crane"

	"github.com/gtn3010/ngisoc/pkg/util"
)

func GetAll(imgName string, digest string, kinds []string, options *[]crane.Option) (map[string]string, error) {
	result := map[string]string{
		"cosign":   `"No Signature"`,
		"notation": `"No Signature"`,
	}
	var wg sync.WaitGroup
	for _, k := range kinds {
		wg.Add(1)
		go func(k string) {
			tag := util.ImageSuffixTag(imgName, digest, util.SuffixSigTag(k))
			sig, err := crane.Get(tag, *options...)
			if err == nil {
				result[k] = string(sig.Manifest)
			}
			defer wg.Done()
		}(k)
	}
	wg.Wait()
	return result, nil
}
