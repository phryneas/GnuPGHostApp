export default class ListenerQueue {
    private queue: Function[] = [];

    listener(...args: any[]) {
        if (this.queue.length === 0)
            return;
        let callback = this.queue.shift();
        callback(...args);
    }

    queueListener(listener: Function) {
        this.queue.push(listener);
    }
}