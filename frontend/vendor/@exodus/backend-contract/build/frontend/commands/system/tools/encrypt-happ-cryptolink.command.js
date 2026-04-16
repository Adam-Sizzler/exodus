"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EncryptHappCryptoLinkCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var EncryptHappCryptoLinkCommand;
(function (EncryptHappCryptoLinkCommand) {
    EncryptHappCryptoLinkCommand.url = api_1.REST_API.SYSTEM.TOOLS.ENCRYPT_HAPP_CRYPTO_LINK;
    EncryptHappCryptoLinkCommand.TSQ_url = EncryptHappCryptoLinkCommand.url;
    EncryptHappCryptoLinkCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.TOOLS.ENCRYPT_HAPP_CRYPTO_LINK, 'post', 'Encrypt Happ Crypto Link');
    EncryptHappCryptoLinkCommand.RequestSchema = zod_1.z.object({
        linkToEncrypt: zod_1.z.string().url(),
    });
    EncryptHappCryptoLinkCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            encryptedLink: zod_1.z.string(),
        }),
    });
})(EncryptHappCryptoLinkCommand || (exports.EncryptHappCryptoLinkCommand = EncryptHappCryptoLinkCommand = {}));
