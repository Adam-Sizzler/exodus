"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenerateX25519Command = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var GenerateX25519Command;
(function (GenerateX25519Command) {
    GenerateX25519Command.url = api_1.REST_API.SYSTEM.TOOLS.GENERATE_X25519;
    GenerateX25519Command.TSQ_url = GenerateX25519Command.url;
    GenerateX25519Command.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.TOOLS.GENERATE_X25519, 'get', 'Generate 30 X25519 keypairs');
    GenerateX25519Command.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            keypairs: zod_1.z.array(zod_1.z.object({
                publicKey: zod_1.z.string(),
                privateKey: zod_1.z.string(),
            })),
        }),
    });
})(GenerateX25519Command || (exports.GenerateX25519Command = GenerateX25519Command = {}));
