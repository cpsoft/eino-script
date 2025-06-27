package loaders

import (
	"eino-script/engine/types"
	"errors"
)

type LoaderCreateFunc func(data map[string]interface{}) (types.LoaderInterface, error)

var loaderRegistry = map[string]LoaderCreateFunc{
	"web": CreateWebLoader,
}

func CreateLoaderNode(data map[string]interface{}) (types.LoaderInterface, error) {
	loaderType, ok := data["type"].(string)
	if !ok {
		return nil, errors.New("loaderType is not a string")
	}

	loaderFunc, ok := loaderRegistry[loaderType]
	if !ok {
		return nil, errors.New("loaderType is not a string")
	}

	return loaderFunc(data)
}
