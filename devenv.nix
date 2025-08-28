{ pkgs, lib, config, inputs, ... }:

let
  unstable = import inputs.unstable { system = pkgs.stdenv.system; };
in {
  # https://devenv.sh/packages/
  packages = [ pkgs.git unstable.go-mockery ];

  languages.go = {
    enable = true;
  };
}
