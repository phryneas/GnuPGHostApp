import "babel-polyfill";
import ListenerQueue from './ListenerQueue';
import {Validity, Key, SubKey, UserID, EncryptedData, DecryptedData} from './HostAppTypes';

const listenerQueue = Symbol('ListenerQueue');
const port = Symbol('Port');

class NativeOpenGpgMeClient {

    constructor(runtime) {
        this[listenerQueue] = new ListenerQueue();
        this[port] = runtime.connectNative('de.phryneas.gpg.hostapp');
        console.info("hostApp connected");

        this[port].onDisconnect.addListener(function () {
            console.info("hostApp disconnected");
        });

        this[port].onMessage.addListener(this[listenerQueue].listener.bind(this[listenerQueue]));
    }

    /**
     * @param {string|Uint8Array} data
     * @param {Array.<string|Key>} publicKeys
     * @param {Array.<string|Key>} privateKeys
     * @param {boolean} armor
     * @param {boolean} detached
     * @param {Uint8Array} signature
     * @returns {Promise.<EncryptedData>}
     */
    encrypt({data, publicKeys, privateKeys, armor, detached, signature} = {}) {
        let isBytes = data instanceof Uint8Array;
        return this.sendToHostApp("encrypt", {
            dataString: !isBytes ? data : "",
            dataBytes: isBytes ? Array.from(data) : null,
            publicKeys: Key.wrapArray(publicKeys).map(key => key.fingerPrint),
            privateKeys: Key.wrapArray(privateKeys).map(key => key.fingerPrint),
            armor,
            detached,
            signature
        }).then(response => EncryptedData.wrap(response.data.encrypt));
    }

    /**
     * @param {string|Uint8Array} data
     * @param publicKeys
     * @param format
     * @param signature
     * @returns {Promise.<DecryptedData>}
     */
    decrypt({data, publicKeys, privateKey, format, signature} = {}) {
        let isBytes = data instanceof Uint8Array;

        return this.sendToHostApp("decrypt", {
            dataString: !isBytes ? data : "",
            dataBytes: isBytes ? Array.from(data) : null,
            publicKeys: Key.wrapArray(publicKeys).map(key => key.fingerPrint),
            format: format || 'utf8',
            signature
        }).then(response => DecryptedData.wrap(response.data.decrypt));
    }

    /**
     * @param keyID
     * @param fingerPrint
     * @param UID
     * @param name
     * @param comment
     * @param email
     * @param secretOnly
     * @returns {Promise.<FindKeysData>}
     */
    findKeys({keyID, fingerPrint, UID, name, comment, email, secretOnly = false} = {}) {
        return this.sendToHostApp("findKeys", {
            keyID, fingerPrint, UID, name, comment, email, secretOnly
        }).then(response => response.data.findKeys);
    }

    /**
     * @param {string} action
     * @param {Object} data
     * @returns {Promise.<HostResponse>}
     */
    sendToHostApp(action, data) {
        return new Promise((resolve, reject) => {
            this[listenerQueue].queueListener(/** @param {HostResponse} response */(response) => {
                console.info('received from HostApp', response);
                if (response.status === "ok") {
                    resolve(response);
                } else {
                    reject(response.message);
                }
            });

            let request = {action, data: {}};
            request.data[action] = data;

            console.info('sending to HostApp', request);
            this[port].postMessage(request);
        });
    }
}

export {NativeOpenGpgMeClient};
export * from './HostAppTypes';
export default NativeOpenGpgMeClient;