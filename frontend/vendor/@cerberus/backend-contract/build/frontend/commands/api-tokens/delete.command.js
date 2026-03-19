"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteApiTokenCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteApiTokenCommand;
(function (DeleteApiTokenCommand) {
    DeleteApiTokenCommand.url = api_1.REST_API.API_TOKENS.DELETE;
    DeleteApiTokenCommand.TSQ_url = DeleteApiTokenCommand.url(':uuid');
    DeleteApiTokenCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.API_TOKENS_ROUTES.DELETE(':uuid'), 'delete', 'Delete an API token by UUID', 'This endpoint is forbidden to use via "API-key". It can be used only with an admin JWT-token.');
    DeleteApiTokenCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteApiTokenCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.boolean(),
    });
})(DeleteApiTokenCommand || (exports.DeleteApiTokenCommand = DeleteApiTokenCommand = {}));
