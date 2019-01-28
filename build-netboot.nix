# Copyright (C) 2019 Néfix Estrada <nefixestrada@gmail.com>
# Author: Néfix Estrada <nefixestrada@gmail.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

{ pkgs ? import <nixpkgs> {}, system, ... }: 
  
let   
  configuration = { pkgs, config, lib, ... }: {  
    services.xserver = {
	  enable = true;
	  displayManager = {
		slim = {
		  enable = true;
		  autoLogin = true;
		  defaultUser = "isard";
		};
		sessionCommands = ''
		  # Execute the X WM
		  ratpoison &

		  # Read all the boot parameters and set them as variables
		  for i in $(cat /proc/cmdline | tr ' ' '\n'); do
			name=$(echo -n "$i" | cut -f1 -d=)
			value=$(echo -n "$i" | cut -f2 -d=)

            if [ "$name" = "tkn" -o "$name" = "id" -o "$name" = "base_url" ]; then
              eval $name=\$value
            fi
		  done

		  # Download the console.vv file
		  wget "$base_url/pxe/viewer?tkn=$tkn&id=$id" -O console.vv

		  exec remote-viewer -fk console.vv
		'';
	  };
	};
	hardware.pulseaudio = {
	  enable = true;
	};
	environment.systemPackages = with pkgs; [ ratpoison virt-viewer wget ];
	users.extraUsers.isard = {
	  isNormalUser = true;
	  uid = 1100;
	};
  };  

  netboot = (import (pkgs.path + "/nixos/lib/eval-config.nix") {
    inherit system;
    modules = [   
      (pkgs.path + "/nixos/modules/installer/netboot/netboot-minimal.nix")
      configuration  
    ];
  });

  ipxeScript = pkgs.writeTextDir "netboot.ipxe" ''
    #!ipxe
    kernel {{.BaseURL}}/pxe/boot/vmlinuz?arch=''${buildarch:uristring} base_url={{.BaseURL}} tkn={{.Token}} id={{.VMID}} init=${netboot.config.system.build.toplevel}/init ${toString netboot.config.boot.kernelParams}
    initrd {{.BaseURL}}/pxe/boot/initrd?arch=''${buildarch:uristring}
    boot
  '';
in
  
pkgs.symlinkJoin {
  name = "netboot";   
  paths = with netboot.config.system.build; [ netbootRamdisk  kernel ipxeScript ];   
}
