//var openpgp = require("openpgp");

var openpgp = require("../../node_module/dist/src/opengpg-dropin").default;

console.log('loaded', openpgp);

var browser = browser || chrome;

var port = browser.runtime.connectNative('de.phryneas.gpg.hostapp');
port.onMessage.addListener(function (msg) {
    console.log("Received", msg);
});
port.onDisconnect.addListener(function () {
    console.log("Disconnected");
});

browser.browserAction.onClicked.addListener(function () {
    console.log('clicked');

    var options, encrypted;

    var pubkey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1

mI0EWQROJwEEAMtfVOyRBvd1BToqLxBnpbFRoZ3JebXjevOrojC21t9Qcxf7oTSi
LevdzvuUQUUF9zoJY0G83UEYLAQHwzaqVYRi9kEfkAX4jYfvuRxCCBNwYtsXVC8b
XfQQLVJhVjuhPKtV2AYmsGyJgpVvKjU3HkC7CRr+6PKhugenfaMaXqblABEBAAG0
QlRlc3QgS2V5IChvbmx5IGZvciB0ZXN0IHB1cnBvc2VzKSA8Z251cGdob3N0YXBw
X3Rlc3RzQGV4YW1wbGUuY29tPoi4BBMBAgAiBQJZBE4nAhsDBgsJCAcDAgYVCAIJ
CgsEFgIDAQIeAQIXgAAKCRAGfRdmFX9klfwNBACc2LaE6OFyBiXra405jujKociE
TNWYveuIB7p6mGnh+ssoswWPKd02uO5OxQayBbM3WA0mDPe3PtBXwbjFG6QnSv5C
eVZejtQvax06uyw48jd1naqz609iNx/buc5NP6rQ50WzmaPk6C3anPd3nICOZufz
TuQd0ZILly1xS8bRH7iNBFkETicBBADB7ZHODrmDqJ5mY4ybQI7FN1bTdh24Hpje
FcHTRrvQApN4Ttm/IM07cvKWUppQwuzMwvpjxPPloB/oImpr36wDPmDN6lotsQK8
W5HHKTEAUoDJoLGXgVuafnetr+q8hfvi/jsuw1GKGU2cJkQdm9Bw7z1ppmlTLprh
TbY3s7GsvQARAQABiJ8EGAECAAkFAlkETicCGwwACgkQBn0XZhV/ZJU8cwP/bTIo
cM3fr3iWU3bdo1zmXnzSf7kIsrHUTfZqshfJyDIYJgQaTdGav9Uq/Ncwjxlrnw40
DuIouEvacGLPXUUnDMPXKBPPRwNvVdrKx1fZVH4jERI16P0Fjq2u2Gisvb3WBJMM
lyEL7Mb18KJLTFMGJc3P6nu61b4wLrKccOvOKjc=
=xXDW
-----END PGP PUBLIC KEY BLOCK-----
`;
    var privkey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: GnuPG v1

lQHYBFkETicBBADLX1TskQb3dQU6Ki8QZ6WxUaGdyXm143rzq6IwttbfUHMX+6E0
oi3r3c77lEFFBfc6CWNBvN1BGCwEB8M2qlWEYvZBH5AF+I2H77kcQggTcGLbF1Qv
G130EC1SYVY7oTyrVdgGJrBsiYKVbyo1Nx5Auwka/ujyoboHp32jGl6m5QARAQAB
AAP7B8AKyvcd5llByT0pTP0+Lbs4JvyyFDHmkhmk1Slql9kHgc73jjtt95Kc3C6C
rEA1czM/YpZxchUbPE4VbORh3Ne7AUTNoI1r1Lz0NyCLxJ+iCW2LUp8WH0L4jiAR
HkOejCxdWDIhQhJ1iLSjIB1FKiPdsxxaM+h44Rn3KLGpmWMCANHix+PvvMm0MvaK
IDZZlOiZrhrs9iLG14cbpt4idQjUEJraVweEsDRZZIsVY+mYILC9NfIR91PvUdqW
cIFv/jcCAPgOMdwNnp0nUeubUkZS5BG83oMQ2OJfnzdxzrWo5vY3u2LwgZz/VNWI
zxYEcVHnJyk3pyIgztniH9rrVSi2lcMB/Ry5r2WCdkj8h9w8C3FLH8YWSloGhIUN
OfS5V/A0FidK2g/SO15SHqcc72tSWiIfO5teVnHcPpQqXEKVkV+P0GCbBbRCVGVz
dCBLZXkgKG9ubHkgZm9yIHRlc3QgcHVycG9zZXMpIDxnbnVwZ2hvc3RhcHBfdGVz
dHNAZXhhbXBsZS5jb20+iLgEEwECACIFAlkETicCGwMGCwkIBwMCBhUIAgkKCwQW
AgMBAh4BAheAAAoJEAZ9F2YVf2SV/A0EAJzYtoTo4XIGJetrjTmO6MqhyIRM1Zi9
64gHunqYaeH6yyizBY8p3Ta47k7FBrIFszdYDSYM97c+0FfBuMUbpCdK/kJ5Vl6O
1C9rHTq7LDjyN3WdqrPrT2I3H9u5zk0/qtDnRbOZo+ToLdqc93ecgI5m5/NO5B3R
kguXLXFLxtEfnQHYBFkETicBBADB7ZHODrmDqJ5mY4ybQI7FN1bTdh24HpjeFcHT
RrvQApN4Ttm/IM07cvKWUppQwuzMwvpjxPPloB/oImpr36wDPmDN6lotsQK8W5HH
KTEAUoDJoLGXgVuafnetr+q8hfvi/jsuw1GKGU2cJkQdm9Bw7z1ppmlTLprhTbY3
s7GsvQARAQABAAP5AdWph3WEM8aomPdgISffMeZwH9gCN/eyIoe6KbGFnVYo5v53
+OLqjiFsQhfN9e2iJ93AWKlIVWfKZXvN3e9jxS/d0UYBCUOcjVvjbpXRJd0WQL8S
3d8MVs3Vga4t2RJlNRLfXk7lwbyerOvXkschtr5WYm7BEdF9sNaNhOSI1dECANZO
UUSf57lrsazt38BGjFdidKG4Sf90FHIzilwejT3+tp+vaj3Dj/nvtPgtQufJvuhV
5rrTqNspr0yXA3Swj8cCAOeoUCClGsH9CoARimdf2XVKKyhshoXyUfNdHNnF4zTh
UKv9Pc4y4F31NxHOwM/mil35XHaGkH6WOIKgQkM+51sB/1GlwLiZ3hzrEJDQFver
sTOv4h4FpxGLoLmgCTFhDwyrtRs4vVohkZVJWz3Jr5SjNsI1dRqsuIbqpfG4ZwQf
QkarkIifBBgBAgAJBQJZBE4nAhsMAAoJEAZ9F2YVf2SVPHMD/20yKHDN3694llN2
3aNc5l580n+5CLKx1E32arIXycgyGCYEGk3Rmr/VKvzXMI8Za58ONA7iKLhL2nBi
z11FJwzD1ygTz0cDb1XaysdX2VR+IxESNej9BY6trthorL291gSTDJchC+zG9fCi
S0xTBiXNz+p7utW+MC6ynHDrzio3
=ARxr
-----END PGP PRIVATE KEY BLOCK-----`;
    var passphrase = '';

    var privKeyObj = openpgp.key.readArmored(privkey).keys[0];
    privKeyObj.decrypt(null);


    options = {
        data: 'Hello, World!',                             // input as String (or Uint8Array)
        publicKeys: openpgp.key.readArmored(pubkey).keys,  // for encryption
        privateKeys: privKeyObj // for signing (optional)
    };

    openpgp.encrypt(options)
        .then(function (ciphertext) {
            console.log("ciphertext", ciphertext);
            return ciphertext.data;
        })
        .then(function (encrypted) {
            options = {
                message: encrypted, //TODO openpgp.message.readArmored(encrypted),     // parse armored message
                publicKeys: openpgp.key.readArmored(pubkey).keys,    // for verification (optional)
                privateKey: privKeyObj // for decryption
            };

            console.log(options);
            return openpgp.decrypt(options);

        })
        .then(function (plaintext) {
            console.log("plaintext", plaintext);
            return plaintext.data;
        });


});