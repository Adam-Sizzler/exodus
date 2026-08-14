"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeletePasskeyCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeletePasskeyCommand;
(function (DeletePasskeyCommand) {
    DeletePasskeyCommand.url = api_1.REST_API.PASSKEYS.DELETE_PASSKEY;
    DeletePasskeyCommand.TSQ_url = DeletePasskeyCommand.url;
    DeletePasskeyCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.PASSKEYS_ROUTES.DELETE_PASSKEY, 'delete', 'Delete a passkey by ID', { scope: 'delete', kind: 'write' });
    DeletePasskeyCommand.RequestBodySchema = zod_1.z.object({
        id: zod_1.z.string(),
    });
})(DeletePasskeyCommand || (exports.DeletePasskeyCommand = DeletePasskeyCommand = {}));
