{ pkgs, lib, config, inputs, ... }:

let
  unstable = import inputs.unstable { system = pkgs.stdenv.system; };
in {
  # https://devenv.sh/packages/
  packages = [ pkgs.git unstable.go-mockery unstable.air ];

  languages.go = {
    enable = true;
  };
}
