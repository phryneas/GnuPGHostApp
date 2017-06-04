import {HostResponse} from './HostAppTypes';

export type listenerFunction = (response: HostResponse.HostResponse) => void ;

export default class ListenerQueue {
    private queue: listenerFunction[] = [];

    listener(response: HostResponse.HostResponse) {
        if (this.queue.length === 0)
            return;
        let callback = this.queue.shift();
        callback(response);
    }

    queueListener(listener: listenerFunction) {
        this.queue.push(listener);
    }
}