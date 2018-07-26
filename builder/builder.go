package builder

type AppBuilder interface {
	Build(path string) error
}

func GetBuilder(stack string) AppBuilder {
	switch stack {
	case "node":
		return new(NodeBuilder)
	case "dotnet":
		return new(DotnetBuilder)
	}
	return nil
}
