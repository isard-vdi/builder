#
# Build stage
#

# Use Golang 1.11 as build stage
FROM golang:1.11 as build

# Copy Isard Build
COPY . /go/src/github.com/isard-vdi/builder

# Move to the correct directory
WORKDIR /go/src/github.com/isard-vdi/builder

# Compile the binaries
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-s -w" -o builder cmd/builder/main.go
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-s -w" -o force-build cmd/builder/force-build.go

#
# Base stage
#

# Use Nix as base stage
FROM nixos/nix

# Use the nixpkgs unstable channel
RUN nix-channel --add https://nixos.org/channels/nixpkgs-unstable nixpkgs
RUN nix-channel --update

# Copy the compiled binaries from the build stage
COPY --from=build /go/src/github.com/isard-vdi/builder/builder /app/builder
COPY --from=build /go/src/github.com/isard-vdi/builder/force-build /bin/force-build

# Create the data directory
RUN mkdir /data

# Copy the build Nix expressions
COPY build-netboot.nix /data/build-netboot.nix
COPY build-ipxe.nix /data/build-ipxe.nix

# Move to the correct directory
WORKDIR /data

# Expose the volume
VOLUME [ "/data/public" ]

# Expose the required port
EXPOSE 1312

# Run the service
CMD [ "/app/builder" ]
