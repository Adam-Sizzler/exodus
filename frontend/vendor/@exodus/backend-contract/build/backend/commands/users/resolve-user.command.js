"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResolveUserCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var ResolveUserCommand;
(function (ResolveUserCommand) {
    ResolveUserCommand.url = api_1.REST_API.USERS.RESOLVE;
    ResolveUserCommand.TSQ_url = ResolveUserCommand.url;
    ResolveUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.RESOLVE, 'post', 'Resolve a user', { scope: 'resolve', kind: 'read' }, 'Resolve a user by ID, Short UUID or username. Exactly one of the fields must be provided.');
    ResolveUserCommand.RequestBodySchema = zod_1.z
        .object({
        id: zod_1.z.number().optional(),
        shortUuid: zod_1.z.string().optional(),
        username: zod_1.z.string().optional(),
    })
        .refine((data) => {
        const provided = [data.id, data.shortUuid, data.username].filter((v) => v !== undefined);
        return provided.length === 1;
    }, {
        error: 'Exactly one of id, shortUuid, or username must be provided',
    });
    ResolveUserCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            id: zod_1.z.number(),
            username: zod_1.z.string(),
            shortUuid: zod_1.z.string(),
        }),
    });

})(ResolveUserCommand || (exports.ResolveUserCommand = ResolveUserCommand = {}));
