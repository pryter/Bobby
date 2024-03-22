package resources

import (
	. "Bobby/pkg/resources"
	"github.com/rs/zerolog/log"
)

// ResourceTree singleton is a declared resourceTree for the application
type ResourceTree struct {
	Configs ConfigResources
}

type ConfigResources struct {
	NetworkConfigFile File
	Folder
}

var resources *ResourceTree

func GetResources() *ResourceTree {
	if resources == nil {
		log.Panic().Msg("ResourceTree must instantiate once before using.")
	}

	return resources
}

func InitResourceTree(resourceBasePath string) error {

	configFolder := Folder{ResourcePath: resourceBasePath, Name: "configs"}
	configFolder.CreateIfNotExist()

	resources = &ResourceTree{
		Configs: ConfigResources{
			NetworkConfigFile: configFolder.MapFile("network-credentials.json"),
			Folder:            configFolder,
		},
	}

	return nil
}
