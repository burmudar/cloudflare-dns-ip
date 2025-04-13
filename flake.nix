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
            cloudflare-dns-ip = pkgs.buildGo122Module rec {
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
              vendorHash = "sha256-JFvC9V0xS8SZSdLsOtpyTrFzXjYAOaPQaJHdcnJzK3s=";
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
        cloudflare-dns-ip = self.packages.${final.system}.cloudflare-dns-ip;
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
                description = "the zone where the dns record is defined in cloudflare eg. 'burmudar.dev'";
              };
              record = mkOption {
                type = types.str;
                description = "the dns record to keep updated eg. www";
              };
              ttl = mkOption {
                type = types.int;
                description = "TTL of the record measured in secs";
                default = 25 * 60; # 25 mins * 60 secs = 1800
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
            system.activationScripts.cloudflare-lib-dir = lib.stringAfter [ "var" ] ''
              mkdir -p /var/lib/cloudflare-dns-ip/
              chown ${cfg.user}:${cfg.group} /var/lib/cloudflare-dns-ip/
            '';
            systemd.services.cloudflare-dns-ip = {
              script = let
              token = cfg.tokenPath;
              zone = cfg.zone;
              record = cfg.record;
              ttl = toString cfg.ttl;
              in
                ''
                  ${pkgs.cloudflare-dns-ip}/bin/cloudflare-dns-ip update -t ${token} -z ${zone} -r ${record} --ttl ${ttl}
                '';
              # without 'wantedBy' this unit won't be automatically started at boot
              wants = [ "network-online.target" ];
              after = [ "network-online.target" ];
              serviceConfig = {
                Type = "oneshot";
                User = cfg.user;
                Group = cfg.group;
              };
            };
            systemd.timers.cloudflare-dns-ip = {
              enable = true;
              wantedBy = [ "timers.target" ];
              wants = [ "network-online.target" ];
              after = [ "network-online.target" ];
              timerConfig = {
                OnCalendar = "*-*-* *:30:00";
                Persistent = true;
              };
            };
          };

        };

      formatter = forAllSystems (system: nixpkgsFor.${system}.nixpkgs-fmt);
    };
}
