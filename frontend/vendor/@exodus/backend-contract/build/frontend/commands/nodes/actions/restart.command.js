"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.RestartNodeCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var RestartNodeCommand;
(function (RestartNodeCommand) {
    RestartNodeCommand.url = api_1.REST_API.NODES.ACTIONS.RESTART;
    RestartNodeCommand.TSQ_url = RestartNodeCommand.url(':uuid');
    RestartNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.ACTIONS.RESTART(':uuid'), 'post', 'Restart node');
    RestartNodeCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    RestartNodeCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(RestartNodeCommand || (exports.RestartNodeCommand = RestartNodeCommand = {}));
