// +build release

package cmd

import (
	"fmt"
)

func activeProfile(dataDir string, isDemo bool) profile {
	dsn := fmt.Sprintf("file:%s/bytebase.db", dataDir)
	seedDir := "seed/release"
	forceResetSeed := false
	if isDemo {
		dsn = fmt.Sprintf("file:%s/bytebase_demo.db", dataDir)
		seedDir = "seed/test"
		forceResetSeed = true
	}
	return profile{
		mode:           "release",
		dsn:            dsn,
		seedDir:        seedDir,
		forceResetSeed: forceResetSeed,
	}
}
