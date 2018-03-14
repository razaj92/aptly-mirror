package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// REPO - Handles the Repo Config
type REPO struct {
	Name       string
	URL        string
	Release    string
	Components string
	Arch       string
	GPGKeys    string
}

// GPG - Handles the GPG Config
type GPG struct {
	Key     string
	Servers []string
}

// APTLY - Handles the AWS Config
type APTLY struct {
	Endpoint string
}

// Global Config Vars
var repos []REPO
var aptly APTLY
var gpg GPG

// Global CLI Flags
var aptlyPath string
var cfgPath string
var debug bool
var publish bool

// Main Function
func main() {

	// CLI: mirror command
	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "Mirrors repos from config file",
		Long:  `Mirrors the repos described in the config file using aptly`,
		Run: func(cmd *cobra.Command, args []string) {
			mirrorRepo()
		},
	}

	cmdRun.Flags().BoolVarP(&publish, "publish", "p", false, "Publish mirrored repos to Aptly Endpoint")

	// Configure CLI options
	var rootCmd = &cobra.Command{Use: "aptly_mirror"}
	rootCmd.PersistentFlags().StringVarP(&aptlyPath, "aptly-path", "", "aptly", "Path to Aptly Binary")
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config-path", "", "/etc", "Path to ./aptly_mirrors.yaml file containing config")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug mode: prints command outputs..")
	rootCmd.AddCommand(cmdRun)

	// RUN CLI
	rootCmd.Execute()

}

// aptly mirror create/update
func mirrorRepo() {

	// Declare config file
	viper.SetConfigName("aptly_mirrors")
	viper.AddConfigPath(cfgPath)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error: Config file not found..")
	} else {

		// Read repos from config file
		err = viper.UnmarshalKey("repos", &repos)
		if err != nil {
			log.Fatal("Cannot parse repos..")
		}

		// Read aptly config from config file
		err = viper.UnmarshalKey("aptly", &aptly)
		if err != nil {
			log.Fatal("Cannot parse repos..")
		}

		// Read gpg config from config file
		err = viper.UnmarshalKey("gpg", &gpg)
		if err != nil {
			log.Fatal("Cannot parse repos..")
		}

	}

	// Get list of existing mirrors before starting loop
	mirrorlist, err := exec.Command("aptly", "mirror", "list").CombinedOutput()
	if err != nil {
		fmt.Println("ERROR: Cannot check aptly mirrors list..")
		log.Fatal(err)
	}

	keyring := "/root/.gnupg/trustedkeys.gpg"
	if os.Getenv("GNUPGHOME") != "" {
		keyring = os.Getenv("GNUPGHOME") + "/trustedkeys.gpg"
	}

	// Loop through each aptly mirror
	for _, repo := range repos {

		fmt.Printf(">>>>>>>> %s - STARTING REPO: %s\n", time.Now().Format("2006-01-02 15:04:05"), repo.Name)

		// If current repo doesnt already exist as a mirror
		if !(strings.Contains(string(mirrorlist[:]), repo.Name)) {
			fmt.Printf(">>>> Mirror for %s does not exist..\n", repo.Name)

			// Import GPG Key
			gpgImport := true
			fmt.Println(">>>> ..Importing GPG Keys")

			for _, gpgserver := range gpg.Servers {

				args := []string{"--no-default-keyring", "--keyring", "trustedkeys.gpg", "--keyserver", "hkp://" + gpgserver + ":80", "--recv-keys", repo.GPGKeys}

				cmd := exec.Command("gpg", args...)

				stdout, err := cmd.CombinedOutput()

				if debug {
					fmt.Printf("%s", stdout)
				}

				if err != nil {
					fmt.Printf(">>>> Could not retrieve key from %s, trying next gpg server\n", gpgserver)
					gpgImport = false
				} else {
					gpgImport = true
					break
				}

			}

			if !gpgImport {
				fmt.Printf(">>>> ERROR: Could not retrieve gpg key, skipping mirror..\n")
				continue
			} else {
				fmt.Printf(">>>> ..Retrieved GPG key\n")
			}

			// Create Mirror
			fmt.Println(">>>> ..Creating Mirror")
			runCommand("aptly", []string{"-keyring=" + keyring, "-architectures=" + repo.Arch, "mirror", "create", repo.Name, repo.URL, repo.Release, repo.Components})
		}

		// Update Mirror
		fmt.Printf(">>>> Updating Mirror for %s..\n", repo.Name)
		runCommand("aptly", []string{"-keyring=" + keyring, "-architectures=" + repo.Arch, "mirror", "update", repo.Name})

		if publish {

			// Create Snapshot from mirror using timestamp
			t := time.Now().Format("20060102150405")
			currentSnap := fmt.Sprintf("snap-%s-%s", repo.Name, t)
			fmt.Printf(">>>> Snapshoting Mirror as %s\n", currentSnap)
			runCommand("aptly", []string{"snapshot", "create", currentSnap, "from", "mirror", repo.Name})

			//Check for previous snapshots
			snapshotlist, err := exec.Command("aptly", "snapshot", "list").CombinedOutput()
			if err != nil {
				fmt.Println("ERROR: Cannot check aptly mirrors list..")
				log.Fatal(err)
			}

			snapshots := strings.Split(string(snapshotlist[:]), "\n")
			c := 0
			for _, snapshot := range snapshots {
				if strings.Contains(snapshot, fmt.Sprintf("[%s]", repo.Name)) && !(strings.Contains(snapshot, t)) {
					c++
				}
			}

			// Publish Snapshot!
			if c == 0 {
				fmt.Printf(">>>> Publising Snapshot: %s\n", currentSnap)
				runCommand("aptly", []string{"-batch=true", "-gpg-key=" + gpg.Key, "-architectures=" + repo.Arch, "publish", "snapshot", currentSnap, aptly.Endpoint + ":" + repo.Name})
			} else {
				fmt.Printf(">>>> Switching Published Snapshot to: %s\n", currentSnap)
				runCommand("aptly", []string{"-batch=true", "-gpg-key=" + gpg.Key, "publish", "switch", repo.Release, aptly.Endpoint + ":" + repo.Name, currentSnap})
			}

		}

		fmt.Printf(">>>>>>>> %s - FINISHED REPO: %s\n\n\n", time.Now().Format("2006-01-02 15:04:05"), repo.Name)

	}
}

// Function to Run OS Commands and display output
func runCommand(c string, args []string) {
	cmd := exec.Command(c, args...)

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", stdout)
		log.Fatal(err)
	}

	if debug {
		fmt.Printf("%s", stdout)
	}
	// return string(stdout[:])
}
