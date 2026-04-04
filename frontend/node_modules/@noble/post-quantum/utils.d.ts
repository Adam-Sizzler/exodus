/**
 * Utilities for hex, bytearray and number handling.
 * @module
 */
/*! noble-post-quantum - MIT License (c) 2024 Paul Miller (paulmillr.com) */
import { type CHash, type TypedArray, concatBytes, randomBytes as randb } from '@noble/hashes/utils.js';
export { abytes } from '@noble/hashes/utils.js';
export { concatBytes };
export declare const randomBytes: typeof randb;
export declare function equalBytes(a: Uint8Array, b: Uint8Array): boolean;
export declare function copyBytes(bytes: Uint8Array): Uint8Array;
export type CryptoKeys = {
    info?: {
        type?: string;
    };
    lengths: {
        seed?: number;
        publicKey?: number;
        secretKey?: number;
    };
    keygen: (seed?: Uint8Array) => {
        secretKey: Uint8Array;
        publicKey: Uint8Array;
    };
    getPublicKey: (secretKey: Uint8Array) => Uint8Array;
};
export type VerOpts = {
    context?: Uint8Array;
};
export type SigOpts = VerOpts & {
    extraEntropy?: Uint8Array | false;
};
export declare function validateOpts(opts: object): void;
export declare function validateVerOpts(opts: VerOpts): void;
export declare function validateSigOpts(opts: SigOpts): void;
/** Generic interface for signatures. Has keygen, sign and verify. */
export type Signer = CryptoKeys & {
    lengths: {
        signRand?: number;
        signature?: number;
    };
    sign: (msg: Uint8Array, secretKey: Uint8Array, opts?: SigOpts) => Uint8Array;
    verify: (sig: Uint8Array, msg: Uint8Array, publicKey: Uint8Array, opts?: VerOpts) => boolean;
};
export type KEM = CryptoKeys & {
    lengths: {
        cipherText?: number;
        msg?: number;
        msgRand?: number;
    };
    encapsulate: (publicKey: Uint8Array, msg?: Uint8Array) => {
        cipherText: Uint8Array;
        sharedSecret: Uint8Array;
    };
    decapsulate: (cipherText: Uint8Array, secretKey: Uint8Array) => Uint8Array;
};
export interface Coder<F, T> {
    encode(from: F): T;
    decode(to: T): F;
}
export interface BytesCoder<T> extends Coder<T, Uint8Array> {
    encode: (data: T) => Uint8Array;
    decode: (bytes: Uint8Array) => T;
}
export type BytesCoderLen<T> = BytesCoder<T> & {
    bytesLen: number;
};
type UnCoder<T> = T extends BytesCoder<infer U> ? U : never;
type SplitOut<T extends (number | BytesCoderLen<any>)[]> = {
    [K in keyof T]: T[K] extends number ? Uint8Array : UnCoder<T[K]>;
};
export declare function splitCoder<T extends (number | BytesCoderLen<any>)[]>(label: string, ...lengths: T): BytesCoder<SplitOut<T>> & {
    bytesLen: number;
};
export declare function vecCoder<T>(c: BytesCoderLen<T>, vecLen: number): BytesCoderLen<T[]>;
export declare function cleanBytes(...list: (TypedArray | TypedArray[])[]): void;
export declare function getMask(bits: number): number;
export declare const EMPTY: Uint8Array;
export declare function getMessage(msg: Uint8Array, ctx?: Uint8Array): Uint8Array;
export declare function checkHash(hash: CHash, requiredStrength?: number): void;
export declare function getMessagePrehash(hash: CHash, msg: Uint8Array, ctx?: Uint8Array): Uint8Array;
//# sourceMappingURL=utils.d.ts.map