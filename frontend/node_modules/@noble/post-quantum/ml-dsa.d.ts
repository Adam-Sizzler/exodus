import { type CryptoKeys, type Signer, type SigOpts, type VerOpts } from './utils.ts';
export type DSAInternalOpts = {
    externalMu?: boolean;
};
/** Signer API, containing internal methods */
export type DSAInternal = CryptoKeys & {
    lengths: Signer['lengths'];
    sign: (msg: Uint8Array, secretKey: Uint8Array, opts?: SigOpts & DSAInternalOpts) => Uint8Array;
    verify: (sig: Uint8Array, msg: Uint8Array, pubKey: Uint8Array, opts?: VerOpts & DSAInternalOpts) => boolean;
};
export type DSA = Signer & {
    internal: DSAInternal;
};
/** Various lattice params. */
export type DSAParam = {
    K: number;
    L: number;
    D: number;
    GAMMA1: number;
    GAMMA2: number;
    TAU: number;
    ETA: number;
    OMEGA: number;
};
/** Internal params for different versions of ML-DSA  */
export declare const PARAMS: Record<string, DSAParam>;
/** ML-DSA-44 for 128-bit security level. Not recommended after 2030, as per ASD. */
export declare const ml_dsa44: DSA;
/** ML-DSA-65 for 192-bit security level. Not recommended after 2030, as per ASD. */
export declare const ml_dsa65: DSA;
/** ML-DSA-87 for 256-bit security level. OK after 2030, as per ASD. */
export declare const ml_dsa87: DSA;
//# sourceMappingURL=ml-dsa.d.ts.map