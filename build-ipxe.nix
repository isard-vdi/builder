{ pkgs ? import <nixpkgs> {}, baseUrl, caCert, ... }:

pkgs.ipxe.overrideAttrs(old: {
  script = pkgs.writeText "embed.ipxe" ''
    #!ipxe
    dhcp
    echo
    echo Connecting to IsardVDI...
    echo
    chain ${baseUrl}/pxe/boot
  '';

  makeFlags = old.makeFlags ++ [
    ''EMBED=''${script}''
    ''TRUST=''${caCert}''
  ];

  enabledOptions = old.enabledOptions ++ [
    "REBOOT_CMD"
    "POWEROFF_CMD"
  ];

  installPhase = ''
    ${old.installPhase}
    make $makeFlags bin-x86_64-efi/ipxe.efi bin-i386-efi/ipxe.efi

    mkdir $out/x86_64
    mkdir $out/i386
    cp -v bin-x86_64-efi/ipxe.efi $out/x86_64/ipxe.efi
    cp -v bin-i386-efi/ipxe.efi $out/i386/ipxe.efi
  '';
})
