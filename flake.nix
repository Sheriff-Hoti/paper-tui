{
  description = "TUI Wallpaper Selector";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.05";
  };

  outputs =
    {
      self,
      nixpkgs,
    }:
    let
      #System types to support.
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

      version = "0.1.0";
      pname = "paper-tui";
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.buildGoModule {
            inherit pname;
            inherit version;

            src = ./.;

            vendorHash = "sha256-OJONkd8N51VC6OlydKHpb/AkC6SoZ+zsJ4tfIDNbKQ8=";
          };
        }
      );

      defaultPackage = forAllSystems (system: self.packages.${system}.default);

      devShell = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        with pkgs;
        mkShell {
          buildInputs = [
            go
            gopls
          ];
        }
      );

      nixosModule = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          config,
          lib,
          pkgs,
          ...
        }:
        with lib;
        {
          options.programs.paper-tui = {
            enable = mkEnableOption "paper-tui enable";

            package = mkOption {
              type = types.package;
              default = self.packages.${system}.default;
              description = "paper-tui package to use";
            };

          };

        }
      );

      homeModule = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          config,
          lib,
          pkgs,
          ...
        }:
        with lib;
        let
          jsonFormat = pkgs.formats.json { };
        in
        {
          options.programs.paper-tui = {
            enable = mkEnableOption "paper-tui enable";

            package = mkOption {
              type = types.package;
              default = self.packages.${system}.default;
              description = "paper-tui package to use";
            };

            settings = mkOption {
              type = lib.types.nullOr jsonFormat.type;
              default = null;
              description = "JSON configuration settings for paper-tui";
              example = {
                wallpapers_dir = "~/Pictures/Wallpapers";
                backend = "swayb";
              };
            };
          };

          config = mkIf config.programs.paper-tui.enable (
            let
              cfg = config.programs.paper-tui;
            in
            {
              home.packages = [ cfg.package ];

              xdg.configFile = lib.mkIf (cfg.settings != null) {
                "paper-tui/config.json".source = jsonFormat.generate "paper-tui-config.json" cfg.settings;
              };
            }
          );

        }
      );
    };
}
