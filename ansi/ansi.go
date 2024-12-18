package ansi

var colors = map[string]string{
	// BASIC COLORS
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"reset":   "\033[0m",

	// BRIGHT/PASTEL COLORS
	"black_bright":   "\033[90m",
	"red_bright":     "\033[91m",
	"green_bright":   "\033[92m",
	"yellow_bright":  "\033[93m",
	"blue_bright":    "\033[94m",
	"magenta_bright": "\033[95m",
	"cyan_bright":    "\033[96m",
	"white_bright":   "\033[97m",

	// BRIGHT + BOLD COLORS
	"black_bright_bold":   "\033[1;90m",
	"red_bright_bold":     "\033[1;91m",
	"green_bright_bold":   "\033[1;92m",
	"yellow_bright_bold":  "\033[1;93m",
	"blue_bright_bold":    "\033[1;94m",
	"magenta_bright_bold": "\033[1;95m",
	"cyan_bright_bold":    "\033[1;96m",
	"white_bright_bold":   "\033[1;97m",
}

// GetColor returns the ANSI escape code for a given color
// in O(1) time
func GetColor(color string) string {
	return colors[color]
}

// Colorize colorizes a string with a given color and adds
// a reset code at the end, in O(1) time
func Colorize(color string, text string) string {
	return GetColor(color) + text + GetColor("reset")
}
