package there

import (
	"log"
	"testing"
)

func TestMode(t *testing.T) {
	var mode Mode = DebugMode

	if !mode.IsDebug() {
		log.Fatalln("is not debug")
	}

	if mode.IsProduction() {
		log.Fatalln("is production despite being in debug")
	}

	mode.SetProduction()

	if mode.IsDebug() {
		log.Fatalln("is debug despite being in production")
	}

	if !mode.IsProduction() {
		log.Fatalln("is not production")
	}

	mode.SetDebug()

	if !mode.IsDebug() {
		log.Fatalln("is not debug")
	}

	if mode.IsProduction() {
		log.Fatalln("is production despite being in debug")
	}

	mode = "something else"

	if mode.IsDebug() || mode.IsProduction() {
		log.Fatalln("is debug or production despite being not debug or production")
	}

}

func TestModeRouter(t *testing.T) {
	router := NewRouter()

	if !router.IsDebugMode() {
		log.Fatalln("is not debug")
	}

	if router.IsProductionMode() {
		log.Fatalln("is production despite being in debug")
	}

	router.SetProductionMode()

	if router.IsDebugMode() {
		log.Fatalln("is debug despite being in production")
	}

	if !router.IsProductionMode() {
		log.Fatalln("is not production")
	}

	router.SetDebugMode()

	if !router.IsDebugMode() {
		log.Fatalln("is not debug")
	}

	if router.IsProductionMode() {
		log.Fatalln("is production despite being in debug")
	}

	router.mode = "something else"

	if router.IsDebugMode() || router.IsProductionMode() {
		log.Fatalln("is debug or production despite being not debug or production")
	}

}
