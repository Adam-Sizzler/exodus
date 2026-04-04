/**
 * Post-Quantum Hybrid Cryptography
 *
 * The current implementation is flawed and likely redundant. We should offer
 * a small, generic API to compose hybrid schemes instead of reimplementing
 * protocol-specific logic (SSH, GPG, etc.) with ad hoc encodings.
 *
 * 1. Core Issues
 *    - sign/verify: implemented as two separate operations with different keys.
 *    - EC getSharedSecret: could be refactored into a proper KEM.
 *    - Multiple calls: keys, signatures, and shared secrets could be
 *      concatenated to reduce the number of API invocations.
 *    - Reinvention: most libraries add strange domain separations and
 *      encodings instead of simple byte concatenation.
 *
 * 2. API Goals
 *    - Provide primitives to build hybrids generically.
 *    - Avoid embedding SSH- or GPG-specific formats in the core API.
 *
 * 3. Edge Cases
 *    • Variable-length signatures:
 *      - DER-encoded (Weierstrass curves).
 *      - Falcon (unpadded).
 *      - Concatenation works only if length is fixed; otherwise a length
 *        prefix is required (but that breaks compatibility).
 *
 *    • getSharedSecret:
 *      - Default: non-KEM (authenticated ECDH).
 *      - KEM conversion: generate a random SK to remove implicit auth.
 *
 * 4. Common Pitfalls
 *    - Seed expansion:
 *      • Expanding a small seed into multiple keys reduces entropy.
 *      • API should allow identity mapping (no expansion).
 *
 *    - Skipping full point encoding:
 *      • Some omit the compression byte (parity) for WebCrypto compatibility.
 *      • Better: hash the raw secret; coordinate output is already non-uniform.
 *      • Some curves (e.g., X448) produce secrets that must be re-hashed to match
 *        symmetric-key lengths.
 *
 *    - Combiner inconsistencies:
 *      • Different domain separations and encodings across libraries.
 *      • Should live at the application layer, since key lengths vary.
 *
 * 5. Protocol Examples
 *    - SSH:
 *      • Concatenate keys.
 *      • Combiner: SHA-512.
 *
 *    - GPG:
 *      • Concatenate keys.
 *      • Combiner: SHA3-256(kemShare || ecdhShare || ciphertext || pubKey || algId || domSep || len(domSep))
 *
 *    - TLS:
 *      • Transcript-based derivation (HKDF).
 *
 * 6. Relevant Specs & Implementations
 *    - IETF Hybrid KEM drafts:
 *      • draft-irtf-cfrg-hybrid-kems
 *      • draft-connolly-cfrg-xwing-kem
 *      • draft-westerbaan-tls-xyber768d00
 *
 *    - PQC Libraries:
 *      • superdilithium (cyph/pqcrypto.js) – low adoption.
 *      • hybrid-pqc (DogeProtocol, quantumcoinproject) – complex encodings.
 *
 * 7. Signatures
 *    - Ed25519: fixed-size, easy to support.
 *    - Variable-size: introduces custom format requirements; best left to
 *      higher-level code.
 *
 * @module
 */
/*! noble-post-quantum - MIT License (c) 2024 Paul Miller (paulmillr.com) */
import { type EdDSA } from '@noble/curves/abstract/edwards.js';
import { type MontgomeryECDH } from '@noble/curves/abstract/montgomery.js';
import { type ECDSA } from '@noble/curves/abstract/weierstrass.js';
import { type CHash, type CHashXOF } from '@noble/hashes/utils.js';
import { type KEM, type Signer } from './utils.ts';
type CurveECDH = ECDSA | MontgomeryECDH;
type CurveSign = ECDSA | EdDSA;
export declare function ecdhKem(curve: CurveECDH, allowZeroKey?: boolean): KEM;
export declare function ecSigner(curve: CurveSign, allowZeroKey?: boolean): Signer;
export type ExpandSeed = (seed: Uint8Array, len: number) => Uint8Array;
type XOF = CHashXOF<any, {
    dkLen: number;
}>;
export declare function expandSeedXof(xof: XOF): ExpandSeed;
export type Combiner = (publicKeys: Uint8Array[], cipherTexts: Uint8Array[], sharedSecrets: Uint8Array[]) => Uint8Array;
export declare function combineKEMS(realSeedLen: number | undefined, // how much bytes expandSeed expects
realMsgLen: number | undefined, // how much bytes combiner returns
expandSeed: ExpandSeed, combiner: Combiner, ...kems: KEM[]): KEM;
export declare function combineSigners(realSeedLen: number | undefined, expandSeed: ExpandSeed, ...signers: Signer[]): Signer;
export declare function QSF(label: string, pqc: KEM, curveKEM: KEM, xof: XOF, kdf: CHash): KEM;
export declare const QSF_ml_kem768_p256: KEM;
export declare const QSF_ml_kem1024_p384: KEM;
export declare function createKitchenSink(label: string, pqc: KEM, curveKEM: KEM, xof: XOF, hash: CHash): KEM;
export declare const KitchenSink_ml_kem768_x25519: KEM;
export declare const ml_kem768_x25519: KEM;
export declare const ml_kem768_p256: KEM;
export declare const ml_kem1024_p384: KEM;
export declare const XWing: KEM;
export declare const MLKEM768X25519: KEM;
export declare const MLKEM768P256: KEM;
export declare const MLKEM1024P384: KEM;
export declare const QSFMLKEM768P256: KEM;
export declare const QSFMLKEM1024P384: KEM;
export declare const KitchenSinkMLKEM768X25519: KEM;
export {};
//# sourceMappingURL=hybrid.d.ts.map