"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.RestartAllNodesCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var RestartAllNodesCommand;
(function (RestartAllNodesCommand) {
    RestartAllNodesCommand.url = api_1.REST_API.NODES.ACTIONS.RESTART_ALL;
    RestartAllNodesCommand.TSQ_url = RestartAllNodesCommand.url;
    RestartAllNodesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.ACTIONS.RESTART_ALL, 'post', 'Restart all nodes');
    RestartAllNodesCommand.RequestBodySchema = zod_1.z.preprocess((val) => val || {}, zod_1.z.object({
        forceRestart: zod_1.z.boolean().optional(),
    }));
    RestartAllNodesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(RestartAllNodesCommand || (exports.RestartAllNodesCommand = RestartAllNodesCommand = {}));
