"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteUsersFromExternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var DeleteUsersFromExternalSquadCommand;
(function (DeleteUsersFromExternalSquadCommand) {
    DeleteUsersFromExternalSquadCommand.url = api_1.REST_API.EXTERNAL_SQUADS.BULK_ACTIONS.REMOVE_USERS;
    DeleteUsersFromExternalSquadCommand.TSQ_url = DeleteUsersFromExternalSquadCommand.url(':uuid');
    DeleteUsersFromExternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.BULK_ACTIONS.REMOVE_USERS(':uuid'), 'delete', 'Delete users from external squad');
    DeleteUsersFromExternalSquadCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteUsersFromExternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(DeleteUsersFromExternalSquadCommand || (exports.DeleteUsersFromExternalSquadCommand = DeleteUsersFromExternalSquadCommand = {}));
