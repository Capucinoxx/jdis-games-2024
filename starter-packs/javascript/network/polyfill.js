import crypto from 'crypto';

if (!globalThis.crypto) {
    globalThis.crypto = {
        getRandomValues: (buffer) => crypto.randomFillSync(buffer)
    };
}
