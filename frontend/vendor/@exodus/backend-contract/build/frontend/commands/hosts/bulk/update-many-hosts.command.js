"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateManyHostsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const update_command_1 = require("../update.command");
var UpdateManyHostsCommand;
(function (UpdateManyHostsCommand) {
    UpdateManyHostsCommand.url = api_1.REST_API.HOSTS.BULK.UPDATE;
    UpdateManyHostsCommand.TSQ_url = UpdateManyHostsCommand.url;
    UpdateManyHostsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.BULK.UPDATE, 'patch', 'Update many hosts', { scope: 'bulk-update', kind: 'write' });
    UpdateManyHostsCommand.RequestBodySchema = update_command_1.UpdateHostCommand.RequestBodySchema.omit({ uuid: true })
        .partial()
        .extend({
        uuids: zod_1.z.array(zod_1.z.uuid()).min(1),
    });
})(UpdateManyHostsCommand || (exports.UpdateManyHostsCommand = UpdateManyHostsCommand = {}));
