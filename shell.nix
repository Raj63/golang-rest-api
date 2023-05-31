{ pkgs ? import (builtins.fetchGit {
         # Descriptive name to make the store path easier to identify                
         name = "purplease-dev-go";                                                 
         url = "https://github.com/NixOS/nixpkgs";                       
         ref = "refs/heads/nixpkgs-unstable";                     
         rev = "7cf5ccf1cdb2ba5f08f0ac29fc3d04b0b59a07e4"; 
}) {} }:

with pkgs;

mkShell {
  buildInputs = [
    clang-tools
    gitlint
    gnupg
    go_1_19
    go-tools
    go-mockery
    gogetdoc
    golangci-lint
    goreleaser
    gosec
    gotools
    #gocritic
    gofumpt
    golint
    #goreturns
    mysql80
    openapi-generator-cli
    postgresql
    pre-commit
  ];

  shellHook =
    ''
      # Setup the binaries installed via `go install` to be accessible globally.
      export PATH="$(go env GOPATH)/bin:$PATH"

      # Install pre-commit hooks.
      pre-commit install

      # Install Go binaries.
      which enumer || go install github.com/dmarkham/enumer@v1.5.3
      which gocritic || go install github.com/go-critic/go-critic/cmd/gocritic@latest
      which goreturns || go install github.com/sqs/goreturns@latest
      which swag || go install github.com/swaggo/swag/cmd/swag@latest
      
      # Add the repo shared gitconfig
      git config --local include.path ../.gitconfig

      # Clear the terminal screen.
      clear
    '';
}
