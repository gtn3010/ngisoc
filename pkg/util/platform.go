package util

import v1 "github.com/google/go-containerregistry/pkg/v1"

type PlatformValue struct {
	platform *v1.Platform
}

func (pv *PlatformValue) Set(platform string) error {
	p, err := parsePlatform(platform)
	if err != nil {
		return err
	}
	pv.platform = p
	return nil
}

func (pv *PlatformValue) String() string {
	return platformToString(pv.platform)
}

func (pv *PlatformValue) Type() string {
	return "platform"
}

func platformToString(p *v1.Platform) string {
	if p == nil {
		return "all"
	}
	return p.String()
}

func parsePlatform(platform string) (*v1.Platform, error) {
	if platform == "all" {
		return nil, nil
	}

	return v1.ParsePlatform(platform)
}
