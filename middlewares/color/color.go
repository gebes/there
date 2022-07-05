package color

import "fmt"

const (
	BlackColor   = "\033[1;30m%s\033[0m"
	RedColor     = "\033[1;31m%s\033[0m"
	GreenColor   = "\033[1;32m%s\033[0m"
	YellowColor  = "\033[1;33m%s\033[0m"
	BlueColor    = "\033[1;34m%s\033[0m"
	MagentaColor = "\033[1;35m%s\033[0m"
	CyanColor    = "\033[1;36m%s\033[0m"
	WhiteColor   = "\033[1;37m%s\033[0m"
	ResetColor   = "\033[1;0m%s\033[0m"
)

func Black(s string) string {
	return fmt.Sprintf(BlackColor, s)
}

func Red(s string) string {
	return fmt.Sprintf(RedColor, s)
}

func Green(s string) string {
	return fmt.Sprintf(GreenColor, s)
}

func Yellow(s string) string {
	return fmt.Sprintf(YellowColor, s)
}

func Blue(s string) string {
	return fmt.Sprintf(BlueColor, s)
}

func Magenta(s string) string {
	return fmt.Sprintf(MagentaColor, s)
}

func Cyan(s string) string {
	return fmt.Sprintf(CyanColor, s)
}

func White(s string) string {
	return fmt.Sprintf(WhiteColor, s)
}

func Reset(s string) string {
	return fmt.Sprintf(ResetColor, s)
}
