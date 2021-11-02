package utils

func GreenString(str string) string {
	return colorString(str, "\033[32m")
}

func RedString(str string) string {
	return colorString(str, "\033[31m")
}

func YellowString(str string) string {
	return colorString(str, "\033[33m")
}

func colorString(str string, colorPrefix string) string {
	return colorPrefix + str + "\033[0m"
}
