"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResetNodeTrafficCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var ResetNodeTrafficCommand;
(function (ResetNodeTrafficCommand) {
    ResetNodeTrafficCommand.url = api_1.REST_API.NODES.ACTIONS.RESET_TRAFFIC;
    ResetNodeTrafficCommand.TSQ_url = ResetNodeTrafficCommand.url(':uuid');
    ResetNodeTrafficCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.ACTIONS.RESET_TRAFFIC(':uuid'), 'post', 'Reset Node Traffic');
    ResetNodeTrafficCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    ResetNodeTrafficCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(ResetNodeTrafficCommand || (exports.ResetNodeTrafficCommand = ResetNodeTrafficCommand = {}));
