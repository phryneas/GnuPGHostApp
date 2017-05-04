//var openpgp = require("openpgp");

const lib = require("../../node_module/dist/src/NativeOpenGpgMeClient");
const NativeOpenGpgMeClient = lib.default;

console.log(lib);
let browser = window.browser || window.chrome;
const gpg = new NativeOpenGpgMeClient(browser.runtime);

browser.browserAction.onClicked.addListener(function () {
    gpg.encrypt({
        data: "das ist ein Test",
        public_keys: ["1E43F132357B5AD55CECCCC3067D1766157F6495"],
        private_keys: ["1E43F132357B5AD55CECCCC3067D1766157F6495"],
        armor: true
    }).then(result => {
        console.log(result);
        return result.data;
    }).then(encrypted => gpg.decrypt({
            message: encrypted,
            public_keys: ["1E43F132357B5AD55CECCCC3067D1766157F6495"]

        })
    ).then(result => {
        console.log(result)
    });
});
