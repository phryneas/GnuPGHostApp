import Encodings from './Encodings';

export namespace HostResponse {
    export interface HostResponse {
        status: string;
        message?: string;
        data: HostResponse.Data;
    }

    export interface Data {
        encrypt?: EncryptedData;
        decrypt?: DecryptedData;
        findKeys?: FindKeysData;
    }

    export interface EncryptedData {
        dataString?: string;
        dataBytes?: string;
        signature?: string;
    }

    export interface DecryptedData {
        dataString: string;
        dataBytes: string;
        signatures: Signature[];
    }

    export interface SubKey {
        revoked: boolean;
        expired: boolean;
        disabled: boolean;
        invalid: boolean;
        secret: boolean;
        keyID: string;
        fingerPrint: string;
        created: Date;
        expires: Date;
        cardNumber: string;
    }

    export interface UserID {
        revoked: boolean;
        invalid: boolean;
        validity: boolean;
        UID: string;
        name: string;
        comment: string;
        email: string;
    }

    export interface Signature {
        keyId: string;
        valid: boolean;
    }

    export interface FindKeysData {
        [key: string]: Key
    }

    export interface Key {
        revoked: boolean;
        expired: boolean;
        disabled: boolean;
        secret: boolean;
        canEncrypt: boolean;
        canSign: boolean;
        canCertify: boolean;
        canAuthenticate: boolean;
        ownerTrust: Validity;
        subKeys: SubKey[];
        userIDs: UserID[];
    }

    export enum Validity{
        ValidityUnknown = 0,
        ValidityUndefined = 1,
        ValidityNever = 2,
        ValidityMarginal = 3,
        ValidityFull = 4,
        ValidityUltimate = 5
    }
}

export namespace HostRequest {
    export interface HostRequest {
        action: Action,
        data: {
            encrypt?: EncryptData,
            decrypt?: DecryptData,
            findKeys?: FindKeysData
        }
    }
    export type Action = "decrypt" | "encrypt" | "findKeys" | "test";
    export type GenericData = EncryptData | DecryptData | FindKeysData;
    export interface EncryptData {
        dataString?: string,
        dataBytes?: string,
        publicKeys: string[],
        privateKeys: string[],
        armor: boolean,
        detached?: boolean,
        signature?: Uint8Array
    }
    export interface DecryptData {
        dataString?: string,
        dataBytes?: string,
        publicKeys: string[],
        format: HostRequest.DataType,
        signature: string
    }
    export interface FindKeysData {
        keyID?: string,
        fingerPrint?: string,
        UID?: string,
        name?: string,
        comment?: string,
        email?: string,
        secretOnly: boolean
    }

    export type DataType = "utf8" | "binary";
}

export class EncryptedData {
    public dataString?: string;
    public dataBytes?: Uint8Array;
    public signature?: Uint8Array;

    constructor(data: HostResponse.EncryptedData) {
        this.dataString = data.dataString;
        if (data.dataBytes) {
            this.dataBytes = Encodings.base64ToUint8Array(data.dataBytes);
        }
        if (data.signature) {
            this.signature = Encodings.base64ToUint8Array(data.signature);
        }
    }
}


export class DecryptedData {
    public dataString: string;
    public dataBytes: Uint8Array;
    public signatures: HostResponse.Signature[];

    constructor(data: HostResponse.DecryptedData) {
        this.dataString = data.dataString;
        this.signatures = data.signatures;

        if (data.dataBytes) {
            this.dataBytes = Encodings.base64ToUint8Array(data.dataBytes);
        }
    }
}

export class FindKeysData {
    [key: string]: Key;

    constructor(data: HostResponse.FindKeysData) {
        for (let [key, value] of (<any>Object).entries(data)) {
            this[key] = new Key(value);
        }
    }
}

export class Key {
    revoked: boolean;
    expired: boolean;
    disabled: boolean;
    secret: boolean;
    canEncrypt: boolean;
    canSign: boolean;
    canCertify: boolean;
    canAuthenticate: boolean;
    ownerTrust: HostResponse.Validity;
    subKeys: HostResponse.SubKey[];
    userIDs: HostResponse.UserID[];

    constructor(data: HostResponse.Key) {
        this.revoked = data.revoked;
        this.expired = data.expired;
        this.disabled = data.disabled;
        this.secret = data.secret;
        this.canEncrypt = data.canEncrypt;
        this.canSign = data.canSign;
        this.canCertify = data.canCertify;
        this.canAuthenticate = data.canAuthenticate;
        this.ownerTrust = data.ownerTrust;
        this.subKeys = data.subKeys;
        this.userIDs = data.userIDs;
    }

    get fingerPrint(): string | null {
        if (!this.subKeys || this.subKeys.length === 0)
            return null;
        return this.subKeys[0].fingerPrint;
    }

    get UID(): string | null {
        if (!this.userIDs || this.userIDs.length === 0)
            return null;
        return this.userIDs[0].UID;
    }
}

