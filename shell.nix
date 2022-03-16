{ pkgs ? import <nixpkgs> {} }:

with pkgs;

mkShell {
  buildInputs = [
      go  
      nodejs
      yarn
  ];
}