with import <nixpkgs> {};

stdenv.mkDerivation {
  name = "native-opengpgme-client";

  buildInputs = [
    nodejs
  ];

shellHook =
  ''
'';
}
