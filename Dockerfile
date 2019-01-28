#
# Build stage
#

# Use Golang 1.11 as build stage
FROM golang:1.11 as build

# Copy Isard Build
COPY . /go/src/github.com/isard-vdi/builder

# Move to the correct directory
WORKDIR /go/src/github.com/isard-vdi/builder

# Compile the binary
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-s -w" -o builder .

#
# Base stage
#

# Use Nix as base stage
FROM nixos/nix

# Use the nixpkgs unstable channel
RUN nix-channel --add https://nixos.org/channels/nixpkgs-unstable nixpkgs
RUN nix-channel --update

# Copy the compiled binary from the build stage
COPY --from=build /go/src/github.com/isard-vdi/builder/builder /app/builder

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
EXPOSE 3000

# Run the service
CMD [ "/app/builder" ]
