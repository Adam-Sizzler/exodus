"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AddUsersToInternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var AddUsersToInternalSquadCommand;
(function (AddUsersToInternalSquadCommand) {
    AddUsersToInternalSquadCommand.url = api_1.REST_API.INTERNAL_SQUADS.BULK_ACTIONS.ADD_USERS;
    AddUsersToInternalSquadCommand.TSQ_url = AddUsersToInternalSquadCommand.url(':uuid');
    AddUsersToInternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.BULK_ACTIONS.ADD_USERS(':uuid'), 'post', 'Add all users to internal squad');
    AddUsersToInternalSquadCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    AddUsersToInternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(AddUsersToInternalSquadCommand || (exports.AddUsersToInternalSquadCommand = AddUsersToInternalSquadCommand = {}));
