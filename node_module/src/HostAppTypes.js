/**
 * @typedef {Object} HostResponse
 * @property {string} status
 * @property {string} message
 * @property {HostResponseData} data
 */

/**
 * @typedef {Object} HostResponseData
 * @property {EncryptedData} encrypt
 * @property {DecryptedData} decrypt
 * @property {FindKeysData} findKeys
 */

/**
 * @typedef {Object} EncryptedData
 * @property {string} data
 * @property {Uint8Array} message
 * @property {Uint8Array} signature
 */

/**
 * @typedef {Object} DecryptedData
 * @property {string} dataString
 * @property {Uint8Array} dataBytes
 * @property {Array.<{keyid:string, valid:boolean}>} signatures
 */

/** @typedef {string} FingerPrint */

/**
 * @typedef {Object} FindKeysData
 * @property {Object<FingerPrint, Key>} keys
 */


/** @enum {int} */
export const Validity = {
    ValidityUnknown: 0,
    ValidityUndefined: 1,
    ValidityNever: 2,
    ValidityMarginal: 3,
    ValidityFull: 4,
    ValidityUltimate: 5
};


export class Wrappable {
    constructor(data) {
        Object.assign(this, data);
    }

    /**
     * @param {(Wrappable|Object)} item
     * @returns {Wrappable}
     */
    static wrap(item = {}) {
        return typeof item instanceof this ? item : new this(item);
    }

    /**
     * @param {Array.<Wrappable|Object>} arr
     * @returns {Array.<Wrappable>}
     */
    static wrapArray(arr = []) {
        if (!Array.isArray(arr)){
            return this.wrapArray([arr]);
        }
        return arr.filter(item => typeof item !== "undefined").map(this.wrap.bind(this));
    }
}


/**
 * @property {boolean} revoked
 * @property {boolean} expired
 * @property {boolean} disabled
 * @property {boolean} secret
 * @property {boolean} canEncrypt
 * @property {boolean} canSign
 * @property {boolean} canCertify
 * @property {boolean} canAuthenticate
 * @property {Validity} ownerTrust
 * @property {Array.<SubKey>} subKeys
 * @property {Array.<UserID>} userIDs
 */
export class Key extends Wrappable {
    /**
     *
     * @param {Array.<Object|SubKey>} subKeys
     * @param {Array.<Object|UserID>} userIDs
     */
    constructor({subKeys, userIDs} = {}) {
        super();
        if (this.constructor === Key) {
            this.subKeys = SubKey.wrapArray(subKeys);
            this.userIDs = UserID.wrapArray(userIDs);
        }
    }

    /**
     * @returns {string|null}
     */
    get fingerPrint() {
        if (!this.subKeys || this.subKeys.length === 0)
            return null;
        return this.subKeys[0].fingerPrint;
    }

    /**
     * @returns {string|null}
     */
    get UID() {
        if (!this.userIDs || this.userIDs.length === 0)
            return null;
        return this.userIDs[0].UID;
    }
}

/**
 * @property {boolean} revoked
 * @property {boolean} expired
 * @property {boolean} disabled
 * @property {boolean} invalid
 * @property {boolean} secret
 * @property {string} keyID
 * @property {string} fingerPrint
 * @property {Date} created
 * @property {Date} expires
 * @property {string} cardNumber
 */
export class SubKey extends Wrappable {
}

/**
 * @property {boolean} revoked
 * @property {boolean} invalid
 * @property {Validity} validity
 * @property {string} UID
 * @property {string} name
 * @property {string} comment
 * @property {string} email
 */
export class UserID extends Wrappable {
}
