export default class Encodings {
    static base64ToUint8Array(from: string): Uint8Array {
        let raw = window.atob(from);
        let array = new Uint8Array(new ArrayBuffer(raw.length));
        for (let i = 0; i < raw.length; i++) {
            array[i] = raw.charCodeAt(i);
        }
        return array;
    }

    static uint8ArrayToBase64(from: Uint8Array): string {
        let decoder = new TextDecoder('utf8');
        return btoa(decoder.decode(from));
    }
}
