import AsyncWorker from './async_proxy_dropin';

import openpgp from 'openpgp';

/**
 *
 * @param {Object[]|Object|null} objects
 * @returns {Array}
 */
function toArray(objects) {
    if (!Array.isArray(objects)) {
        if (!objects) {
            objects = [];
        } else {
            objects = [objects];
        }
    }
    return objects;
}



let asyncWorker = new AsyncWorker();

export function initWorker() {
    return true;
}


export function getWorker() {
    return asyncWorker;
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
export function encrypt({data, publicKeys = [], privateKeys = [], passwords = [], filename, armor = true, detached = false, signature = null} = {}) {
    let request = {
        "action": "encrypt",
        "data": {
            "encrypt": {
                "data_string": null,
                "data_bytes": null,
                "public_keys": toArray(publicKeys).map((key) => key.primaryKey.fingerprint),
                "private_keys": toArray(privateKeys).map((key) => key.primaryKey.fingerprint),
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

    return this.getWorker()
        .sendToHostApp(request)
        .then(response => response.data.encrypt);
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
export function decrypt({message, privateKey, publicKeys, sessionKey, password, format = 'utf8', signature = null} = {}) {
    // TODO: wrap Message && this.message - this should not be handled via direct call.

    let request = {
        "action": "decrypt",
        "data": {
            "decrypt": {
                "message": message, //TODO
                "public_keys": toArray(publicKeys).map((key) => key.primaryKey.fingerprint),
                "private_key": privateKey ? privateKey.primaryKey.fingerprint : null,
                "format": format,
                "signature": null //TODO
            }
        }
    };

    return this.getWorker()
        .sendToHostApp(request)
        .then(response => response.data.decrypt);
}