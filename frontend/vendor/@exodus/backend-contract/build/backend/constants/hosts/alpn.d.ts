export declare const ALPN: {
    readonly H3: "h3";
    readonly H2: "h2";
    readonly HTTP_1_1: "http/1.1";
    readonly H_COMBINED: "h2,http/1.1";
    readonly H3_H2_H1_COMBINED: "h3,h2,http/1.1";
    readonly H3_H2_COMBINED: "h3,h2";
};
export type TAlpn = [keyof typeof ALPN][number];
export declare const ALPN_VALUES: ("h3" | "h2" | "http/1.1" | "h2,http/1.1" | "h3,h2,http/1.1" | "h3,h2")[];
//# sourceMappingURL=alpn.d.ts.map