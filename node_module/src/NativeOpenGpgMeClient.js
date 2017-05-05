import "babel-polyfill";
import ListenerQueue from './ListenerQueue';
import {Validity, Key, SubKey, UserID} from './HostAppTypes';

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
     * @param {Array.<string|Key>} public_keys
     * @param {Array.<string|Key>} private_keys
     * @param {boolean} armor
     * @param {boolean} detached
     * @param {Uint8Array} signature
     * @returns {Promise.<EncryptedData>}
     */
    encrypt({data, public_keys, private_keys, armor, detached, signature} = {}) {
        let is_bytes = data instanceof Uint8Array;

        return this.sendToHostApp("encrypt", {
            data_string: !is_bytes ? data : "",
            data_bytes: is_bytes ? data : null,
            public_keys,
            private_keys,
            armor,
            detached,
            signature
        }).then(response => response.data.encrypt);
    }

    /**
     * @param {string|Uint8Array} data TODO: Uint8Array
     * @param public_keys
     * @param format
     * @param signature
     * @returns {Promise.<DecryptedData>}
     */
    decrypt({message, public_keys, private_key, format, signature} = {}) {
        let is_bytes = message instanceof Uint8Array;

        return this.sendToHostApp("decrypt", {
            message,
            public_keys,
            format,
            signature
        }).then(response => response.data.decrypt);
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