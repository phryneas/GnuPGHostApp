import {NativeOpenGpgMeClient} from "../src/NativeOpenGpgMeClient";

let mockRuntime = {
    connectNative: () => {
        let listeners = {onDisconnect: [], onMessage: []};
        let triggerListeners = (listenerType, data) => listeners[listenerType].forEach(listener => listener(data));
        return {
            onDisconnect: {
                addListener: () => {
                }
            },
            onMessage: {
                addListener: (listener) => {
                    listeners["onMessage"].push(listener);
                }
            },
            postMessage: (message) => {
                switch (message.action) {
                    case "encrypt":
                        return triggerListeners("onMessage", {
                            "status": "ok",
                            "message": "",
                            "data": {
                                "encrypt": {
                                    "data": "-----BEGIN PGP MESSAGE-----\r\n\r\nhIwDRsAl1a6NFu4BBACD6wg94HlfAkvCGbYZXN7b6QtUkY8r6rtt3WktsAyeXsDh\r\n+cEUDgeCIqUWnm3WaQG3wUylO631GHJ/Set9GeGmXOCOloikMundEwR+jFa675a7\r\ndr7F+3oo0PDPKsY03a3FhUlVAVY+C+BniOKZZNxbzIaT3/a5gM6ijfPlGxRcp9JI\r\nAS4kK/xbkMO4akzoDOM6+IqX57TFviheNoIcPrOU13IbL76mgaxcOGt8bu2B2cFV\r\nyfJ3yGV3Hm9dg7LGmgdF67okMsT4hk3p\r\n=JFiG\r\n-----END PGP MESSAGE-----\r\n",
                                    "message": null,
                                    "signature": null
                                }
                            }
                        });
                    case "decrypt":
                        return triggerListeners("onMessage", {
                            "status": "ok",
                            "message": "",
                            "data": {
                                "decrypt": {
                                    "signatures": null,
                                    "dataString": "Hello, World!",
                                    "dataBytes": null
                                }
                            }
                        });
                    default:
                        return triggerListeners("onMessage", {
                            "status": "error",
                            "message": "unknown action",
                            "data": {}
                        });
                }
            }
        };
    }
};

describe("NativeOpenGpgMeClient", () => {

    it("contains specs with an expectation", () => {
        expect(true).toBe(true);
    });

    /** @var {NativeOpenGpgMeClient} gpg */
    let gpg;
    /** @var {string} encrypted */
    let encrypted;
    let testString = "Hello, World!";
    let testKeyFingerPrint = "1E43F132357B5AD55CECCCC3067D1766157F6495";

    it("constructs", () => {
        gpg = new NativeOpenGpgMeClient(mockRuntime);
        console.log(gpg);
        expect(gpg).not.toBeNull();
    });

    it("encrypts armored", (finished) => {
        gpg.encrypt({
            data: testString,
            public_keys: [testKeyFingerPrint],
            armor: true
        }).then((result) => {
            encrypted = result.data;
            expect(result.data).toContain("-----BEGIN PGP MESSAGE-----");
        }).then(finished);
    });

    it("decrypts armored", (finished) => {
        gpg.decrypt({
            message: encrypted,
            public_keys: [testKeyFingerPrint],
            returnFormat: "utf8"
        }).then((result) => {
            expect(result.data_string).toBe(testString);
        }).then(finished);
    });
});

