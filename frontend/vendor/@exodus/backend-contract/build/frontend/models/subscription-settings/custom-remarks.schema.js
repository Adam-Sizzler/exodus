"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.CustomRemarksSchema = void 0;
const zod_1 = __importDefault(require("zod"));
exports.CustomRemarksSchema = zod_1.default.object({
    expiredUsers: zod_1.default.array(zod_1.default.string()).min(1),
    limitedUsers: zod_1.default.array(zod_1.default.string()).min(1),
    disabledUsers: zod_1.default.array(zod_1.default.string()).min(1),
    emptyHosts: zod_1.default.array(zod_1.default.string()).min(1),
    HWIDMaxDevicesExceeded: zod_1.default.array(zod_1.default.string()).min(1),
    HWIDNotSupported: zod_1.default.array(zod_1.default.string()).min(1),
});
