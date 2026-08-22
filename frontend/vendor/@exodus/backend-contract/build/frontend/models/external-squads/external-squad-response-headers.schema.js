"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExternalSquadResponseHeadersRemoveSchema = exports.ExternalSquadResponseHeadersAddSchema = void 0;
const zod_1 = require("zod");
exports.ExternalSquadResponseHeadersAddSchema = zod_1.z.record(zod_1.z.string(), zod_1.z.string());
exports.ExternalSquadResponseHeadersRemoveSchema = zod_1.z.array(zod_1.z.string());
