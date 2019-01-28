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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// copyFile copies the content of the file to the destination
func copyFile(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error copying %s: error reading the file: %v", src, err)
	}
	defer f.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error copying %s: error creating the file: %v", dst, err)
	}
	defer out.Close()

	_, err = io.Copy(out, f)
	if err != nil {
		return fmt.Errorf("error copying %s: %v", dst, err)
	}

	if err = out.Sync(); err != nil {
		return fmt.Errorf("error syncing %s: %v", dst, err)
	}

	return nil
}

// publishNetboot publishes the compiled Netboot image to the public folder. It creates a link if the files haven't changed and copies them if they have.
func publishNetboot(arch string) error {
	basePath := filepath.Join("public", arch)
	latestPath := filepath.Join(basePath, "latest")

	today := time.Now().Format("2006-01-02")
	yesteday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	todayPath := filepath.Join(basePath, today)
	yestedayPath := filepath.Join(basePath, yesteday)

	// get the path to the actual binaries
	vmlinuzPath, err := os.Readlink("build-netboot-" + arch + "-result/bzImage")
	if err != nil {
		return fmt.Errorf("error reading the vmlinuz link destination: %v", err)
	}

	initrdPath, err := os.Readlink("build-netboot-" + arch + "-result/initrd")
	if err != nil {
		return fmt.Errorf("error reading the initrd link destination: %v", err)
	}

	netbootPath, err := os.Readlink("build-netboot-" + arch + "-result/netboot.ipxe")
	if err != nil {
		return fmt.Errorf("error reading the netboot.ipxe link destination: %v", err)
	}

	// check the SHA256's
	vmlinuzSHA256, err := getSHA256(vmlinuzPath)
	if err != nil {
		return fmt.Errorf("error checking the SHA256 of the vmlinuz: %v", err)
	}

	initrdSHA256, err := getSHA256(initrdPath)
	if err != nil {
		return fmt.Errorf("error checking the SHA256 of the initrd: %v", err)
	}

	netbootSHA256, err := getSHA256(netbootPath)
	if err != nil {
		return fmt.Errorf("error checking the SHA256 of the netboot: %v", err)
	}

	// create today's directory
	if _, err := os.Stat(todayPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("unexpected error checking the directory status: %v", err)
		}

		if err = os.MkdirAll(todayPath, 0755); err != nil {
			return fmt.Errorf("error creating the public directory: %v", err)
		}
	}

	// check if yesteday's build exist
	if _, err := os.Stat(yestedayPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("unexpected error checking %s directory status: %v", yestedayPath, err)
		}

		if err := copyFile(vmlinuzPath, filepath.Join(todayPath, "vmlinuz")); err != nil {
			return err
		}

		if err := copyFile(initrdPath, filepath.Join(todayPath, "initrd")); err != nil {
			return err
		}

		if err := copyFile(netbootPath, filepath.Join(todayPath, "netboot.ipxe")); err != nil {
			return err
		}

	} else {
		// check yesteday's SHA256's
		oldVmlinuzSHA256, err := getSHA256(filepath.Join(yestedayPath, "vmlinuz"))
		if err != nil {
			return fmt.Errorf("error checking the SHA256 of the old vmlinuz: %v", err)
		}

		oldInitrdSHA256, err := getSHA256(filepath.Join(yestedayPath, "initrd"))
		if err != nil {
			return fmt.Errorf("error checking the SHA256 of the old initrd: %v", err)
		}

		oldNetbootSHA256, err := getSHA256(filepath.Join(yestedayPath, "netboot.ipxe"))
		if err != nil {
			return fmt.Errorf("error checking the SHA256 of the old netboot: %v", err)
		}

		// check if files have changed
		if vmlinuzSHA256 != oldVmlinuzSHA256 {
			if err := copyFile(vmlinuzPath, filepath.Join(todayPath, "vmlinuz")); err != nil {
				return err
			}
		} else {
			if err := os.Symlink(filepath.Join("..", yesteday, "vmlinuz"), filepath.Join(todayPath, "vmlinuz")); err != nil {
				return err
			}
		}

		if initrdSHA256 != oldInitrdSHA256 {
			if err := copyFile(initrdPath, filepath.Join(todayPath, "initrd")); err != nil {
				return err
			}
		} else {
			if err := os.Symlink(filepath.Join("..", yesteday, "initrd"), filepath.Join(todayPath, "initrd")); err != nil {
				return err
			}
		}

		if netbootSHA256 != oldNetbootSHA256 {
			if err := copyFile(netbootPath, filepath.Join(todayPath, "netboot.ipxe")); err != nil {
				return err
			}
		} else {
			if err := os.Symlink(filepath.Join("..", yesteday, "netboot.ipxe"), filepath.Join(todayPath, "netboot.ipxe")); err != nil {
				return err
			}
		}
	}

	sha256sums := [][2]string{
		[2]string{"vmlinuz", vmlinuzSHA256},
		[2]string{"initrd", initrdSHA256},
		[2]string{"netboot.ipxe", netbootSHA256},
	}

	if err := ioutil.WriteFile(filepath.Join(todayPath, "sha256sum.txt"), []byte(sha256sum(sha256sums)), 0644); err != nil {
		return fmt.Errorf("error writing the SHA256 sums file: %v", err)
	}

	// TODO: GPG sign the SHA256 file

	if _, err := os.Stat(latestPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("unexpected error checking the 'latest' directory status: %v", err)
		}
	} else {
		if err = os.Remove(latestPath); err != nil {
			return fmt.Errorf("error removing the latest link: %v", err)
		}
	}

	if err = os.Symlink(today, latestPath); err != nil {
		return fmt.Errorf("error creating the latest link: %v", err)
	}

	return nil
}
