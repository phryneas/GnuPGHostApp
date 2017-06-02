import "babel-polyfill";
import ListenerQueue from './ListenerQueue';
import {Key, EncryptedData, DecryptedData, HostResponse, HostRequest} from './HostAppTypes';

import GenericData = HostRequest.GenericData;

class NativeOpenGpgMeClient {

    private listenerQueue: ListenerQueue;
    private port: chrome.runtime.Port;

    constructor(runtime: typeof chrome.runtime) {
        this.listenerQueue = new ListenerQueue();
        this.port = runtime.connectNative('de.phryneas.gpg.hostapp');
        console.info("hostApp connected");

        this.port.onDisconnect.addListener(function () {
            console.info("hostApp disconnected");
        });

        this.port.onMessage.addListener(this.listenerQueue.listener.bind(this.listenerQueue));
    }

    encrypt({
                data,
                publicKeys = [],
                privateKeys = [],
                armor = true,
                detached = false,
                signature = null
            }: {
                data: string | Uint8Array,
                publicKeys: (string | Key)[],
                privateKeys: (string | Key)[],
                armor: boolean,
                detached: boolean,
                signature?: Uint8Array
            }) {
        return this.sendToHostApp({
            action: "encrypt",
            data: {
                encrypt: {
                    dataString: typeof data === "string" ? data : "",
                    dataBytes: data instanceof Uint8Array ? (<any>Array).from(data) : null, // TODO better typing
                    publicKeys: publicKeys.map(key => typeof key === "string" ? key : key.fingerPrint),
                    privateKeys: privateKeys.map(key => typeof key === "string" ? key : key.fingerPrint),
                    armor,
                    detached,
                    signature
                }
            }
        }).then(response => new EncryptedData(response.data.encrypt));
    }

    decrypt({
                data,
                publicKeys,
                format,
                signature
            }: {
                data: string | Uint8Array,
                publicKeys: (string | Key)[],
                format: HostRequest.DataType,
                signature: string
            }) {

        return this.sendToHostApp({
            action: "decrypt",
            data: {
                decrypt: {
                    dataString: typeof data === "string" ? data : "",
                    dataBytes: data instanceof Uint8Array ? (<any>Array).from(data) : null, // TODO better typing
                    publicKeys: publicKeys.map(key => typeof key === "string" ? key : key.fingerPrint),
                    format,
                    signature
                }
            }
        }).then(response => new DecryptedData(response.data.decrypt));
    }

    findKeys({keyID, fingerPrint, UID, name, comment, email, secretOnly = false}: HostRequest.FindKeysData): Promise<HostResponse.FindKeysData> {
        return this.sendToHostApp({
            action: "findKeys",
            data: {
                findKeys: {
                    keyID,
                    fingerPrint,
                    UID, name,
                    comment,
                    email,
                    secretOnly
                }
            }
        }).then(response => response.data.findKeys);
    }

    sendToHostApp(request: HostRequest.HostRequest): Promise<HostResponse.HostResponse> {
        return new Promise((resolve, reject) => {
            this.listenerQueue.queueListener((response: HostResponse.HostResponse) => {
                console.info('received from HostApp', response);
                if (response.status === "ok") {
                    resolve(response);
                } else {
                    reject(response.message);
                }
            });

            console.info('sending to HostApp', request);
            this.port.postMessage(request);
        });
    }
}

export {NativeOpenGpgMeClient};
export * from './HostAppTypes';
export default NativeOpenGpgMeClient;