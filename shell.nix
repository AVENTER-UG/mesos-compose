{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  buildInputs = [
    go
    syft
    grype
    docker
    trivy
  ];
}
