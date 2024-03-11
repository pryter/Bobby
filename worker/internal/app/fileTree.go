package app

import . "Bobby/pkg/resources"

type ConfigResources struct {
	NetworkConfigFile File
}

type ResourceTree struct {
	Configs ConfigResources
}

func NewResourceTree(resourceBasePath string) ResourceTree {

	configFolder := Folder{ResourcePath: resourceBasePath, Name: "configs"}

	return ResourceTree{
		Configs: ConfigResources{
			NetworkConfigFile: configFolder.MapFile("network-credentials.json"),
		},
	}
}
