package app

import (
	. "Bobby/pkg/resources"
	"github.com/rs/zerolog/log"
)

// ResourceTree singleton is a declared resourceTree for the application
type ResourceTree struct {
	WorkerData Folder
}

var resources *ResourceTree

func GetResources() *ResourceTree {
	if resources == nil {
		log.Panic().Msg("ResourceTree must instantiate once before using.")
	}

	return resources
}

func InitResourceTree(resourceBasePath string) error {

	workerFolder := Folder{ResourcePath: resourceBasePath, Name: "workers"}
	workerFolder.CreateIfNotExist()

	resources = &ResourceTree{
		WorkerData: workerFolder,
	}

	return nil
}
