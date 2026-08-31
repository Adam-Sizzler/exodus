export declare const ERRORS: {
    readonly INTERNAL_SERVER_ERROR: {
        readonly code: "A001";
        readonly message: "Server error";
        readonly httpCode: 500;
    };
    readonly LOGIN_ERROR: {
        readonly code: "A002";
        readonly message: "Login error";
        readonly httpCode: 500;
    };
    readonly UNAUTHORIZED: {
        readonly code: "A003";
        readonly message: "Unauthorized";
        readonly httpCode: 401;
    };
    readonly FORBIDDEN_ROLE_ERROR: {
        readonly code: "A004";
        readonly message: "Forbidden role error";
        readonly httpCode: 403;
    };
    readonly CREATE_API_TOKEN_ERROR: {
        readonly code: "A005";
        readonly message: "Create API token error";
        readonly httpCode: 500;
    };
    readonly DELETE_API_TOKEN_ERROR: {
        readonly code: "A006";
        readonly message: "Delete API token error";
        readonly httpCode: 500;
    };
    readonly REQUESTED_TOKEN_NOT_FOUND: {
        readonly code: "A007";
        readonly message: "Requested token not found";
        readonly httpCode: 404;
    };
    readonly FIND_ALL_API_TOKENS_ERROR: {
        readonly code: "A008";
        readonly message: "Find all API tokens error";
        readonly httpCode: 500;
    };
    readonly GET_PUBLIC_KEY_ERROR: {
        readonly code: "A009";
        readonly message: "Get public key error";
        readonly httpCode: 500;
    };
    readonly ENABLE_NODE_ERROR: {
        readonly code: "A010";
        readonly message: "Enable node error";
        readonly httpCode: 500;
    };
    readonly NODE_NOT_FOUND: {
        readonly code: "A011";
        readonly message: "Node not found";
        readonly httpCode: 404;
    };
    readonly CONFIG_NOT_FOUND: {
        readonly code: "A012";
        readonly message: "Configuration not found";
        readonly httpCode: 404;
    };
    readonly UPDATE_CONFIG_ERROR: {
        readonly code: "A013";
        readonly message: "Error updating configuration";
        readonly httpCode: 500;
    };
    readonly GET_CONFIG_ERROR: {
        readonly code: "A014";
        readonly message: "Error retrieving configuration";
        readonly httpCode: 500;
    };
    readonly DELETE_MANY_INBOUNDS_ERROR: {
        readonly code: "A015";
        readonly message: "Delete many inbounds error";
        readonly httpCode: 500;
    };
    readonly CREATE_MANY_INBOUNDS_ERROR: {
        readonly code: "A016";
        readonly message: "Create many inbounds error";
        readonly httpCode: 500;
    };
    readonly FIND_ALL_INBOUNDS_ERROR: {
        readonly code: "A017";
        readonly message: "Find all inbounds error";
        readonly httpCode: 500;
    };
    readonly CREATE_USER_ERROR: {
        readonly code: "A018";
        readonly message: "Failed to create user";
        readonly httpCode: 500;
    };
    readonly USER_USERNAME_ALREADY_EXISTS: {
        readonly code: "A019";
        readonly message: "User username already exists";
        readonly httpCode: 400;
    };
    readonly USER_SHORT_UUID_ALREADY_EXISTS: {
        readonly code: "A020";
        readonly message: "User short UUID already exists";
        readonly httpCode: 400;
    };
    readonly USER_SUBSCRIPTION_UUID_ALREADY_EXISTS: {
        readonly code: "A021";
        readonly message: "User subscription UUID already exists";
        readonly httpCode: 400;
    };
    readonly CREATE_USER_WITH_INBOUNDS_ERROR: {
        readonly code: "A022";
        readonly message: "User creation successful, but inbound creation failed. User not created.";
        readonly httpCode: 500;
    };
    readonly CANT_GET_CREATED_USER_WITH_INBOUNDS: {
        readonly code: "A023";
        readonly message: "User creation successful, but failed to get created user with inbounds.";
        readonly httpCode: 500;
    };
    readonly GET_ALL_USERS_ERROR: {
        readonly code: "A024";
        readonly message: "Get all users error";
        readonly httpCode: 500;
    };
    readonly USER_NOT_FOUND: {
        readonly code: "A025";
        readonly message: "User not found";
        readonly httpCode: 404;
    };
    readonly GET_USER_BY_ERROR: {
        readonly code: "A026";
        readonly message: "Get user by error";
        readonly httpCode: 500;
    };
    readonly REVOKE_USER_SUBSCRIPTION_ERROR: {
        readonly code: "A027";
        readonly message: "Revoke user subscription error";
        readonly httpCode: 500;
    };
    readonly DISABLE_USER_ERROR: {
        readonly code: "A028";
        readonly message: "Disable user error";
        readonly httpCode: 500;
    };
    readonly USER_ALREADY_DISABLED: {
        readonly code: "A029";
        readonly message: "User already disabled";
        readonly httpCode: 400;
    };
    readonly USER_ALREADY_ENABLED: {
        readonly code: "A030";
        readonly message: "User already enabled";
        readonly httpCode: 400;
    };
    readonly ENABLE_USER_ERROR: {
        readonly code: "A031";
        readonly message: "Enable user error";
        readonly httpCode: 500;
    };
    readonly CREATE_NODE_ERROR: {
        readonly code: "A032";
        readonly message: "Create node error";
        readonly httpCode: 500;
    };
    readonly NODE_NAME_ALREADY_EXISTS: {
        readonly code: "A033";
        readonly message: "Node name already exists";
        readonly httpCode: 400;
    };
    readonly NODE_ADDRESS_ALREADY_EXISTS: {
        readonly code: "A034";
        readonly message: "Node address already exists";
        readonly httpCode: 400;
    };
    readonly NODE_ERROR_WITH_MSG: {
        readonly code: "N001";
        readonly message: "";
        readonly httpCode: 500;
        readonly withMessage: (message: string) => {
            code: string;
            message: string;
            httpCode: number;
        };
    };
    readonly NODE_ERROR_500_WITH_MSG: {
        readonly code: "N002";
        readonly message: "";
        readonly httpCode: 500;
        readonly withMessage: (message: string) => {
            code: string;
            message: string;
            httpCode: number;
        };
    };
    readonly RESTART_NODE_ERROR: {
        readonly code: "A035";
        readonly message: "Restart node error";
        readonly httpCode: 500;
    };
    readonly GET_CONFIG_WITH_USERS_ERROR: {
        readonly code: "A036";
        readonly message: "Get config with users error";
        readonly httpCode: 500;
    };
    readonly DELETE_USER_ERROR: {
        readonly code: "A037";
        readonly message: "Delete user error";
        readonly httpCode: 500;
    };
    readonly UPDATE_NODE_ERROR: {
        readonly code: "A038";
        readonly message: "Update node error";
        readonly httpCode: 500;
    };
    readonly UPDATE_USER_ERROR: {
        readonly code: "A039";
        readonly message: "Update user error";
        readonly httpCode: 500;
    };
    readonly INCREMENT_USED_TRAFFIC_ERROR: {
        readonly code: "A040";
        readonly message: "Increment used traffic error";
        readonly httpCode: 500;
    };
    readonly GET_ALL_NODES_ERROR: {
        readonly code: "A041";
        readonly message: "Get all nodes error";
        readonly httpCode: 500;
    };
    readonly GET_ONE_NODE_ERROR: {
        readonly code: "A042";
        readonly message: "Get one node error";
        readonly httpCode: 500;
    };
    readonly DELETE_NODE_ERROR: {
        readonly code: "A043";
        readonly message: "Delete node error";
        readonly httpCode: 500;
    };
    readonly CREATE_HOST_ERROR: {
        readonly code: "A044";
        readonly message: "Create host error";
        readonly httpCode: 500;
    };
    readonly HOST_REMARK_ALREADY_EXISTS: {
        readonly code: "A045";
        readonly message: "Host remark already exists";
        readonly httpCode: 400;
    };
    readonly HOST_NOT_FOUND: {
        readonly code: "A046";
        readonly message: "Host not found";
        readonly httpCode: 404;
    };
    readonly DELETE_HOST_ERROR: {
        readonly code: "A047";
        readonly message: "Delete host error";
        readonly httpCode: 500;
    };
    readonly GET_USER_STATS_ERROR: {
        readonly code: "A048";
        readonly message: "Get user stats error";
        readonly httpCode: 500;
    };
    readonly UPDATE_USER_WITH_INBOUNDS_ERROR: {
        readonly code: "A049";
        readonly message: "Update user with inbounds error";
        readonly httpCode: 500;
    };
    readonly GET_ALL_HOSTS_ERROR: {
        readonly code: "A050";
        readonly message: "Get all hosts error";
        readonly httpCode: 500;
    };
    readonly REORDER_HOSTS_ERROR: {
        readonly code: "A051";
        readonly message: "Reorder hosts error";
        readonly httpCode: 500;
    };
    readonly UPDATE_HOST_ERROR: {
        readonly code: "A052";
        readonly message: "Update host error";
        readonly httpCode: 500;
    };
    readonly CREATE_CONFIG_ERROR: {
        readonly code: "A053";
        readonly message: "Create config error";
        readonly httpCode: 500;
    };
    readonly ENABLED_NODES_NOT_FOUND: {
        readonly code: "A054";
        readonly message: "Enabled nodes not found";
        readonly httpCode: 409;
    };
    readonly GET_NODES_USAGE_BY_RANGE_ERROR: {
        readonly code: "A055";
        readonly message: "Get nodes usage by range error";
        readonly httpCode: 500;
    };
    readonly RESET_USER_TRAFFIC_ERROR: {
        readonly code: "A056";
        readonly message: "Reset user traffic error";
        readonly httpCode: 500;
    };
    readonly REORDER_NODES_ERROR: {
        readonly code: "A057";
        readonly message: "Reorder nodes error";
        readonly httpCode: 500;
    };
    readonly GET_ALL_INBOUNDS_ERROR: {
        readonly code: "A058";
        readonly message: "Get all inbounds error";
        readonly httpCode: 500;
    };
    readonly BULK_DELETE_USERS_BY_STATUS_ERROR: {
        readonly code: "A059";
        readonly message: "Bulk delete users by status error";
        readonly httpCode: 500;
    };
    readonly UPDATE_INBOUND_ERROR: {
        readonly code: "A060";
        readonly message: "Update inbound error";
        readonly httpCode: 500;
    };
    readonly CONFIG_VALIDATION_ERROR: {
        readonly code: "A061";
        readonly message: "Config validation error";
        readonly httpCode: 500;
        readonly withMessage: (message: string) => {
            code: string;
            message: string;
            httpCode: number;
        };
    };
    readonly USERS_NOT_FOUND: {
        readonly code: "A062";
        readonly message: "Users not found";
        readonly httpCode: 404;
    };
    readonly GET_USER_BY_UNIQUE_FIELDS_NOT_FOUND: {
        readonly code: "A063";
        readonly message: "User with specified params not found";
        readonly httpCode: 404;
    };
    readonly UPDATE_EXCEEDED_TRAFFIC_USERS_ERROR: {
        readonly code: "A064";
        readonly message: "Update exceeded traffic users error";
        readonly httpCode: 500;
    };
    readonly ADMIN_NOT_FOUND: {
        readonly code: "A065";
        readonly message: "Admin not found";
        readonly httpCode: 404;
    };
    readonly CREATE_ADMIN_ERROR: {
        readonly code: "A066";
        readonly message: "Create admin error";
        readonly httpCode: 500;
    };
    readonly GET_AUTH_STATUS_ERROR: {
        readonly code: "A067";
        readonly message: "Get auth status error";
        readonly httpCode: 500;
    };
    readonly FORBIDDEN: {
        readonly code: "A068";
        readonly message: "Forbidden";
        readonly httpCode: 403;
    };
    readonly DISABLE_NODE_ERROR: {
        readonly code: "A069";
        readonly message: "Disable node error";
        readonly httpCode: 500;
    };
    readonly GET_ONE_HOST_ERROR: {
        readonly code: "A070";
        readonly message: "Get one host error";
        readonly httpCode: 500;
    };
    readonly SUBSCRIPTION_SETTINGS_NOT_FOUND: {
        readonly code: "A071";
        readonly message: "Subscription settings not found";
        readonly httpCode: 404;
    };
    readonly GET_SUBSCRIPTION_SETTINGS_ERROR: {
        readonly code: "A072";
        readonly message: "Get subscription settings error";
        readonly httpCode: 500;
    };
    readonly UPDATE_SUBSCRIPTION_SETTINGS_ERROR: {
        readonly code: "A073";
        readonly message: "Update subscription settings error";
        readonly httpCode: 500;
    };
    readonly ADD_INBOUND_TO_USERS_ERROR: {
        readonly code: "A074";
        readonly message: "Add inbound to users error";
        readonly httpCode: 500;
    };
    readonly REMOVE_INBOUND_FROM_USERS_ERROR: {
        readonly code: "A075";
        readonly message: "Remove inbound from users error";
        readonly httpCode: 500;
    };
    readonly INBOUND_NOT_FOUND: {
        readonly code: "A076";
        readonly message: "Inbound not found";
        readonly httpCode: 404;
    };
    readonly ADD_INBOUND_TO_NODES_ERROR: {
        readonly code: "A077";
        readonly message: "Add inbound to nodes error";
        readonly httpCode: 500;
    };
    readonly REMOVE_INBOUND_FROM_NODES_ERROR: {
        readonly code: "A078";
        readonly message: "Remove inbound from nodes error";
        readonly httpCode: 500;
    };
    readonly DELETE_HOSTS_ERROR: {
        readonly code: "A079";
        readonly message: "Delete hosts error";
        readonly httpCode: 500;
    };
    readonly BULK_ENABLE_HOSTS_ERROR: {
        readonly code: "A080";
        readonly message: "Bulk enable hosts error";
        readonly httpCode: 500;
    };
    readonly BULK_DISABLE_HOSTS_ERROR: {
        readonly code: "A081";
        readonly message: "Bulk disable hosts error";
        readonly httpCode: 500;
    };
    readonly BULK_DELETE_USERS_BY_USER_IDS_ERROR: {
        readonly code: "A084";
        readonly message: "Bulk delete users by user IDs error";
        readonly httpCode: 500;
    };
    readonly BULK_REVOKE_USERS_SUBSCRIPTION_ERROR: {
        readonly code: "A085";
        readonly message: "Bulk revoke users subscription error";
        readonly httpCode: 500;
    };
    readonly BULK_RESET_USER_TRAFFIC_ERROR: {
        readonly code: "A086";
        readonly message: "Bulk reset user traffic error";
        readonly httpCode: 500;
    };
    readonly BULK_UPDATE_USERS_ERROR: {
        readonly code: "A087";
        readonly message: "Bulk update users error";
        readonly httpCode: 500;
    };
    readonly BULK_ADD_INBOUNDS_TO_USERS_ERROR: {
        readonly code: "A088";
        readonly message: "Bulk add inbounds to users error";
        readonly httpCode: 500;
    };
    readonly BULK_UPDATE_ALL_USERS_ERROR: {
        readonly code: "A089";
        readonly message: "Bulk update all users error";
        readonly httpCode: 500;
    };
    readonly INVALID_USER_STATUS_ERROR: {
        readonly code: "A089";
        readonly message: "LIMITED and EXPIRED statuses are not allowed to be set manually.";
        readonly httpCode: 400;
    };
    readonly KEYPAIR_CREATION_ERROR: {
        readonly code: "A090";
        readonly message: "Keypair creation error";
        readonly httpCode: 500;
    };
    readonly GET_USER_USAGE_BY_RANGE_ERROR: {
        readonly code: "A091";
        readonly message: "Get user usage by range error";
        readonly httpCode: 500;
    };
    readonly KEYPAIR_NOT_FOUND: {
        readonly code: "A092";
        readonly message: "Keypair not found. Restart app.";
        readonly httpCode: 500;
    };
    readonly ACTIVATE_ALL_INBOUNDS_ERROR: {
        readonly code: "A093";
        readonly message: "Activate all inbounds error";
        readonly httpCode: 500;
    };
    readonly CREATE_HWID_USER_DEVICE_ERROR: {
        readonly code: "A096";
        readonly message: "Create hwid user device error";
        readonly httpCode: 500;
    };
    readonly CHECK_HWID_EXISTS_ERROR: {
        readonly code: "A097";
        readonly message: "Check hwid exists error";
        readonly httpCode: 500;
    };
    readonly USER_HWID_DEVICE_ALREADY_EXISTS: {
        readonly code: "A098";
        readonly message: "User hwid device already exists";
        readonly httpCode: 400;
    };
    readonly USER_HWID_DEVICE_LIMIT_REACHED: {
        readonly code: "A099";
        readonly message: "User hwid device limit reached";
        readonly httpCode: 400;
    };
    readonly GET_USER_HWID_DEVICES_ERROR: {
        readonly code: "A100";
        readonly message: "Get user hwid devices error";
        readonly httpCode: 500;
    };
    readonly DELETE_HWID_USER_DEVICE_ERROR: {
        readonly code: "A101";
        readonly message: "Delete hwid user device error";
        readonly httpCode: 500;
    };
    readonly UPSERT_HWID_USER_DEVICE_ERROR: {
        readonly code: "A102";
        readonly message: "Upsert hwid user device error";
        readonly httpCode: 500;
    };
    readonly GET_ALL_TAGS_ERROR: {
        readonly code: "A103";
        readonly message: "Get all tags error";
        readonly httpCode: 500;
    };
    readonly GETTING_ALL_SUBSCRIPTIONS_ERROR: {
        readonly code: "A104";
        readonly message: "Getting all subscriptions error";
        readonly httpCode: 500;
    };
    readonly TRIGGER_THRESHOLD_NOTIFICATION_ERROR: {
        readonly code: "A105";
        readonly message: "Trigger threshold notification error";
        readonly httpCode: 500;
    };
    readonly BULK_DELETE_BY_STATUS_ERROR: {
        readonly code: "A106";
        readonly message: "Bulk delete by status error";
        readonly httpCode: 500;
    };
    readonly CLEAN_OLD_USAGE_RECORDS_ERROR: {
        readonly code: "A107";
        readonly message: "Clean old usage records error";
        readonly httpCode: 500;
    };
    readonly VACUUM_TABLE_ERROR: {
        readonly code: "A108";
        readonly message: "Vacuum table error";
        readonly httpCode: 500;
    };
    readonly GET_CONFIG_PROFILES_ERROR: {
        readonly code: "A109";
        readonly message: "Get config profiles error";
        readonly httpCode: 500;
    };
    readonly GET_CONFIG_PROFILE_BY_UUID_ERROR: {
        readonly code: "A110";
        readonly message: "Get config profile by UUID error";
        readonly httpCode: 500;
    };
    readonly CONFIG_PROFILE_NOT_FOUND: {
        readonly code: "A111";
        readonly message: "Config profile not found";
        readonly httpCode: 404;
    };
    readonly CREATE_CONFIG_PROFILE_ERROR: {
        readonly code: "A112";
        readonly message: "Create config profile error";
        readonly httpCode: 500;
    };
    readonly INBOUNDS_WITH_SAME_TAG_ALREADY_EXISTS: {
        readonly code: "A113";
        readonly message: "Inbounds with same tag already exists in database. Inbound tags must be unique.";
        readonly httpCode: 409;
    };
    readonly CONFIG_PROFILE_NAME_ALREADY_EXISTS: {
        readonly code: "A114";
        readonly message: "Config profile name already exists in database. Config profile names must be unique.";
        readonly httpCode: 409;
    };
    readonly GET_INBOUNDS_BY_PROFILE_UUID_ERROR: {
        readonly code: "A115";
        readonly message: "Get inbounds by profile UUID error";
        readonly httpCode: 500;
    };
    readonly GET_INTERNAL_SQUADS_ERROR: {
        readonly code: "A116";
        readonly message: "Get internal squads error";
        readonly httpCode: 500;
    };
    readonly GET_INTERNAL_SQUAD_BY_UUID_ERROR: {
        readonly code: "A117";
        readonly message: "Get internal squad by UUID error";
        readonly httpCode: 500;
    };
    readonly INTERNAL_SQUAD_NOT_FOUND: {
        readonly code: "A118";
        readonly message: "Internal squad not found";
        readonly httpCode: 404;
    };
    readonly CREATE_INTERNAL_SQUAD_ERROR: {
        readonly code: "A119";
        readonly message: "Create internal squad error";
        readonly httpCode: 500;
    };
    readonly INTERNAL_SQUAD_NAME_ALREADY_EXISTS: {
        readonly code: "A120";
        readonly message: "Internal squad name already exists";
        readonly httpCode: 409;
    };
    readonly UPDATE_INTERNAL_SQUAD_ERROR: {
        readonly code: "A121";
        readonly message: "Update internal squad error";
        readonly httpCode: 500;
    };
    readonly DELETE_INTERNAL_SQUAD_ERROR: {
        readonly code: "A122";
        readonly message: "Delete internal squad error";
        readonly httpCode: 500;
    };
    readonly CREATE_USER_WITH_INTERNAL_SQUAD_ERROR: {
        readonly code: "A123";
        readonly message: "Create user with internal squad error";
        readonly httpCode: 500;
    };
    readonly CONFIG_PROFILE_INBOUND_NOT_FOUND_IN_SPECIFIED_PROFILE: {
        readonly code: "A124";
        readonly message: "Config profile inbound not found in specified profile";
        readonly httpCode: 404;
    };
    readonly GET_USER_ACCESSIBLE_NODES_ERROR: {
        readonly code: "A125";
        readonly message: "Get user accessible nodes error";
        readonly httpCode: 500;
    };
    readonly GET_INFRA_PROVIDERS_ERROR: {
        readonly code: "A126";
        readonly message: "Get infra providers error";
        readonly httpCode: 500;
    };
    readonly GET_INFRA_PROVIDER_BY_UUID_ERROR: {
        readonly code: "A127";
        readonly message: "Get infra provider by UUID error";
        readonly httpCode: 500;
    };
    readonly INFRA_PROVIDER_NOT_FOUND: {
        readonly code: "A128";
        readonly message: "Infra provider not found";
        readonly httpCode: 404;
    };
    readonly DELETE_INFRA_PROVIDER_BY_UUID_ERROR: {
        readonly code: "A129";
        readonly message: "Delete infra provider by UUID error";
        readonly httpCode: 500;
    };
    readonly CREATE_INFRA_PROVIDER_ERROR: {
        readonly code: "A130";
        readonly message: "Create infra provider error";
        readonly httpCode: 500;
    };
    readonly UPDATE_INFRA_PROVIDER_ERROR: {
        readonly code: "A131";
        readonly message: "Update infra provider error";
        readonly httpCode: 500;
    };
    readonly CREATE_INFRA_BILLING_HISTORY_RECORD_ERROR: {
        readonly code: "A132";
        readonly message: "Create infra billing history record error";
        readonly httpCode: 500;
    };
    readonly GET_INFRA_BILLING_HISTORY_RECORDS_ERROR: {
        readonly code: "A133";
        readonly message: "Get infra billing history records error";
        readonly httpCode: 500;
    };
    readonly DELETE_INFRA_BILLING_HISTORY_RECORD_BY_UUID_ERROR: {
        readonly code: "A134";
        readonly message: "Delete infra billing history record by UUID error";
        readonly httpCode: 500;
    };
    readonly GET_BILLING_NODES_ERROR: {
        readonly code: "A135";
        readonly message: "Get billing nodes error";
        readonly httpCode: 500;
    };
    readonly UPDATE_INFRA_BILLING_NODE_ERROR: {
        readonly code: "A136";
        readonly message: "Update infra billing node error";
        readonly httpCode: 500;
    };
    readonly CREATE_INFRA_BILLING_NODE_ERROR: {
        readonly code: "A137";
        readonly message: "Create infra billing node error";
        readonly httpCode: 500;
    };
    readonly DELETE_INFRA_BILLING_NODE_BY_UUID_ERROR: {
        readonly code: "A138";
        readonly message: "Delete infra billing node by UUID error";
        readonly httpCode: 500;
    };
    readonly GET_BILLING_NODES_FOR_NOTIFICATIONS_ERROR: {
        readonly code: "A139";
        readonly message: "Get billing nodes for notifications error";
        readonly httpCode: 500;
    };
    readonly ADD_USERS_TO_INTERNAL_SQUAD_ERROR: {
        readonly code: "A140";
        readonly message: "Add users to internal squad error";
        readonly httpCode: 500;
    };
    readonly INTERNAL_SQUAD_BULK_ACTIONS_ERROR: {
        readonly code: "A141";
        readonly message: "Internal squad bulk actions error";
        readonly httpCode: 500;
    };
    readonly REMOVE_USERS_FROM_INTERNAL_SQUAD_ERROR: {
        readonly code: "A142";
        readonly message: "Remove users from internal squad error";
        readonly httpCode: 500;
    };
    readonly DELETE_CONFIG_PROFILE_BY_UUID_ERROR: {
        readonly code: "A143";
        readonly message: "Delete config profile by UUID error";
        readonly httpCode: 500;
    };
    readonly RESERVED_INTERNAL_SQUAD_NAME: {
        readonly code: "A144";
        readonly message: "This name is reserved by Exodus. Please use a different name.";
        readonly httpCode: 400;
    };
    readonly RESERVED_CONFIG_PROFILE_NAME: {
        readonly code: "A145";
        readonly message: "This name is reserved by Exodus. Please use a different name.";
        readonly httpCode: 400;
    };
    readonly UPDATE_CONFIG_PROFILE_ERROR: {
        readonly code: "A146";
        readonly message: "Update config profile error";
        readonly httpCode: 500;
    };
    readonly OAUTH2_PROVIDER_NOT_FOUND: {
        readonly code: "A147";
        readonly message: "OAuth2 provider not found";
        readonly httpCode: 404;
    };
    readonly OAUTH2_AUTHORIZE_ERROR: {
        readonly code: "A148";
        readonly message: "OAuth2 authorize error";
        readonly httpCode: 500;
    };
    readonly NODE_IS_DISABLED: {
        readonly code: "A149";
        readonly message: "Node is disabled";
        readonly httpCode: 400;
    };
    readonly SYNC_ACTIVE_PROFILE_ERROR: {
        readonly code: "A150";
        readonly message: "Sync active profile error";
        readonly httpCode: 500;
    };
    readonly GET_ALL_HOST_TAGS_ERROR: {
        readonly code: "A151";
        readonly message: "Get all host tags error";
        readonly httpCode: 500;
    };
    readonly NAME_OR_CONFIG_REQUIRED: {
        readonly code: "A152";
        readonly message: "Name or config is required";
        readonly httpCode: 400;
    };
    readonly NAME_OR_INBOUNDS_REQUIRED: {
        readonly code: "A153";
        readonly message: "Name or inbounds is required";
        readonly httpCode: 400;
    };
    readonly GET_INTERNAL_SQUAD_ACCESSIBLE_NODES_ERROR: {
        readonly code: "A154";
        readonly message: "Get internal squad accessible nodes error";
        readonly httpCode: 500;
    };
    readonly DELETE_HWID_USER_DEVICES_ERROR: {
        readonly code: "A155";
        readonly message: "Delete hwid user devices error";
        readonly httpCode: 500;
    };
    readonly CREATE_USER_SUBSCRIPTION_REQUEST_HISTORY_ERROR: {
        readonly code: "A156";
        readonly message: "Create user subscription request history error";
        readonly httpCode: 500;
    };
    readonly GET_USER_SUBSCRIPTION_REQUEST_HISTORY_ERROR: {
        readonly code: "A157";
        readonly message: "Get user subscription request history error";
        readonly httpCode: 500;
    };
    readonly GET_ALL_HWID_DEVICES_ERROR: {
        readonly code: "A158";
        readonly message: "Get all hwid devices error";
        readonly httpCode: 500;
    };
    readonly GET_HWID_DEVICES_STATS_ERROR: {
        readonly code: "A159";
        readonly message: "Get hwid devices stats error";
        readonly httpCode: 500;
    };
    readonly GET_USER_SUBSCRIPTION_REQUEST_HISTORY_STATS_ERROR: {
        readonly code: "A160";
        readonly message: "Get user subscription request history stats error";
        readonly httpCode: 500;
    };
    readonly GET_SNIPPETS_ERROR: {
        readonly code: "A161";
        readonly message: "Get snippets error";
        readonly httpCode: 500;
    };
    readonly SNIPPET_NOT_FOUND: {
        readonly code: "A162";
        readonly message: "Snippet not found";
        readonly httpCode: 404;
    };
    readonly DELETE_SNIPPET_BY_NAME_ERROR: {
        readonly code: "A163";
        readonly message: "Delete snippet by name error";
        readonly httpCode: 500;
    };
    readonly SNIPPET_NAME_ALREADY_EXISTS: {
        readonly code: "A164";
        readonly message: "Snippet name already exists";
        readonly httpCode: 400;
    };
    readonly UPDATE_SNIPPET_ERROR: {
        readonly code: "A165";
        readonly message: "Update snippet error";
        readonly httpCode: 500;
    };
    readonly SNIPPET_CANNOT_BE_EMPTY: {
        readonly code: "A166";
        readonly message: "Snippet cannot be empty";
        readonly httpCode: 400;
    };
    readonly SNIPPET_CANNOT_CONTAIN_EMPTY_OBJECTS: {
        readonly code: "A167";
        readonly message: "Snippet cannot contain empty objects";
        readonly httpCode: 400;
    };
    readonly GET_ALL_SUBSCRIPTION_TEMPLATES_ERROR: {
        readonly code: "A168";
        readonly message: "Get all subscription templates error";
        readonly httpCode: 500;
    };
    readonly GET_SUBSCRIPTION_TEMPLATE_BY_UUID_ERROR: {
        readonly code: "A169";
        readonly message: "Get subscription template by UUID error";
        readonly httpCode: 500;
    };
    readonly SUBSCRIPTION_TEMPLATE_NOT_FOUND: {
        readonly code: "A170";
        readonly message: "Subscription template not found";
        readonly httpCode: 404;
    };
    readonly UPDATE_SUBSCRIPTION_TEMPLATE_ERROR: {
        readonly code: "A171";
        readonly message: "Update subscription template error";
        readonly httpCode: 500;
    };
    readonly RESERVED_TEMPLATE_NAME: {
        readonly code: "A172";
        readonly message: "This name is reserved. Please use a different name.";
        readonly httpCode: 400;
    };
    readonly TEMPLATE_JSON_NOT_ALLOWED_FOR_YAML_TEMPLATE: {
        readonly code: "A173";
        readonly message: "Template JSON is not allowed for YAML template";
        readonly httpCode: 400;
    };
    readonly TEMPLATE_YAML_NOT_ALLOWED_FOR_JSON_TEMPLATE: {
        readonly code: "A174";
        readonly message: "Template YAML is not allowed for JSON template";
        readonly httpCode: 400;
    };
    readonly TEMPLATE_JSON_AND_YAML_CANNOT_BE_UPDATED_SIMULTANEOUSLY: {
        readonly code: "A175";
        readonly message: "Template JSON and YAML cannot be updated simultaneously";
        readonly httpCode: 400;
    };
    readonly TEMPLATE_NAME_ALREADY_EXISTS_FOR_THIS_TYPE: {
        readonly code: "A176";
        readonly message: "Template name already exists for this type";
        readonly httpCode: 400;
    };
    readonly DELETE_SUBSCRIPTION_TEMPLATE_ERROR: {
        readonly code: "A177";
        readonly message: "Delete subscription template error";
        readonly httpCode: 500;
    };
    readonly RESERVED_TEMPLATE_CANNOT_BE_DELETED: {
        readonly code: "A178";
        readonly message: "Reserved template cannot be deleted";
        readonly httpCode: 400;
    };
    readonly CREATE_SUBSCRIPTION_TEMPLATE_ERROR: {
        readonly code: "A179";
        readonly message: "Create subscription template error";
        readonly httpCode: 500;
    };
    readonly TEMPLATE_TYPE_NOT_ALLOWED: {
        readonly code: "A180";
        readonly message: "Template type not allowed";
        readonly httpCode: 400;
    };
    readonly GET_EXTERNAL_SQUADS_ERROR: {
        readonly code: "A181";
        readonly message: "Get external squads error";
        readonly httpCode: 500;
    };
    readonly EXTERNAL_SQUAD_NOT_FOUND: {
        readonly code: "A182";
        readonly message: "External squad not found";
        readonly httpCode: 404;
    };
    readonly CREATE_EXTERNAL_SQUAD_ERROR: {
        readonly code: "A183";
        readonly message: "Create external squad error";
        readonly httpCode: 500;
    };
    readonly UPDATE_EXTERNAL_SQUAD_ERROR: {
        readonly code: "A184";
        readonly message: "Update external squad error";
        readonly httpCode: 500;
    };
    readonly DELETE_EXTERNAL_SQUAD_ERROR: {
        readonly code: "A185";
        readonly message: "Delete external squad error";
        readonly httpCode: 500;
    };
    readonly ADD_USERS_TO_EXTERNAL_SQUAD_ERROR: {
        readonly code: "A186";
        readonly message: "Add users to external squad error";
        readonly httpCode: 500;
    };
    readonly REMOVE_USERS_FROM_EXTERNAL_SQUAD_ERROR: {
        readonly code: "A187";
        readonly message: "Remove users from external squad error";
        readonly httpCode: 500;
    };
    readonly GET_EXTERNAL_SQUAD_BY_UUID_ERROR: {
        readonly code: "A188";
        readonly message: "Get external squad by UUID error";
        readonly httpCode: 500;
    };
    readonly EXTERNAL_SQUAD_NAME_ALREADY_EXISTS: {
        readonly code: "A189";
        readonly message: "External squad name already exists";
        readonly httpCode: 400;
    };
    readonly NAME_OR_TEMPLATES_REQUIRED: {
        readonly code: "A190";
        readonly message: "Name or templates are required";
        readonly httpCode: 400;
    };
    readonly PASSKEY_NOT_FOUND: {
        readonly code: "A191";
        readonly message: "Passkey not found";
        readonly httpCode: 404;
    };
    readonly GET_EXODUS_SETTINGS_ERROR: {
        readonly code: "A192";
        readonly message: "Get Exodus settings error";
        readonly httpCode: 500;
    };
    readonly UPDATE_EXODUS_SETTINGS_ERROR: {
        readonly code: "A193";
        readonly message: "Update Exodus settings error";
        readonly httpCode: 500;
    };
    readonly PASSKEYS_NOT_CONFIGURED: {
        readonly code: "A194";
        readonly message: "Passkeys not configured";
        readonly httpCode: 400;
    };
    readonly PASSKEYS_NOT_ENABLED: {
        readonly code: "A195";
        readonly message: "Passkeys not enabled. Please enable it first.";
        readonly httpCode: 400;
    };
    readonly GENERATE_PASSKEY_REGISTRATION_OPTIONS: {
        readonly code: "A196";
        readonly message: "Generate passkey registration options error";
        readonly httpCode: 500;
    };
    readonly VERIFY_PASSKEY_REGISTRATION_ERROR: {
        readonly code: "A197";
        readonly message: "Verify passkey registration error";
        readonly httpCode: 500;
    };
    readonly GET_ACTIVE_PASSKEYS_ERROR: {
        readonly code: "A198";
        readonly message: "Get active passkeys error";
        readonly httpCode: 500;
    };
    readonly DELETE_PASSKEY_ERROR: {
        readonly code: "A199";
        readonly message: "Delete passkey error";
        readonly httpCode: 500;
    };
    readonly VALIDATE_EXODUS_SETTINGS_ERROR: {
        readonly code: "A199";
        readonly message: "Validate Exodus settings error";
        readonly httpCode: 500;
        readonly withMessage: (message: string) => {
            code: string;
            message: string;
            httpCode: number;
        };
    };
    readonly GET_COMPUTED_CONFIG_PROFILE_BY_UUID_ERROR: {
        readonly code: "A200";
        readonly message: "Get computed config profile by UUID error";
        readonly httpCode: 500;
    };
    readonly RESET_NODE_TRAFFIC_ERROR: {
        readonly code: "A201";
        readonly message: "Reset node traffic error";
        readonly httpCode: 500;
    };
    readonly UPDATE_PASSKEY_ERROR: {
        readonly code: "A202";
        readonly message: "Update passkey error";
        readonly httpCode: 500;
    };
    readonly GENERIC_REORDER_ERROR: {
        readonly code: "A203";
        readonly message: "Generic reorder error";
        readonly httpCode: 500;
    };
    readonly HWID_DEVICE_NOT_FOUND: {
        readonly code: "A204";
        readonly message: "HWID device not found";
        readonly httpCode: 404;
    };
    readonly BULK_EXTEND_EXPIRATION_DATE_ERROR: {
        readonly code: "A205";
        readonly message: "Bulk extend expiration date error";
        readonly httpCode: 500;
    };
    readonly SUBSCRIPTION_PAGE_CONFIG_NOT_FOUND: {
        readonly code: "A206";
        readonly message: "Subscription page config not found";
        readonly httpCode: 404;
    };
    readonly GET_SUBSCRIPTION_PAGE_CONFIG_BY_UUID_ERROR: {
        readonly code: "A207";
        readonly message: "Get subscription page config by UUID error";
        readonly httpCode: 500;
    };
    readonly GET_ALL_SUBSCRIPTION_PAGE_CONFIGS_ERROR: {
        readonly code: "A208";
        readonly message: "Get all subscription page configs error";
        readonly httpCode: 500;
    };
    readonly RESERVED_CONFIG_NAME: {
        readonly code: "A209";
        readonly message: "Reserved config name";
        readonly httpCode: 400;
    };
    readonly CONFIG_NAME_ALREADY_EXISTS: {
        readonly code: "A210";
        readonly message: "Config name already exists";
        readonly httpCode: 400;
    };
    readonly UPDATE_SUBSCRIPTION_PAGE_CONFIG_ERROR: {
        readonly code: "A211";
        readonly message: "Update subscription page config error";
        readonly httpCode: 500;
    };
    readonly RESERVED_SUBPAGE_CONFIG_CANT_BE_DELETED: {
        readonly code: "A212";
        readonly message: "Reserved subpage config cannot be deleted";
        readonly httpCode: 400;
    };
    readonly DELETE_SUBSCRIPTION_PAGE_CONFIG_ERROR: {
        readonly code: "A213";
        readonly message: "Delete subscription page config error";
        readonly httpCode: 500;
    };
    readonly CREATE_SUBSCRIPTION_PAGE_CONFIG_ERROR: {
        readonly code: "A214";
        readonly message: "Create subscription page config error";
        readonly httpCode: 500;
    };
    readonly INVALID_SUBSCRIPTION_PAGE_CONFIG: {
        readonly code: "A215";
        readonly message: "Invalid subscription page config";
        readonly httpCode: 400;
    };
    readonly INVALID_EXODUS_INJECTOR: {
        readonly code: "A216";
        readonly message: "Invalid Exodus injector";
        readonly httpCode: 400;
    };
    readonly JOB_CREATION_FAILED: {
        readonly code: "A217";
        readonly message: "Job creation failed";
        readonly httpCode: 500;
    };
    readonly JOB_RESULT_FETCH_FAILED: {
        readonly code: "A218";
        readonly message: "Job result fetch failed or job not found";
        readonly httpCode: 404;
    };
    readonly CONNECTED_NODES_NOT_FOUND: {
        readonly code: "A219";
        readonly message: "Connected nodes not found";
        readonly httpCode: 404;
    };
    readonly GET_ALL_NODE_PLUGINS_ERROR: {
        readonly code: "A219";
        readonly message: "Get all node plugins error";
        readonly httpCode: 500;
    };
    readonly NODE_PLUGIN_NOT_FOUND: {
        readonly code: "A220";
        readonly message: "Node plugin not found";
        readonly httpCode: 404;
    };
    readonly GET_NODE_PLUGIN_BY_UUID_ERROR: {
        readonly code: "A221";
        readonly message: "Get node plugin by UUID error";
        readonly httpCode: 500;
    };
    readonly INVALID_NODE_PLUGIN_CONFIG: {
        readonly code: "A222";
        readonly message: "Invalid node plugin config";
        readonly httpCode: 400;
        readonly withMessage: (message: string) => {
            code: string;
            message: string;
            httpCode: number;
        };
    };
    readonly NODE_PLUGIN_NAME_ALREADY_EXISTS: {
        readonly code: "A223";
        readonly message: "Node plugin name already exists";
        readonly httpCode: 400;
    };
    readonly UPDATE_NODE_PLUGIN_ERROR: {
        readonly code: "A224";
        readonly message: "Update node plugin error";
        readonly httpCode: 500;
    };
    readonly CREATE_NODE_PLUGIN_ERROR: {
        readonly code: "A225";
        readonly message: "Create node plugin error";
        readonly httpCode: 500;
    };
    readonly METADATA_NOT_FOUND: {
        readonly code: "A226";
        readonly message: "Metadata not found";
        readonly httpCode: 404;
    };
    readonly GET_TORRENT_BLOCKER_REPORTS_ERROR: {
        readonly code: "A227";
        readonly message: "Get torrent blocker reports error";
        readonly httpCode: 500;
    };
    readonly UPDATE_HOSTS_ERROR: {
        readonly code: "A228";
        readonly message: "Update hosts error";
        readonly httpCode: 500;
    };
    readonly INVALID_API_TOKEN_SCOPE: {
        readonly code: "A229";
        readonly message: "One or more provided API token scopes are invalid";
        readonly httpCode: 400;
    };
    readonly CREATE_INFRA_BILLING_NODE_MISSING_TARGET: {
        readonly code: "A230";
        readonly message: "Either nodeUuid or name must be provided";
        readonly httpCode: 400;
    };
    readonly CUSTOM_RAW_REMARK_VALIDATION_ERROR: {
        readonly code: "A231";
        readonly message: "Invalid custom raw remark";
        readonly httpCode: 500;
        readonly withMessage: (message: string) => {
            code: string;
            message: string;
            httpCode: number;
        };
    };
    readonly GET_INTERNAL_SQUAD_USAGE_ERROR: {
        readonly code: "A232";
        readonly message: "Get internal squad usage error";
        readonly httpCode: 500;
    };
    readonly ADD_MANY_USERS_TO_INTERNAL_SQUAD_ERROR: {
        readonly code: "A233";
        readonly message: "Add many users to internal squad error";
        readonly httpCode: 500;
    };
    readonly REMOVE_MANY_USERS_FROM_INTERNAL_SQUAD_ERROR: {
        readonly code: "A234";
        readonly message: "Remove many users from internal squad error";
        readonly httpCode: 500;
    };
    readonly GET_STATS_DIGEST_INVALID_RANGE: {
        readonly code: "A235";
        readonly message: "Start date must be before or equal to end date";
        readonly httpCode: 400;
    };
    readonly GET_STATS_DIGEST_ERROR: {
        readonly code: "A236";
        readonly message: "Get stats digest error";
        readonly httpCode: 500;
    };
    readonly SYNC_SNIPPET_ERROR: {
        readonly code: "A237";
        readonly message: "Sync snippet error";
        readonly httpCode: 500;
    };
    readonly NODE_INTEGRATION_NOT_FOUND: {
        readonly code: "A238";
        readonly message: "Node integration not found";
        readonly httpCode: 404;
    };
    readonly GET_ALL_NODE_INTEGRATIONS_ERROR: {
        readonly code: "A239";
        readonly message: "Get all node integrations error";
        readonly httpCode: 500;
    };
    readonly GET_NODE_INTEGRATION_BY_UUID_ERROR: {
        readonly code: "A240";
        readonly message: "Get node integration by UUID error";
        readonly httpCode: 500;
    };
    readonly CREATE_NODE_INTEGRATION_ERROR: {
        readonly code: "A241";
        readonly message: "Create node integration error";
        readonly httpCode: 500;
    };
    readonly UPDATE_NODE_INTEGRATION_ERROR: {
        readonly code: "A242";
        readonly message: "Update node integration error";
        readonly httpCode: 500;
    };
    readonly DELETE_NODE_INTEGRATION_ERROR: {
        readonly code: "A243";
        readonly message: "Delete node integration error";
        readonly httpCode: 500;
    };
    readonly NODE_INTEGRATION_NAME_ALREADY_EXISTS: {
        readonly code: "A244";
        readonly message: "Node integration name already exists";
        readonly httpCode: 400;
    };
    readonly SHARED_LIST_NOT_FOUND: {
        readonly code: "A245";
        readonly message: "Shared list not found";
        readonly httpCode: 404;
    };
    readonly SHARED_LIST_NAME_ALREADY_EXISTS: {
        readonly code: "A246";
        readonly message: "Shared list name already exists";
        readonly httpCode: 400;
    };
    readonly GET_ALL_SHARED_LISTS_ERROR: {
        readonly code: "A247";
        readonly message: "Get all shared lists error";
        readonly httpCode: 500;
    };
    readonly GET_SHARED_LIST_BY_NAME_ERROR: {
        readonly code: "A248";
        readonly message: "Get shared list by name error";
        readonly httpCode: 500;
    };
    readonly CREATE_SHARED_LIST_ERROR: {
        readonly code: "A249";
        readonly message: "Create shared list error";
        readonly httpCode: 500;
    };
    readonly UPDATE_SHARED_LIST_ERROR: {
        readonly code: "A250";
        readonly message: "Update shared list error";
        readonly httpCode: 500;
    };
    readonly DELETE_SHARED_LIST_ERROR: {
        readonly code: "A251";
        readonly message: "Delete shared list error";
        readonly httpCode: 500;
    };
    readonly INVALID_SHARED_LIST_CONFIG: {
        readonly code: "A252";
        readonly message: "Invalid shared list config";
        readonly httpCode: 400;
        readonly withMessage: (message: string) => {
            code: string;
            message: string;
            httpCode: number;
        };
    };
    readonly EVALUATE_VAULT_ERROR: {
        readonly code: "A255";
        readonly message: "Vault evaluation refused";
        readonly httpCode: 429;
    };
    readonly CREATE_SSH_TICKET_ERROR: {
        readonly code: "A254";
        readonly message: "Create SSH ticket error";
        readonly httpCode: 500;
    };
    readonly GET_ENTITY_TAGS_ERROR: {
        readonly code: "A256";
        readonly message: "Get tags error";
        readonly httpCode: 500;
    };
    readonly SET_ENTITY_TAGS_ERROR: {
        readonly code: "A257";
        readonly message: "Set tags error";
        readonly httpCode: 500;
    };
};
//# sourceMappingURL=errors.d.ts.map