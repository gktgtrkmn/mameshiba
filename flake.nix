{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAll = f: nixpkgs.lib.genAttrs systems (s: f nixpkgs.legacyPackages.${s});
    in
    {
      packages = forAll
        (pkgs: rec {
          mameshiba = pkgs.buildGoModule
            {
              pname = "mameshiba";
              version = "0.1.0";
              src = ./.;
              vendorHash = "sha256-nssNEysi9/p+i6IndbEbK6qmDydlrC5Hx/NE+exILE4=";
            };
          default = mameshiba;
        });

      devShells = forAll (pkgs: {
        default = pkgs.mkShell {
          packages = [ pkgs.go pkgs.gopls pkgs.gotools ];
        };
      });
    };
}
