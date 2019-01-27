package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron"
	"github.com/segmentio/ksuid"
)

// supportedArchitectures is an slice that contains all the supoprted architectures to build
var supportedArchitectures = []string{"x86_64", "i386"}
var netbootArchitectures = map[string]string{
	"x86_64": "x86_64-linux",
	"i386":   "i686-linux",
}

// jobs is a map that contains all the jobs that are running or that have been finished less than 24 hours ago
var jobs = map[ksuid.KSUID]*jobStatus{}

// jobStatus is the status of each job
type jobStatus struct {
	hasFinished bool
	started     time.Time
	finished    time.Time
}

func main() {
	c := cron.New()
	c.AddFunc("0 0 2 * * *", buildNetboot)
	c.Start()

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Println("Starting IsardVDI builder at port :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Printf("error listening at port :3000: %v", err)
	}
}

// buildNetboot builds the Netboot images for all the CPU architectures and makes a symlink to the public webroot
func buildNetboot() {
	for _, arch := range supportedArchitectures {
		id := ksuid.New()

		if err := nixBuld(id, "build-netboot.nix", map[string]string{"system": netbootArchitectures[arch]}); err != nil {
			log.Printf("error building netboot for %s: %v", arch, err)
		} else if err := linkToWebRoot(arch); err != nil {
			log.Printf("error moving the netboot for %s to webroot: %v", arch, err)
		}

		jobs[id].finished = time.Now()
		jobs[id].hasFinished = true
	}

	log.Println("successfully built " + time.Now().Format("2006-01-02") + " images")
}

// nixBuild builds the Nix expression
func nixBuld(id ksuid.KSUID, expression string, args map[string]string) error {
	jobs[id] = &jobStatus{
		hasFinished: false,
		started:     time.Now(),
	}

	// parse the args
	var cmdArgs = []string{"-o", expression + "-result", expression}

	for k, v := range args {
		cmdArgs = append(cmdArgs, "--argstr")
		cmdArgs = append(cmdArgs, k)
		cmdArgs = append(cmdArgs, v)
	}

	// build the expression
	out, err := exec.Command("nix-build", cmdArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error building %s: %v\n%s", expression, err, out)
	}

	return nil
}

// linkToWebRoot links the generated binaries to the public server webroot so they can be downloaded. Also creates the required directories
func linkToWebRoot(arch string) error {
	basePath := "public/" + arch + "/"
	pathToday := time.Now().Format("2006-01-02")
	path := basePath + pathToday
	latestPath := basePath + "latest"

	// create the public directory
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("unexpected error checking the directory status: %v", err)
		}

		if err = os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("error creating the public directory: %v", err)
		}
	}

	// make links to the actual build
	vmlinuzPath, err := os.Readlink("build-netboot.nix-result/bzImage")
	if err != nil {
		return fmt.Errorf("error reading the vmlinuz link destination: %v", err)
	}

	initrdPath, err := os.Readlink("build-netboot.nix-result/initrd")
	if err != nil {
		return fmt.Errorf("error reading the initrd link destination: %v", err)
	}

	menuPath, err := os.Readlink("build-netboot.nix-result/netboot.ipxe")
	if err != nil {
		return fmt.Errorf("error reading the menu link destination: %v", err)
	}

	if err := os.Symlink(vmlinuzPath, path+"/vmlinuz"); err != nil {
		return fmt.Errorf("error creating the vmlinuz link: %v", err)
	}

	if err = os.Symlink(initrdPath, path+"/initrd"); err != nil {
		return fmt.Errorf("error creating the initrd link: %v", err)
	}

	if err = os.Symlink(menuPath, path+"/netboot.ipxe"); err != nil {
		return fmt.Errorf("error creating the menu link: %v", err)
	}

	if _, err := os.Stat(latestPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("unexpected error checking the 'latest' directory status: %v", err)
		}
	} else {
		if err = os.Remove(latestPath); err != nil {
			return fmt.Errorf("error removing the latest link: %v", err)
		}
	}

	if err = os.Symlink(pathToday, latestPath); err != nil {
		return fmt.Errorf("error creating the latest link: %v", err)
	}

	return nil
}
