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
        return item instanceof this ? item : new this(item);
    }

    /**
     * @param {Array.<Wrappable|Object>} arr
     * @returns {Array.<Wrappable>}
     */
    static wrapArray(arr = []) {
        return arr.map(this.wrap);
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
     * @param {string|Key} key
     */
    static getFingerprint(key) {
        if (key instanceof Key) {
            return key.fingerprint;
        }
        return key.toString();
    }

    /**
     * @returns {string|null}
     */
    get fingerprint() {
        for (let subKey of Object.values(this.subKeys)) {
            if (!(subKey.revoked || subKey.expired || subKey.disabled || subKey.invalid)) {
                return subKey.fingerprint;
            }
        }
        return null;
    }

    /**
     * @returns {string|null}
     */
    get userId() {
        for (let firstUserID of Object.values(this.userIDs)) {
            if (firstUserID  && !firstUserID.revoked && !firstUserID.invalid && firstUserID.name && firstUserID.email){
                return `${firstUserID.name}${firstUserID.comment ? ` (${firstUserID.comment}) ` : ' '}<${firstUserID.email}>`
            }
        }

        return null;
    }
}

/**
 * @property {boolean} invalid
 * @property {string} keyID
 * @property {string} fingerPrint
 * @property {Date} created
 * @property {Date} expires
 * @property {string} cardNumber
 */
export class SubKey extends Key {
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
