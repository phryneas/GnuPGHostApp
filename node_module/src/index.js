import openpgp, {
    destroyWorker,
    generateKey,
    reformatKey,
    decryptKey,
    sign,
    verify,
    encryptSessionKey,
    decryptSessionKey,
    key,
    signature,
    cleartext,
    util,
    packet,
    MPI,
    S2K,
    Keyid,
    armor,
    enums,
    config,
    crypto,
    Keyring,
    AsyncProxy,
    HKP
} from 'openpgp';

import {
    initWorker,
    getWorker,
    encrypt,
    decrypt
} from './opengpg_dropin'

let overridden_openpgp = Object.assign(
    {},
    openpgp, {
        initWorker,
        getWorker,
        encrypt,
        decrypt
    }
);
export default overridden_openpgp;

export {
    initWorker, getWorker, destroyWorker, generateKey, reformatKey, decryptKey, encrypt, decrypt, sign,
    verify, encryptSessionKey, decryptSessionKey, key, signature, cleartext, util, packet, MPI, S2K, Keyid, armor,
    enums, config, crypto, Keyring, AsyncProxy, HKP
};
