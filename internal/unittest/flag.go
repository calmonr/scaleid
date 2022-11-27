package unittest

func ComposeFlag(name, value string) string {
	return "--" + name + "=" + value
}
