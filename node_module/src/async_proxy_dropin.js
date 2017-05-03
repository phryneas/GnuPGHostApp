import openpgp from './index';

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
 * @property {ListenerQueue} _listenerQueue
 * @property {Port} _port
 */
export default class AsyncProxy {
    constructor() {
        let browser = window.browser || window.chrome;

        this._listenerQueue = new ListenerQueue();
        this._port = browser.runtime.connectNative('de.phryneas.gpg.hostapp');
        console.info("hostApp connected");

        this._port.onDisconnect.addListener(function () {
            console.info("hostApp disconnected");
        });

        this._port.onMessage.addListener(this._listenerQueue.listener.bind(this._listenerQueue));
    }

    delegate(method, options) {
        if (typeof openpgp[method] !== 'function') {
            return Promise.resolve({event: 'method-return', err: 'Unknown Worker Event'});
        }
        return openpgp[method](options);
    }

    sendToHostApp(data){
        return new Promise((resolve) => {
            this._listenerQueue.queueListener(function (response) {
                console.info('received from HostApp', data);
                resolve(response);
            });

            console.info('sending to HostApp', data);
            this._port.postMessage(data);
        });
    }
}
