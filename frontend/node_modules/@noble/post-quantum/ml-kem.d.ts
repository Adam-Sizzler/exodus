import { type KEM } from './utils.ts';
/** FIPS 203: 7. Parameter Sets */
export type KEMParam = {
    N: number;
    K: number;
    Q: number;
    ETA1: number;
    ETA2: number;
    du: number;
    dv: number;
    RBGstrength: number;
};
/** Internal params of ML-KEM versions */
export declare const PARAMS: Record<string, KEMParam>;
/** ML-KEM-512 for 128-bit security level. Not recommended after 2030, as per ASD. */
export declare const ml_kem512: KEM;
/** ML-KEM-768, for 192-bit security level. Not recommended after 2030, as per ASD. */
export declare const ml_kem768: KEM;
/** ML-KEM-1024 for 256-bit security level. OK after 2030, as per ASD. */
export declare const ml_kem1024: KEM;
//# sourceMappingURL=ml-kem.d.ts.map