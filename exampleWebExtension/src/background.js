//var openpgp = require("openpgp");

const lib = require("../../node_module/dist/module.js");
const NativeOpenGpgMeClient = lib.NativeOpenGpgMeClient;

console.log(lib);
let browser = window.browser || window.chrome;
const gpg = new NativeOpenGpgMeClient(browser.runtime);


browser.browserAction.onClicked.addListener(() =>
    gpg.findKeys({email: "gnupghostapp_tests@example.com"})
        .then(result => Object.values(result.keys)[0])
        .then(key => gpg.encrypt({
                data: "das ist ein Test",
                publicKeys: [key],
                privateKeys: [key],
                armor: true
            })
        )
        .then(result => result.data)
        .then(encrypted => gpg.decrypt({
                message: encrypted,
                publicKeys: ["1E43F132357B5AD55CECCCC3067D1766157F6495"]
            })
        )
        .catch(error => {
            console.error(error)
        })
);