import {default as NativeOpenGpgMeClient, FindKeysData, Key} from 'native-opengpgme-client';

let browser = (<any>window).browser || window.chrome;
const gpg = new NativeOpenGpgMeClient(browser.runtime, console);

function log<T>(x: T) : T {
    console.log(x);
    return x;
}

browser.browserAction.onClicked.addListener(() =>
    gpg.findKeys({email: "gnupghostapp_tests@example.com"})
        .then((result: FindKeysData) => result[Object.keys(result)[0]])
        .then((key: Key) => gpg.encrypt({
                data: "das ist ein Test",
                encryptFor:   [key],
                signWith: [key],
                armor: false
            })
        )
        .then(result => result.dataBytes || result.dataString)
        .then(log)
        .then(encrypted => gpg.decrypt({
                data: encrypted,
                verifySignatures: ["1E43F132357B5AD55CECCCC3067D1766157F6495"],
            })
        )
        .then(log)
        .catch(error => {
            console.error(error)
        })
);