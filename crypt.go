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
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// getSHA256 calculates the SHA256 of a file and returns it
func getSHA256(src string) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("error checking the SHA256: error reading %s: %v", src, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("error checking the SHA256: %v", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// sha256sum returns a string with the sha256sum format. In the first array item goes the filename and in the second the SHA256
func sha256sum(sums [][2]string) string {
	var rsp string

	for _, sum := range sums {
		rsp += sum[1] + " *" + sum[0] + "\n"
	}

	return rsp
}
