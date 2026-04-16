"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExternalSquadResponseHeadersSchema = void 0;
const zod_1 = require("zod");
exports.ExternalSquadResponseHeadersSchema = zod_1.z.nullable(zod_1.z.record(zod_1.z.string(), zod_1.z.string()));
