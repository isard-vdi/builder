/*
 * Copyright (C) 2019 Néfix Estrada <nefixestrada@gmail.com>
 * Author: Néfix Estrada <nefixestrada@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

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

// nixBuild builds the Nix expression
func nixBuld(id ksuid.KSUID, expression, result string, args map[string]string) error {
	// parse the args
	var cmdArgs = []string{"-o", result, expression}

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

// buildNetboot builds the Netboot images for all the CPU architectures. It also manages the image publication
func buildNetboot() {
	for _, arch := range supportedArchitectures {
		id := ksuid.New()
		jobs[id] = &jobStatus{
			hasFinished: false,
			started:     time.Now(),
		}

		if err := nixBuld(id, "build-netboot.nix", "build-netboot-"+arch+"-result", map[string]string{"system": netbootArchitectures[arch]}); err != nil {
			log.Printf("error building netboot for %s: %v", arch, err)
		} else if err := publishNetboot(arch); err != nil {
			log.Printf("error moving the netboot for %s to webroot: %v", arch, err)
		}

		jobs[id].finished = time.Now()
		jobs[id].hasFinished = true
	}

	log.Println("successfully built " + time.Now().Format("2006-01-02") + " images")
}
