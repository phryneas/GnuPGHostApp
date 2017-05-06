REM ./include should contain header files from C:\Program Files (x86)\GnuPG\include
REM ./lib should contain files from C:\Program Files (x86)\GnuPG\lib, with extension .imp changed to .a (libgpg-err.a, libassuan.a, libgpgme.a)
REM ./ should contain dll files from C:\Program Files (x86)\GnuPG\bin (libassuan-0.dll  libgpg-error-0.dll  libgpgme-11.dll)

set PATH=%PATH%;C:\Go\bin
set PATH=%PATH%;C:\mingw\bin
set PATH=%PATH%;C:\mingw\mingw32\bin
set GOPATH=%~dp0

set CGO_CFLAGS=-I%~dp0\include
set CGO_LDFLAGS=-L%~dp0\lib

REM go test -v github.com/proglottis/gpgme
go test -v HostApp NativeMessagingHost OpenPgpJsApi
go build HostApp

echo Build complete

pause
