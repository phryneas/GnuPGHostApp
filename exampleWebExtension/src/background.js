//var openpgp = require("openpgp");

const lib = require("../../node_module/dist/module.js");
const NativeOpenGpgMeClient = lib.NativeOpenGpgMeClient;

console.log(lib);
let browser = window.browser || window.chrome;
const gpg = new NativeOpenGpgMeClient(browser.runtime);

browser.browserAction.onClicked.addListener(() =>
    gpg.encrypt({
        data: "das ist ein Test",
        publicKeys: ["1E43F132357B5AD55CECCCC3067D1766157F6495"],
        privateKeys: ["1E43F132357B5AD55CECCCC3067D1766157F6495"],
        armor: true
    }).then(result => {
        console.log(result);
        return result.data;
    }).then(encrypted => gpg.decrypt({
            message: encrypted,
            publicKeys: ["1E43F132357B5AD55CECCCC3067D1766157F6495"]
        })
    ).then(result => {
        console.log(result)
    }).catch(error => {
        console.error(error)
    })
);