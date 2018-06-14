package model

// Flags sets the basic flags for the treads,intensity,etc
type Flags struct {
	CPU int `long:"treads" short:"t" description:"CPU threads count for specified currency"`

	Intensity int `long:"intensity" short:"i" default:"28" description:"GPU mining intensity (NVidia only) (values range: 1..4. Recommended: 2)"`

	ExcludeGPUS string `long:"excluded" short:"e"  description:"GPUs Excluded. (values: 0,1,2,3)"`

	Proxy string `string:"proxy" description:"Proxy server URL. Supports only socks protocols (for example: socks://192.168.0.1:1080"`

	URL      string `long:"url" short:"r" description:"Mining pool URL"`
	Username string `long:"username" short:"u" description:"Mining pool username"`
	Password string `long:"password" short:"p" description:"Mining pool password"`

	Version bool `long:"version" short:"v" description:"Display version and exit"`

	CoinType string `long:"cointype" short:"c" description:"Coin Type eg: (SIA,BTC,ETH,ETC,XRM,VRM)"`
}
