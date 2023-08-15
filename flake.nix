{
  description = "Cloudflare DNS IP updater";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "github:NixOS/nixpkgs";

  outputs = { self, nixpkgs }:
    let

      # to work with older version of flakes
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";


      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in
    rec {

      # Provide some binary packages for selected system types.
      packages = forAllSystems
        (system:
          let
            pkgs = nixpkgsFor.${system};
          in
          {
            cloudflare-dns-ip = pkgs.buildGo119Module rec {
              noCheck = true;
              pname = "cloudflare-dns-ip";
              version = "v0.1";
              src = ./.;
              subPackages = [
                "cmd/cli"
              ];
              postInstall = ''
                mv $out/bin/{cli,${pname}}
              '';
              checkPhase = false;
              vendorSha256 = "sha256-JFvC9V0xS8SZSdLsOtpyTrFzXjYAOaPQaJHdcnJzK3s=";
            };
          }
        );

      # Add dependencies that are only needed for development
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
          baseDeps = with pkgs; [ go gopls gotools go-tools ];
        in
        {
          default = pkgs.mkShell {
            buildInputs = baseDeps;
          };
        });

      # The default package for 'nix build'. This makes sense if the
      # flake provides only one package or there is a clear "main"
      # package.
      defaultPackage = forAllSystems (system: self.packages.${system}.cloudflare-dns-ip);

      overlay = final: prev: {
        cloudflare-dns-ip = self.packages.${final.system}.default;
      };

      nixosModules.default = { config, lib, pkgs, ... }:
        let
          cfg = config.services.cloudflare-dns-ip;
        in
        {
          options = with lib;{
            services.cloudflare-dns-ip = {
              enable = mkEnableOption "Enable cloudflare-dns-ip";
              tokenPath = mkOption {
                type = types.path;
                default = "/var/lib/cloudflare-dns-ip/token";
                description = "path to file containing the cloudflare token";
              };
              user = mkOption {
                type = types.str;
                default = "cloudflare-dns-ip";
                description = "user cloudflare-dns-ip should run as";
              };
              group = mkOption {
                type = types.str;
                default = "cloudflare-dns-ip";
                description = "group cloudflare-dns-ip should run as";
              };
              zone = mkOption {
                type = types.str;
                description = "the zone where the dns record is defined in cloudflare";
              };
              record = mkOption {
                type = types.str;
                description = "the dns record to keep updated";
              };
            };
          };
          config = lib.mkIf cfg.enable {
            users.users."${cfg.user}" = {
              createHome = false;
              group = "${cfg.group}";
              isSystemUser = true;
              isNormalUser = false;
              description = "user for cloudflare-dns-ip service";
            };
            users.groups."${cfg.group}" = { };
            systemd.services.cloudflare-dns-ip = {
              enable = true;
              script =
                ''
                  ${lib.optionalString (cfg.tokenPath != null) ''
                    export CLOUDFLARE_TOKEN="$(head -n1 ${lib.escapeShellArg cfg.tokenPath})"
                  ''}

                  ${pkgs.cloudflare-dns-ip}/bin/cloudflare-dns-ip update -t $CLOUDFLARE_TOKEN -z ${cfg.zone} -r ${cfg.record}
                '';
              wantedBy = [ "multi-user.target" ];
              after = [ "network-online.target" ];
              serviceConfig = {
                Type = "oneshot";
                RemainAfterExit = "yes";
                User = cfg.user;
                Group = cfg.group;
              };
            };
          };

        };

      formatter = forAllSystems (system: nixpkgsFor.${system}.nixpkgs-fmt);
    };
}
