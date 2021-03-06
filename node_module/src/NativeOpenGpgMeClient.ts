import ListenerQueue from './ListenerQueue';
import {Key, EncryptedData, DecryptedData, FindKeysData, HostResponse, HostRequest} from './HostAppTypes';
import Encodings from './Encodings';

export default class NativeOpenGpgMeClient {

    private listenerQueue: ListenerQueue = new ListenerQueue();
    private port: chrome.runtime.Port = null;

    constructor(private runtime: typeof chrome.runtime, private logger?: Console) {
    }

    public connect(application: string = 'de.phryneas.gpg.hostapp') {
        if (this.port != null) {
            this.port.disconnect();
        }
        this.port = this.runtime.connectNative(application);
        this.logger && this.logger.info("hostApp connected to %s", application);


        this.port.onDisconnect.addListener((p: chrome.runtime.Port) => {
            this.port = null;
            this.logger && this.logger.info("hostApp disconnected", (<any>p).error || chrome.runtime.lastError);
        });

        this.port.onMessage.addListener(this.listenerQueue.listener.bind(this.listenerQueue));
    }

    get connected() {
        return this.port != null;
    }

    public encrypt({
                       data,
                       encryptFor,
                       signWith = [],
                       armor = true,
                       detached = false,
                       signature = null
                   }: {
        data: string | Uint8Array,
        encryptFor: (string | Key)[],
        signWith?: (string | Key)[],
        armor?: boolean,
        detached?: boolean,
        signature?: Uint8Array
    }) {
        let encryptData: HostRequest.EncryptData = {
            dataString: typeof data === "string" ? data : null,
            dataBytes: data instanceof Uint8Array ? Encodings.uint8ArrayToBase64(data) : null,
            publicKeys: encryptFor.map(key => typeof key === "string" ? key : key.fingerPrint),
            privateKeys: signWith.map(key => typeof key === "string" ? key : key.fingerPrint),
            armor,
            detached,
            signature
        };

        return this.sendToHostApp({
            action: "encrypt",
            data: {encrypt: encryptData}
        }).then((response: HostResponse.HostResponse) => new EncryptedData(response.data.encrypt));
    }

    public decrypt({
                       data,
                       verifySignatures,
                       returnFormat = "utf8",
                       detachedSignature
                   }: {
        data: string | Uint8Array,
        verifySignatures: (string | Key)[],
        returnFormat?: HostRequest.DataType,
        detachedSignature?: string
    }) {

        let decryptData: HostRequest.DecryptData = {
            dataString: typeof data === "string" ? data : null,
            dataBytes: data instanceof Uint8Array ? Encodings.uint8ArrayToBase64(data) : null,
            publicKeys: verifySignatures.map(key => typeof key === "string" ? key : key.fingerPrint),
            format: returnFormat,
            signature: detachedSignature
        };

        return this.sendToHostApp({
            action: "decrypt",
            data: {decrypt: decryptData}
        }).then((response: HostResponse.HostResponse) => new DecryptedData(response.data.decrypt));
    }

    public findKeys({
                        keyID,
                        fingerPrint,
                        UID,
                        name,
                        comment,
                        email,
                        secretOnly = false
                    }: HostRequest.FindKeysData): Promise<FindKeysData> {
        let findKeysData: HostRequest.FindKeysData = {
            keyID,
            fingerPrint,
            UID,
            name,
            comment,
            email,
            secretOnly
        };

        return this.sendToHostApp({
            action: "findKeys",
            data: {findKeys: findKeysData}
        }).then((response: HostResponse.HostResponse) => new FindKeysData(response.data.findKeys));
    }

    public exportPublicKeys(pattern: string): Promise<HostResponse.ExportPublicKeysData> {
        return this.sendToHostApp({
            action: "exportPublicKeys",
            data: {exportPublicKeys: {pattern}}
        }).then((response: HostResponse.HostResponse) => response.data.exportPublicKeys);
    }

    sendToHostApp(request: HostRequest.HostRequest): Promise<HostResponse.HostResponse> {
        return new Promise((resolve: (response: HostResponse.HostResponse) => void, reject: (error: string) => void) => {
            this.listenerQueue.queueListener((response: HostResponse.HostResponse) => {
                this.logger && this.logger.info('received from HostApp', response);
                if (response.status === "ok") {
                    resolve(response);
                } else {
                    reject(response.message);
                }
            });

            this.logger && this.logger.info('sending to HostApp', request);
            this.port.postMessage(request);
        });
    }
}
