console.log('loaded');

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
    //port.postMessage({ action: "test" });

    port.postMessage({
        "action": "encrypt",
        "data": {
            "encrypt": {
                "data_string": "test",
                "data_bytes": null,
                "public_keys": ["1E43F132357B5AD55CECCCC3067D1766157F6495"],
                "private_keys": null,
                "armor": true,
                "detached": false,
                "signature": null
            }
        }
    });
});