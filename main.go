package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AWS - Handles the AWS Config
type AWS struct {
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

// REPO - Handles the Repo Config
type REPO struct {
	Name       string
	URL        string
	Release    string
	Components string
	Arch       string
	GPGKeys    string
}

// Global Config Vars
var aws AWS
var repos []REPO
var gpgservers []string

var aptlyPath string
var cfgPath string

// Main Function
func main() {

	// CLI: mirror command
	var cmdMirror = &cobra.Command{
		Use:   "mirror",
		Short: "Mirror repos from config file",
		Long:  `Mirrors the repos described in the config file using aptly`,
		Run: func(cmd *cobra.Command, args []string) {
			mirror()
		},
	}

	// CLI: snap command
	var cmdSnap = &cobra.Command{
		Use:   "snap",
		Short: "Snapshot aptly mirror",
		Long:  `Snapshots a mirrored repository in aptly`,
		Run: func(cmd *cobra.Command, args []string) {
			snapshot()
		},
	}

	// CLI: publish command
	var cmdPublish = &cobra.Command{
		Use:   "publish",
		Short: "Publish aptly mirror",
		Long:  `Publishes an aptly snapshot to AWS S3`,
		Run: func(cmd *cobra.Command, args []string) {
			publish()
		},
	}

	// Configure CLI options
	var rootCmd = &cobra.Command{Use: "aptly_mirror"}
	rootCmd.PersistentFlags().StringVarP(&aptlyPath, "aptly-path", "", "aptly", "Path to Aptly Binary")
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "mirrors-path", "", "/etc", "Path to ./aptly_mirrors.yaml file containing Aptly Mirrors")
	rootCmd.AddCommand(cmdMirror, cmdSnap, cmdPublish)
	rootCmd.Execute()

}

// aptly mirror create/update
func mirror() {

	// Declare config file
	viper.SetConfigName("aptly_mirrors")
	viper.AddConfigPath(cfgPath)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found...")
	} else {

		// Read repos from config file
		err := viper.UnmarshalKey("repos", &repos)
		if err != nil {
			fmt.Println("Cannot parse repos...")
		}

		// Read gpg servers from config files
		gpgservers = viper.GetStringSlice("gpg_servers")
		os.Setenv("GNUPGHOME", "/root/.gnupg")

		// Get list of existing mirrors before starting loop
		mirrorlist, err := exec.Command("aptly", "mirror", "list").CombinedOutput()
		if err != nil {
			fmt.Println("ERROR: Cannot check aptly mirrors list..")
			log.Fatal(err)
		}

		// Loop through each aptly mirror
		for _, repo := range repos {

			// If current repo doesnt already exist as a mirror
			if !(strings.Contains(string(mirrorlist[:]), repo.Name)) {

				// Import GPG Key
				for _, gpgserver := range gpgservers {
					fmt.Println("Importing GPG Keys..\n")
					runCommand("gpg", []string{"--no-default-keyring", "--keyring", "trustedkeys.gpg", "--keyserver", "hkp://" + gpgserver + ":80", "--recv-keys", repo.GPGKeys})
				}

				// Create Mirror
				runCommand("aptly", []string{"-keyring=/root/.gnupg/trustedkeys.gpg", "-architectures=" + repo.Arch, "mirror", "create", repo.Name, repo.URL, repo.Release, repo.Components})
			}

			// Update Mirror
			fmt.Printf("Updating Mirror %s..\n", repo.Name)
			runCommand("aptly", []string{"-keyring=/root/.gnupg/trustedkeys.gpg", "-architectures=" + repo.Arch, "mirror", "update", repo.Name})

		}
	}
}

// TODO: aptly snapshot create [mirror]
func snapshot() {}

// TODO: aptly publish snapshot [snapshot]
func publish() {}

// Function to Run OS Commands and display output
func runCommand(c string, args []string) {
	cmd := exec.Command(c, args...)

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", stdout)
		log.Fatal(err)
	}

	fmt.Printf("%s", stdout)

	// return string(stdout[:])
}
