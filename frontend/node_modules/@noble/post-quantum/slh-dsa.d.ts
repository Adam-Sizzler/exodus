import { type CHash } from '@noble/hashes/utils.js';
import { type Signer } from './utils.ts';
/**
 * * N: Security parameter (in bytes). W: Winternitz parameter
 * * H: Hypertree height. D: Hypertree layers
 * * K: FORS trees numbers. A: FORS trees height
 */
export type SphincsOpts = {
    N: number;
    W: number;
    H: number;
    D: number;
    K: number;
    A: number;
    securityLevel: number;
};
export type SphincsHashOpts = {
    isCompressed?: boolean;
    getContext: GetContext;
};
/** Winternitz signature params. */
export declare const PARAMS: Record<string, SphincsOpts>;
/** Address, byte array of size ADDR_BYTES */
export type ADRS = Uint8Array;
export type Context = {
    PRFaddr: (addr: ADRS) => Uint8Array;
    PRFmsg: (skPRF: Uint8Array, random: Uint8Array, msg: Uint8Array) => Uint8Array;
    Hmsg: (R: Uint8Array, pk: Uint8Array, m: Uint8Array, outLen: number) => Uint8Array;
    thash1: (input: Uint8Array, addr: ADRS) => Uint8Array;
    thashN: (blocks: number, input: Uint8Array, addr: ADRS) => Uint8Array;
    clean: () => void;
};
export type GetContext = (opts: SphincsOpts) => (pub_seed: Uint8Array, sk_seed?: Uint8Array) => Context;
export type SphincsSigner = Signer & {
    internal: Signer;
    securityLevel: number;
    prehash: (hash: CHash) => Signer;
};
/** SLH-DSA: 128-bit fast SHAKE version. */
export declare const slh_dsa_shake_128f: SphincsSigner;
/** SLH-DSA: 128-bit short SHAKE version. */
export declare const slh_dsa_shake_128s: SphincsSigner;
/** SLH-DSA: 192-bit fast SHAKE version. */
export declare const slh_dsa_shake_192f: SphincsSigner;
/** SLH-DSA: 192-bit short SHAKE version. */
export declare const slh_dsa_shake_192s: SphincsSigner;
/** SLH-DSA: 256-bit fast SHAKE version. */
export declare const slh_dsa_shake_256f: SphincsSigner;
/** SLH-DSA: 256-bit short SHAKE version. */
export declare const slh_dsa_shake_256s: SphincsSigner;
/** SLH-DSA: 128-bit fast SHA2 version. */
export declare const slh_dsa_sha2_128f: SphincsSigner;
/** SLH-DSA: 128-bit small SHA2 version. */
export declare const slh_dsa_sha2_128s: SphincsSigner;
/** SLH-DSA: 192-bit fast SHA2 version. */
export declare const slh_dsa_sha2_192f: SphincsSigner;
/** SLH-DSA: 192-bit small SHA2 version. */
export declare const slh_dsa_sha2_192s: SphincsSigner;
/** SLH-DSA: 256-bit fast SHA2 version. */
export declare const slh_dsa_sha2_256f: SphincsSigner;
/** SLH-DSA: 256-bit small SHA2 version. */
export declare const slh_dsa_sha2_256s: SphincsSigner;
//# sourceMappingURL=slh-dsa.d.ts.map