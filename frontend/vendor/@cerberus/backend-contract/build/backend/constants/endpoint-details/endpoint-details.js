"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getEndpointDetails = getEndpointDetails;
function getEndpointDetails(controllerUrl, requestMethod, methodDescription, methodLongDescription) {
    return {
        CONTROLLER_URL: controllerUrl,
        REQUEST_METHOD: requestMethod,
        METHOD_DESCRIPTION: methodDescription,
        METHOD_LONG_DESCRIPTION: methodLongDescription,
    };
}
