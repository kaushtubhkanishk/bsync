package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kaushtubhkanishk/bsync/internal/packagemanager"
)

func main() {
	if os.Getenv("SOURCE_BUILD_DIR") == "" {
		err := errors.New("env: SOURCE_BUILD_DIR is not set")
		panic(err)
	}
	if os.Getenv("SOURCE_MANIFEST_PATH") == "" {
		err := errors.New("env: SOURCE_MANIFEST_PATH is not set")
		panic(err)
	}

	packs, err := packagemanager.ReadManifest()
	if err != nil {
		panic(err)
	}

	for i, pack := range packs {
		fmt.Println("Package: ", pack.Name)
		pack.Path = strings.TrimSuffix(pack.Path, "/")

		err = pack.FetchLatestVersion()
		if err != nil {
			panic(err)
		}

		if pack.LatestVersion != pack.CurrentVersion {
			fmt.Printf("New version found, current version: %v, new version: %v\n", pack.CurrentVersion, pack.LatestVersion)

			err = pack.Update()
			if err != nil {
				panic(err)
			}
		}
		packs[i] = pack
	}

	err = packagemanager.UpdateManifest(&packs)
	if err != nil {
		panic(err)
	}
}
