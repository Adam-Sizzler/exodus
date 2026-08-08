"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.UserResponseSchema = void 0;
const zod_1 = require("zod");
const extended_users_schema_1 = require("../../models/extended-users.schema");
exports.UserResponseSchema = (zod_1.z || zod_1.default.z || zod_1.default).object({
    response: extended_users_schema_1.ExtendedUsersSchema,
});
