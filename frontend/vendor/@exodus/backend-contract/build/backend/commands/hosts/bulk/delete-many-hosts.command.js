"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkDeleteHostsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkDeleteHostsCommand;
(function (BulkDeleteHostsCommand) {
    BulkDeleteHostsCommand.url = api_1.REST_API.HOSTS.BULK.DELETE_HOSTS;
    BulkDeleteHostsCommand.TSQ_url = BulkDeleteHostsCommand.url;
    BulkDeleteHostsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.BULK.DELETE_HOSTS, 'post', 'Delete hosts by UUIDs', { scope: 'bulk-delete', kind: 'write' });
    BulkDeleteHostsCommand.RequestBodySchema = zod_1.z.object({
        uuids: zod_1.z.array(zod_1.z.uuid()),
    });

})(BulkDeleteHostsCommand || (exports.BulkDeleteHostsCommand = BulkDeleteHostsCommand = {}));
