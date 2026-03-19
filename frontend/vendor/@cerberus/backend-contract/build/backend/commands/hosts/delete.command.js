"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteHostCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var DeleteHostCommand;
(function (DeleteHostCommand) {
    DeleteHostCommand.url = api_1.REST_API.HOSTS.DELETE;
    DeleteHostCommand.TSQ_url = DeleteHostCommand.url(':uuid');
    DeleteHostCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.DELETE(':uuid'), 'delete', 'Delete a host by UUID');
    DeleteHostCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteHostCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isDeleted: zod_1.z.boolean(),
        }),
    });
})(DeleteHostCommand || (exports.DeleteHostCommand = DeleteHostCommand = {}));
