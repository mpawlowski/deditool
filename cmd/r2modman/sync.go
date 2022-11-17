package r2modman

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/mpawlowski/deditool/v2/lib/extractor"
	"github.com/mpawlowski/deditool/v2/lib/r2modman"
	"github.com/spf13/cobra"
)

var profileExportR2Z string
var workDir string
var installDir string
var forceDownload bool
var thunderstoreCDNHost string
var thunderstoreCdnTimeout time.Duration

func init() {
	syncCmd.PersistentFlags().StringVar(&profileExportR2Z, "r2z-file", "", "Path to a zip export from R2ModmanPlus.")
	syncCmd.PersistentFlags().StringVar(&workDir, "work-dir", "work", "A work directory for downloaded files.")
	syncCmd.PersistentFlags().StringVar(&installDir, "install-dir", "", "Path to a dedicated server installation.")
	syncCmd.PersistentFlags().BoolVar(&forceDownload, "force-download", false, "Force re-downloading mods from Thunderstore.")
	syncCmd.PersistentFlags().StringVar(&thunderstoreCDNHost, "thunderstore-cdn-host", "gcdn.thunderstore.io", "Hostname of the thunderstore CDN to use.")
	syncCmd.PersistentFlags().DurationVar(&thunderstoreCdnTimeout, "thunderstore-cdn-timeout", 30*time.Second, "Timeout while downloading each mod.")
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize mods for a dedicated install.",
	Long:  `Synchronize mods for a dedicated install.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if profileExportR2Z == "" {
			return fmt.Errorf("missing argument r2z-file")
		}

		if installDir == "" {
			return fmt.Errorf("missing argument install-dir")
		}

		modExtractor := extractor.NewExtractor()
		modUtil, err := r2modman.NewModUtil(r2modman.Config{
			InstallDirectory:          installDir,
			WorkDirectory:             workDir,
			ThunderstoreForceDownload: forceDownload,
			ThunderstoreCDN:           thunderstoreCDNHost,
			ThunderstoreCDNTimeout:    thunderstoreCdnTimeout,
		})
		if err != nil {
			return err
		}

		parser, err := r2modman.NewExportParser()
		if err != nil {
			return err
		}

		log.Println("Using profile", profileExportR2Z)
		metadata, err := parser.Parse(profileExportR2Z)
		if err != nil {
			return err
		}

		packages, err := r2modman.GetPackagesMetadata(cmd.Context())
		if err != nil {
			log.Printf("unable to pull thunderstore api: %s", err)
			return err
		}

		log.Printf("Found packages %d packages from thunderstore\n", len(packages))

		for _, v := range metadata.Mods {

			downloadedZipPath := path.Join(workDir, v.Filename())

			thunderstoreMeta, ok := packages[v.ThunderstoreKey()]
			if !ok {
				return fmt.Errorf("thunderstore metadata does not exist for: %s", v.ThunderstoreKey())
			}

			err = modUtil.Download(v, thunderstoreMeta)
			if err != nil {
				return err
			}

			//extract modes to install directory
			packagingType, prefixToStrip, err := r2modman.DeterminePackagingType(downloadedZipPath)
			if err != nil {
				return err
			}

			log.Printf("packaging type %v", packagingType)

			installDir := fmt.Sprintf("%s/%s", installDir, packagingType.Directory())
			err = modExtractor.Extract(downloadedZipPath, installDir, prefixToStrip)
			if err != nil {
				return err
			}
		}

		// extract profile to bepinex in install dir
		bepinDir := path.Join(installDir, "/BepInEx")
		log.Println(fmt.Sprintf("Extracting %s to %s", profileExportR2Z, bepinDir))
		err = modExtractor.Extract(profileExportR2Z, bepinDir, "")
		if err != nil {
			return err
		}

		log.Printf("Mod install finished successfully, configure your start script at %s/start_server_bepinex.sh\n", installDir)

		return nil
	},
}
