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
            cloudflare-dns-ip = pkgs.buildGoModule {
              noCheck = true;
              pname = "cloudflare-dns-ip";
              version = "v0.1";
              src = ./.;
              subPackages = [
                "github.com/cloudflare-dns/cmd/cli"
              ];
              vendorSha256 = "";
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

      overlay = self.overlays.default;

      formatter = forAllSystems (system: nixpkgsFor.${system}.nixpkgs-fmt);
    };
}
