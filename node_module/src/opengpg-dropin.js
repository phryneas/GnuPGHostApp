let port, listenerQueue;

/**
 * @property {function[]} queue
 */
class ListenerQueue {
    constructor(){
        this.queue = [];
    }

    listener(...args){
        if (this.queue.length == 0)
            return;
        let callback = this.queue.shift();
        callback(...args);
    }

    queueListener(listener) {
        this.queue.push(listener);
    }
}


export function initWorker() {
    if (port){
        return true;
    }

    let browser = browser || chrome;
    if (!browser || !browser.runtime || typeof browser.runtime.connectNative != "function") {
        return false;
    }

    listenerQueue = new ListenerQueue()
    port = browser.runtime.connectNative('de.phryneas.gpg.hostapp')
    console.info("hostApp connected")

    if (!port){
        return false;
    }

    port.onDisconnect.addListener(function () {
        console.info("hostApp disconnected");
    });

    port.onMessage.addListener(listenerQueue.listener.bind(listenerQueue));

    return true;
}


function convertKeys(keys){
    if (!Array.isArray(keys)){
        if (!keys){
            keys = [];
        } else {
            keys = [ keys ];
        }
    }
    return keys.map((key) => key.primaryKey.fingerprint);
}

export function encrypt({data, publicKeys = [], privateKeys = [], passwords = [], filename, armor=true, detached=false, signature=null}) {
    return new Promise(function(resolve){
        listenerQueue.queueListener(function(response){
            console.log("listener received:", response);
            resolve({data: response.data.encrypt.data});
        });

        let message = {
            "action": "encrypt",
            "data": {
                "encrypt": {
                    "data_string": null,
                    "data_bytes": null,
                    "public_keys": convertKeys(publicKeys),
                    "private_keys": convertKeys(privateKeys),
                    "armor": armor,
                    "detached": detached,
                    "signature": null //TODO
                }
            }
        };

        if (typeof data == "string"){
            message.data.encrypt.data_string = data;
        } else {
            message.data.encrypt.data_bytes = data;
        }

        port.postMessage(message);
    });
}

export function decrypt({ message, privateKey, publicKeys, sessionKey, password, format='utf8', signature=null }) {
    return new Promise(function(resolve){
        listenerQueue.queueListener(function(response){
            console.log("listener received:", response);
            resolve({data: response.data.decrypt.data});
        });

        let message = {
            "action": "decrypt",
            "data": {
                "decrypt": {
                    "message": format == 'utf-8' ? message : message, //TODO?
                    "public_keys": convertKeys(publicKeys),
                    "private_key": privateKey?privateKey.primaryKey.fingerprint:null,
                    "format": format,
                    "signature": null //TODO
                }
            }
        };

        port.postMessage(message);
    });
}