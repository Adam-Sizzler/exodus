export const INFRA_BILLING_CONTROLLER = 'infra-billing' as const;

export const INFRA_BILLING_ROUTES = {
    GET_PROVIDERS: 'providers', // Get list of all providers // get
    CREATE_PROVIDER: 'providers', // Create new provider // post
    UPDATE_PROVIDER: 'providers', // Update provider by uuid // patch
    DELETE_PROVIDER: (uuid: string) => `providers/${uuid}`, // Delete provider by uuid // delete
    GET_PROVIDER_BY_UUID: (uuid: string) => `providers/${uuid}`, // Get provider by uuid // get

    GET_BILLING_NODES: 'nodes', // Get list of all nodes billing // get
    CREATE_BILLING_NODE: 'nodes', // Create new node billing // post
    UPDATE_BILLING_NODE: 'nodes', // Update node billing by uuid // patch
    DELETE_BILLING_NODE: (uuid: string) => `nodes/${uuid}`, // Delete node billing by uuid // delete

    GET_BILLING_HISTORY: 'history', // Get list of all nodes billing history // get
    CREATE_BILLING_HISTORY: 'history', // Create new node billing history // post
    DELETE_BILLING_HISTORY: (uuid: string) => `history/${uuid}`, // Delete node billing history by uuid // delete
} as const;
