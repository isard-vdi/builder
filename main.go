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
	"log"
	"net/http"

	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	c.AddFunc("0 0 2 * * *", buildNetboot)
	c.Start()

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Println("Starting IsardVDI builder at port :1312")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Printf("error listening at port :1312: %v", err)
	}
}
