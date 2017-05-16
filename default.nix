with import <nixpkgs> {};

stdenv.mkDerivation {
  name = "gpgmeHostApp";

  buildInputs = [
	go
	godep
	gpgme
  ];

shellHook =
  ''
export GOROOT=$(cd $(dirname $(which go)); cd ../share/go; pwd)
export GOPATH=$(pwd)
export GNUPGHOME=$(cd test_keys/GPGHOME/; pwd)

go test HostApp NativeMessagingHost OpenPgpJsApi
  '';
}
