/**
 * @property {function[]} queue
 */
export default class ListenerQueue {
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