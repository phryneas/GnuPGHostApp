console.log('loaded');

var browser = browser || chrome;

var port = browser.runtime.connectNative('de.phryneas.gpg.hostapp');
port.onMessage.addListener(function(msg) {
    console.log("Received" + msg);
});
port.onDisconnect.addListener(function() {
    console.log("Disconnected");
});

browser.browserAction.onClicked.addListener(function() {
    console.log('clicked');
    port.postMessage({ action: "test" });
});