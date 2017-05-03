/*jshint esversion: 6 */

import * as openpgp from 'openPgp';
/*
import {Key} from 'openpgp/src/key';
import {Signature} from 'openpgp/src/signature';
import {Message} from 'openpgp/src/message';
*/

/**
 * @property {function[]} queue
 */
class ListenerQueue {
    constructor() {
        this.queue = [];
    }

    listener(...args) {
        if (this.queue.length === 0)
            return;
        let callback = this.queue.shift();
        callback(...args);
    }

    queueListener(listener) {
        this.queue.push(listener);
    }
}

/**
 *
 * @param {Key[]} keys
 * @returns {Array}
 */
function toArray(keys) {
    if (!Array.isArray(keys)) {
        if (!keys) {
            keys = [];
        } else {
            keys = [keys];
        }
    }
    return keys.map(/** @param {Key} key*/(key) => key.primaryKey.fingerprint);
}

/**
 * @param module
 * @returns {Function}
 */
function wrapModule(module) {
    let ret = function () {
    };
    for (let [key, value] of Object.entries(module)) {
        ret.prototype[key] = value;
    }
    return ret;
}


let wrapped = wrapModule(openpgp);

/**
 * @property {OpenPgpDropIn} delegateTo
 */
class Worker {
    constructor(delegateTo) {
        this.delegateTo = delegateTo;
    }

    delegate(method, options) {
        if (typeof this.delegateTo[method] !== 'function') {
            return Promise.resolve({event: 'method-return', err: 'Unknown Worker Event'});
        }
        return this.delegateTo[method](options);
    }
}

/**
 * @property {runtime.Port} _port
 * @property {ListenerQueue} _listenerQueue
 * @extends {openpgp}
 */
class OpenPgpDropIn extends wrapModule(openpgp) {
    constructor() {
        super();

        let browser = window.browser || window.chrome;

        this._listenerQueue = new ListenerQueue();
        this._port = browser.runtime.connectNative('de.phryneas.gpg.hostapp');
        console.info("hostApp connected");

        this._port.onDisconnect.addListener(function () {
            console.info("hostApp disconnected");
        });

        this._port.onMessage.addListener((data) => console.log(data));
        this._port.onMessage.addListener(this._listenerQueue.listener.bind(this._listenerQueue));
    }


    /**
     * @inheritDoc
     */
    initWorker() {
        return true;
    }

    /**
     * @inheritDoc
     */
    getWorker() {
        return new Worker(this);
    }

    /**
     * @see openpgp.encrypt
     * Encrypts message text/data with public keys, passwords or both at once. At least either public keys or passwords
     *   must be specified. If private keys are specified, those will be used to sign the message.
     * @param  {String|Uint8Array} data           text/data to be encrypted as JavaScript binary string or Uint8Array
     * @param  {Key|Array<Key>} publicKeys        (optional) array of keys or single key, used to encrypt the message
     * @param  {Key|Array<Key>} privateKeys       (optional) private keys for signing. If omitted message will not be signed
     * @param  {String|Array<String>} passwords   (optional) array of passwords or a single password to encrypt the message
     * @param  {String} filename                  (optional) a filename for the literal data packet
     * @param  {Boolean} armor                    (optional) if the return values should be ascii armored or the message/signature objects
     * @param  {Boolean} detached                 (optional) if the signature should be detached (if true, signature will be added to returned object)
     * @param  {Signature} signature              (optional) a detached signature to add to the encrypted message
     * @return {Promise<Object>}                  encrypted (and optionally signed message) in the form:
     *                                              {data: ASCII armored message if 'armor' is true,
     *                                                message: full Message object if 'armor' is false, signature: detached signature if 'detached' is true}
     */
    encrypt({data, publicKeys = [], privateKeys = [], passwords = [], filename, armor = true, detached = false, signature = null}={}) {
        return new Promise((resolve) => {
            this._listenerQueue.queueListener(function (response) {
                resolve(response.data.encrypt);
            });

            let request = {
                "action": "encrypt",
                "data": {
                    "encrypt": {
                        "data_string": null,
                        "data_bytes": null,
                        "public_keys": toArray(publicKeys),
                        "private_keys": toArray(privateKeys),
                        "armor": armor,
                        "detached": detached,
                        "signature": null //TODO
                    }
                }
            };

            if (typeof data === "string") {
                request.data.encrypt.data_string = data;
            } else {
                request.data.encrypt.data_bytes = data;
            }
            console.info('sending message', request);
            this._port.postMessage(request);
        });
    }

    /**
     * Decrypts a message with the user's private key, a session key or a password. Either a private key,
     *   a session key or a password must be specified.
     * @param  {Message} message             the message object with the encrypted data
     * @param  {Key} privateKey              (optional) private key with decrypted secret key data or session key
     * @param  {Key|Array<Key>} publicKeys   (optional) array of public keys or single key, to verify signatures
     * @param  {Object} sessionKey           (optional) session key in the form: { data:Uint8Array, algorithm:String }
     * @param  {String} password             (optional) single password to decrypt the message
     * @param  {String} format               (optional) return data format either as 'utf8' or 'binary'
     * @param  {Signature} signature         (optional) detached signature for verification
     * @return {Promise<Object>}             decrypted and verified message in the form:
     *                                         { data:Uint8Array|String, filename:String, signatures:[{ keyid:String, valid:Boolean }] }
     */
    decrypt({message, privateKey, publicKeys, sessionKey, password, format = 'utf8', signature = null}={}) {

         // TODO: wrap Message && this.message - this should not be handled via direct call.
        return new Promise((resolve) => {
            this._listenerQueue.queueListener(function (response) {
                resolve(response.data.decrypt);
            });

            let request = {
                "action": "decrypt",
                "data": {
                    "decrypt": {
                        "message": message, //TODO
                        "public_keys": toArray(publicKeys),
                        "private_key": privateKey ? privateKey.primaryKey.fingerprint : null,
                        "format": format,
                        "signature": null //TODO
                    }
                }
            };
            console.info('sending message', request);
            this._port.postMessage(request);
        });
    }
}

export default new OpenPgpDropIn();