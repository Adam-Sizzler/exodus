"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.numberParamSchema = void 0;
const zod_1 = require("zod");
exports.numberParamSchema = zod_1.z.coerce.number().positive();
