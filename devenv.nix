{ pkgs, lib, config, inputs, ... }:

{
  # https://devenv.sh/packages/
  packages = [ pkgs.git pkgs.claude-code ];

  languages.go = {
    enable = true;
  };
}
