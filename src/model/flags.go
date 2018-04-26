package model

// Flags sets the basic flags for the treads,intensity,etc
type Flags struct {
	CPU         int    `long:"treads" short:"t" description:"CPU threads count for specified currency"`
	Intensity   int    `long:"intensity" short:"i" default:"28" description:"GPU mining intensity (NVidia only) (values range: 1..4. Recommended: 2)"`
	ExcludeGPUS string `string:"e" description:"GPUs Excluded. (values: 0,1,2,3)"`
	Proxy       string `string:"proxy" description:"Proxy server URL. Supports only socks protocols (for example: socks://192.168.0.1:1080"`
	Host        string `string:"o" description:"Mining pool URL"`
	Login       string `string:"u" description:"Mining pool login"`
	Version     bool   `long:"version" short:"v" description:"Display version and exit"`
}
