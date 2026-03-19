"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ALPN_VALUES = exports.ALPN = void 0;
exports.ALPN = {
    H3: 'h3',
    H2: 'h2',
    HTTP_1_1: 'http/1.1',
    H_COMBINED: 'h2,http/1.1',
    H3_H2_H1_COMBINED: 'h3,h2,http/1.1',
    H3_H2_COMBINED: 'h3,h2',
};
exports.ALPN_VALUES = Object.values(exports.ALPN);
