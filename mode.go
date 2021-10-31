package there

const (
	envThereMode = "THERE_MODE"

	// DebugMode provides more logging information
	DebugMode = "debug"
	// ProductionMode should be used when deploying your app to production
	ProductionMode = "production"
)

type Mode string

func (mode *Mode) IsProduction() bool {
	return *mode == ProductionMode
}

func (mode *Mode) IsDebug() bool {
	return *mode == DebugMode
}

func (mode *Mode) SetProduction() {
	*mode = ProductionMode
}

func (mode *Mode) SetDebug() {
	*mode = DebugMode
}
