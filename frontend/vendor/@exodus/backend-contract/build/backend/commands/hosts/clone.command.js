"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CloneHostCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const host_response_1 = require("./host.response");
var CloneHostCommand;
(function (CloneHostCommand) {
    CloneHostCommand.url = api_1.REST_API.HOSTS.ACTIONS.CLONE;
    CloneHostCommand.TSQ_url = CloneHostCommand.url;
    CloneHostCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.ACTIONS.CLONE, 'post', 'Clone host', { scope: 'clone', kind: 'write' });
    CloneHostCommand.RequestBodySchema = zod_1.z.object({
        cloneFromUuid: zod_1.z.uuid(),
    });
    CloneHostCommand.ResponseSchema = host_response_1.HostResponseSchema;
})(CloneHostCommand || (exports.CloneHostCommand = CloneHostCommand = {}));
