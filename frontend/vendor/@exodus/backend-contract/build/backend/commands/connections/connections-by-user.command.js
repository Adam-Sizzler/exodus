"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ConnectionsByUserCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var ConnectionsByUserCommand;
(function (ConnectionsByUserCommand) {
    ConnectionsByUserCommand.url = api_1.REST_API.CONNECTIONS.CONNECTIONS_BY_USER;
    ConnectionsByUserCommand.TSQ_url = ConnectionsByUserCommand.url(':userId');
    ConnectionsByUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONNECTIONS_ROUTES.CONNECTIONS_BY_USER(':userId'), 'post', 'Request Connections for User', { scope: 'by-user', kind: 'read' });
    ConnectionsByUserCommand.RequestParamSchema = zod_1.z.object({
        userId: models_1.numberParamSchema,
    });
    ConnectionsByUserCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            jobId: zod_1.z.string(),
        }),
    });
})(ConnectionsByUserCommand || (exports.ConnectionsByUserCommand = ConnectionsByUserCommand = {}));
