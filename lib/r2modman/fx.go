package r2modman

import (
	"time"
)

type Config struct {
	InstallDirectory          string
	WorkDirectory             string
	ThunderstoreForceDownload bool
	ThunderstoreCDN           string
	ThunderstoreCDNTimeout    time.Duration
}
