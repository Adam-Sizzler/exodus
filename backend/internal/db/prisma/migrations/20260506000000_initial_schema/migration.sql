CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;
CREATE TABLE public.admin (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    username text NOT NULL,
    password_hash text NOT NULL,
    role text NOT NULL,
    session_ttl_minutes integer DEFAULT 60 NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.admin_sessions (
    session_token text NOT NULL,
    admin_uuid uuid NOT NULL,
    expires_at bigint NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.api_tokens (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    token text NOT NULL,
    token_name text NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.config_profile_inbounds (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    profile_uuid uuid NOT NULL,
    tag text NOT NULL,
    type text NOT NULL,
    network text,
    security text,
    port integer,
    raw_inbound jsonb
);
CREATE TABLE public.config_profile_inbounds_to_nodes (
    config_profile_inbound_uuid uuid CONSTRAINT config_profile_inbounds_to__config_profile_inbound_uui_not_null NOT NULL,
    node_uuid uuid NOT NULL
);
CREATE TABLE public.config_profile_snippets (
    name character varying(255) NOT NULL,
    snippet jsonb NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.config_profiles (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    view_position integer NOT NULL,
    name text NOT NULL,
    config json NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.config_profiles_view_position_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.config_profiles_view_position_seq OWNED BY public.config_profiles.view_position;
CREATE TABLE public.exodus_settings (
    id integer DEFAULT 1 CONSTRAINT cerberus_settings_id_not_null NOT NULL,
    passkey_settings jsonb,
    oauth2_settings jsonb,
    tg_auth_settings jsonb,
    password_settings jsonb,
    branding_settings jsonb,
    modules_settings jsonb
);
CREATE TABLE public.external_squads (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    view_position integer NOT NULL,
    name character varying(30) NOT NULL,
    subscription_settings jsonb,
    host_overrides jsonb,
    response_headers jsonb,
    hwid_settings jsonb,
    custom_remarks jsonb,
    subpage_config_uuid uuid,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.external_squads_templates (
    external_squad_uuid uuid NOT NULL,
    template_uuid uuid NOT NULL,
    template_type text NOT NULL
);
CREATE SEQUENCE public.external_squads_view_position_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.external_squads_view_position_seq OWNED BY public.external_squads.view_position;
CREATE TABLE public.hosts (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    view_position integer NOT NULL,
    remark character varying(50) NOT NULL,
    address text NOT NULL,
    port integer NOT NULL,
    path text,
    sni text,
    host text,
    alpn text,
    fingerprint text,
    security_layer text DEFAULT 'DEFAULT'::text NOT NULL,
    mux_params jsonb,
    sockopt_params jsonb,
    is_disabled boolean DEFAULT false NOT NULL,
    server_description character varying(30),
    vless_route_id integer,
    allow_insecure boolean DEFAULT false NOT NULL,
    shuffle_host boolean DEFAULT false NOT NULL,
    mihomo_x25519 boolean DEFAULT false NOT NULL,
    xray_json_template_uuid uuid,
    keep_sni_blank boolean DEFAULT false NOT NULL,
    exclude_from_subscription_types text[] DEFAULT ARRAY[]::text[] NOT NULL,
    tag text,
    is_hidden boolean DEFAULT false NOT NULL,
    override_sni_from_address boolean DEFAULT false NOT NULL,
    config_profile_uuid uuid,
    config_profile_inbound_uuid uuid,
    singbox_mux_params jsonb,
    clash_mux_params jsonb,
    selector_nodes_first boolean DEFAULT false NOT NULL
);
CREATE TABLE public.hosts_to_nodes (
    host_uuid uuid NOT NULL,
    node_uuid uuid NOT NULL
);
CREATE SEQUENCE public.hosts_view_position_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.hosts_view_position_seq OWNED BY public.hosts.view_position;
CREATE TABLE public.hwid_user_devices (
    hwid text NOT NULL,
    user_uuid uuid NOT NULL,
    platform text,
    os_version text,
    device_model text,
    user_agent text,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.infra_billing_history (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    provider_uuid uuid NOT NULL,
    amount double precision NOT NULL,
    billed_at timestamp(3) without time zone NOT NULL
);
CREATE TABLE public.infra_billing_nodes (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    node_uuid uuid NOT NULL,
    provider_uuid uuid NOT NULL,
    next_billing_at timestamp(3) without time zone NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.infra_providers (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    favicon_link text,
    login_url text,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.internal_squad_host_exclusions (
    host_uuid uuid NOT NULL,
    squad_uuid uuid NOT NULL
);
CREATE TABLE public.internal_squad_inbounds (
    internal_squad_uuid uuid NOT NULL,
    inbound_uuid uuid NOT NULL
);
CREATE TABLE public.internal_squad_members (
    internal_squad_uuid uuid NOT NULL,
    user_id bigint NOT NULL
);
CREATE TABLE public.internal_squads (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    view_position integer NOT NULL,
    name text NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.internal_squads_view_position_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.internal_squads_view_position_seq OWNED BY public.internal_squads.view_position;
CREATE TABLE public.keygen (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    priv_key text NOT NULL,
    pub_key text NOT NULL,
    ca_cert text,
    ca_key text,
    client_cert text,
    client_key text,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    grpc_auth_token text DEFAULT encode(public.gen_random_bytes(32), 'hex'::text) NOT NULL
);
CREATE TABLE public.modules_settings (
    id integer DEFAULT 1 NOT NULL,
    haproxy_enabled boolean DEFAULT false NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.nodes (
    id bigint NOT NULL,
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    address text NOT NULL,
    port integer,
    api_schema text DEFAULT 'grpc'::text NOT NULL,
    api_path text DEFAULT ''::text NOT NULL,
    active_config_profile_uuid uuid,
    is_connected boolean DEFAULT false NOT NULL,
    is_connecting boolean DEFAULT false NOT NULL,
    is_disabled boolean DEFAULT false NOT NULL,
    last_status_change timestamp(3) without time zone,
    last_status_message text,
    singbox_version text,
    node_version text,
    singbox_uptime text DEFAULT '0'::text CONSTRAINT nodes_xray_uptime_not_null NOT NULL,
    users_online integer DEFAULT 0,
    consumption_multiplier bigint DEFAULT 1000000000 NOT NULL,
    is_traffic_tracking_active boolean DEFAULT false NOT NULL,
    traffic_reset_day integer DEFAULT 1,
    traffic_limit_bytes bigint DEFAULT 0,
    traffic_used_bytes bigint DEFAULT 0,
    notify_percent integer DEFAULT 0,
    provider_uuid uuid,
    view_position integer NOT NULL,
    country_code text DEFAULT 'XX'::text NOT NULL,
    tags text[] DEFAULT ARRAY[]::text[] NOT NULL,
    cpu_count integer,
    cpu_model text,
    total_ram text,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.nodes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.nodes_id_seq OWNED BY public.nodes.id;
CREATE TABLE public.nodes_traffic_usage_history (
    id bigint NOT NULL,
    node_uuid uuid NOT NULL,
    traffic_bytes bigint NOT NULL,
    reset_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.nodes_traffic_usage_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.nodes_traffic_usage_history_id_seq OWNED BY public.nodes_traffic_usage_history.id;
CREATE TABLE public.nodes_usage_history (
    node_uuid uuid NOT NULL,
    download_bytes bigint NOT NULL,
    upload_bytes bigint NOT NULL,
    total_bytes bigint NOT NULL,
    created_at timestamp(3) without time zone DEFAULT date_trunc('hour'::text, now()) NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.nodes_user_usage_history (
    node_id bigint NOT NULL,
    user_id bigint NOT NULL,
    total_bytes bigint NOT NULL,
    created_at date DEFAULT CURRENT_DATE NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.nodes_view_position_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.nodes_view_position_seq OWNED BY public.nodes.view_position;
CREATE TABLE public.passkeys (
    id text NOT NULL,
    admin_uuid uuid NOT NULL,
    public_key bytea NOT NULL,
    counter bigint NOT NULL,
    device_type text NOT NULL,
    backed_up boolean NOT NULL,
    transports text,
    passkey_provider text,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.srs_lists (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    tag text NOT NULL,
    format text DEFAULT 'binary'::text NOT NULL,
    url text NOT NULL,
    update_interval text DEFAULT '1d'::text NOT NULL,
    path text,
    file_name text NOT NULL,
    view_position integer NOT NULL,
    is_available boolean DEFAULT false NOT NULL,
    last_checked_at timestamp(3) without time zone,
    last_error text,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    is_enabled boolean DEFAULT true NOT NULL
);
CREATE SEQUENCE public.srs_lists_view_position_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.srs_lists_view_position_seq OWNED BY public.srs_lists.view_position;
CREATE TABLE public.sub_nodes (
    id bigint NOT NULL,
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    address text NOT NULL,
    port integer,
    api_schema text DEFAULT 'mtls'::text NOT NULL,
    api_path text DEFAULT '/'::text NOT NULL,
    is_connected boolean DEFAULT false NOT NULL,
    is_connecting boolean DEFAULT false NOT NULL,
    is_disabled boolean DEFAULT false NOT NULL,
    provider_uuid uuid,
    view_position integer NOT NULL,
    tags text[] DEFAULT ARRAY[]::text[] NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    last_status_change timestamp(3) without time zone,
    last_status_message text,
    public_domain text
);
CREATE SEQUENCE public.sub_nodes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.sub_nodes_id_seq OWNED BY public.sub_nodes.id;
CREATE TABLE public.sub_nodes_to_subscription_page_config (
    node_uuid uuid NOT NULL,
    subpage_config_uuid uuid CONSTRAINT sub_nodes_to_subscription_page_con_subpage_config_uuid_not_null NOT NULL
);
CREATE SEQUENCE public.sub_nodes_view_position_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.sub_nodes_view_position_seq OWNED BY public.sub_nodes.view_position;
CREATE TABLE public.subscription_page_config (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    view_position integer NOT NULL,
    name character varying(30) NOT NULL,
    config jsonb NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.subscription_page_config_view_position_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.subscription_page_config_view_position_seq OWNED BY public.subscription_page_config.view_position;
CREATE TABLE public.subscription_settings (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    profile_title text NOT NULL,
    support_link text NOT NULL,
    profile_update_interval integer NOT NULL,
    address text,
    port integer,
    api_schema text,
    api_path text,
    is_profile_webpage_url_enabled boolean DEFAULT true NOT NULL,
    serve_json_at_base_subscription boolean DEFAULT false NOT NULL,
    happ_announce text,
    happ_routing text,
    is_show_custom_remarks boolean DEFAULT true NOT NULL,
    custom_remarks jsonb NOT NULL,
    custom_response_headers jsonb,
    randomize_hosts boolean DEFAULT false NOT NULL,
    response_rules jsonb,
    hwid_settings jsonb,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.subscription_templates (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    view_position integer NOT NULL,
    name character varying(255) DEFAULT 'Default'::character varying NOT NULL,
    template_type text NOT NULL,
    template_yaml text,
    template_json json,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.subscription_templates_view_position_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.subscription_templates_view_position_seq OWNED BY public.subscription_templates.view_position;
CREATE TABLE public.user_subscription_request_history (
    id bigint NOT NULL,
    user_uuid uuid NOT NULL,
    request_ip text,
    user_agent text,
    request_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.user_subscription_request_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.user_subscription_request_history_id_seq OWNED BY public.user_subscription_request_history.id;
CREATE TABLE public.notification_delivery_queue (
    id bigserial NOT NULL,
    event jsonb NOT NULL,
    attempts integer DEFAULT 0 NOT NULL,
    next_attempt_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    last_error text,
    failed_at timestamp(3) without time zone,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.user_traffic (
    t_id bigint NOT NULL,
    used_traffic_bytes bigint DEFAULT 0 NOT NULL,
    lifetime_used_traffic_bytes bigint DEFAULT 0 NOT NULL,
    online_at timestamp(3) without time zone,
    last_connected_node_uuid uuid,
    first_connected_at timestamp(3) without time zone
);
CREATE TABLE public.users (
    t_id bigint NOT NULL,
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    short_uuid text NOT NULL,
    username text NOT NULL,
    status character varying(10) DEFAULT 'ACTIVE'::character varying NOT NULL,
    traffic_limit_bytes bigint DEFAULT 0 NOT NULL,
    traffic_limit_strategy text DEFAULT 'NO_RESET'::text NOT NULL,
    expire_at timestamp(3) without time zone NOT NULL,
    sub_last_user_agent text,
    sub_last_opened_at timestamp(3) without time zone,
    last_traffic_reset_at timestamp(3) without time zone,
    sub_revoked_at timestamp(3) without time zone,
    trojan_password text NOT NULL,
    vless_uuid uuid NOT NULL,
    ss_password text NOT NULL,
    description text,
    tag text,
    telegram_id bigint,
    email text,
    hwid_device_limit integer,
    external_squad_uuid uuid,
    last_triggered_threshold integer DEFAULT 0 NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL,
    updated_at timestamp(3) without time zone DEFAULT now() NOT NULL
);
CREATE SEQUENCE public.users_t_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.users_t_id_seq OWNED BY public.users.t_id;
ALTER TABLE ONLY public.config_profiles ALTER COLUMN view_position SET DEFAULT nextval('public.config_profiles_view_position_seq'::regclass);
ALTER TABLE ONLY public.external_squads ALTER COLUMN view_position SET DEFAULT nextval('public.external_squads_view_position_seq'::regclass);
ALTER TABLE ONLY public.hosts ALTER COLUMN view_position SET DEFAULT nextval('public.hosts_view_position_seq'::regclass);
ALTER TABLE ONLY public.internal_squads ALTER COLUMN view_position SET DEFAULT nextval('public.internal_squads_view_position_seq'::regclass);
ALTER TABLE ONLY public.nodes ALTER COLUMN id SET DEFAULT nextval('public.nodes_id_seq'::regclass);
ALTER TABLE ONLY public.nodes ALTER COLUMN view_position SET DEFAULT nextval('public.nodes_view_position_seq'::regclass);
ALTER TABLE ONLY public.nodes_traffic_usage_history ALTER COLUMN id SET DEFAULT nextval('public.nodes_traffic_usage_history_id_seq'::regclass);
ALTER TABLE ONLY public.srs_lists ALTER COLUMN view_position SET DEFAULT nextval('public.srs_lists_view_position_seq'::regclass);
ALTER TABLE ONLY public.sub_nodes ALTER COLUMN id SET DEFAULT nextval('public.sub_nodes_id_seq'::regclass);
ALTER TABLE ONLY public.sub_nodes ALTER COLUMN view_position SET DEFAULT nextval('public.sub_nodes_view_position_seq'::regclass);
ALTER TABLE ONLY public.subscription_page_config ALTER COLUMN view_position SET DEFAULT nextval('public.subscription_page_config_view_position_seq'::regclass);
ALTER TABLE ONLY public.subscription_templates ALTER COLUMN view_position SET DEFAULT nextval('public.subscription_templates_view_position_seq'::regclass);
ALTER TABLE ONLY public.user_subscription_request_history ALTER COLUMN id SET DEFAULT nextval('public.user_subscription_request_history_id_seq'::regclass);
ALTER TABLE ONLY public.users ALTER COLUMN t_id SET DEFAULT nextval('public.users_t_id_seq'::regclass);
ALTER TABLE ONLY public.admin
    ADD CONSTRAINT admin_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.admin_sessions
    ADD CONSTRAINT admin_sessions_pkey PRIMARY KEY (session_token);
ALTER TABLE ONLY public.api_tokens
    ADD CONSTRAINT api_tokens_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.config_profile_inbounds
    ADD CONSTRAINT config_profile_inbounds_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.config_profile_inbounds_to_nodes
    ADD CONSTRAINT config_profile_inbounds_to_nodes_pkey PRIMARY KEY (config_profile_inbound_uuid, node_uuid);
ALTER TABLE ONLY public.config_profile_snippets
    ADD CONSTRAINT config_profile_snippets_pkey PRIMARY KEY (name);
ALTER TABLE ONLY public.config_profiles
    ADD CONSTRAINT config_profiles_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.exodus_settings
    ADD CONSTRAINT exodus_settings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.external_squads
    ADD CONSTRAINT external_squads_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.external_squads_templates
    ADD CONSTRAINT external_squads_templates_pkey PRIMARY KEY (external_squad_uuid, template_type);
ALTER TABLE ONLY public.hosts
    ADD CONSTRAINT hosts_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.hosts_to_nodes
    ADD CONSTRAINT hosts_to_nodes_pkey PRIMARY KEY (host_uuid, node_uuid);
ALTER TABLE ONLY public.hwid_user_devices
    ADD CONSTRAINT hwid_user_devices_pkey PRIMARY KEY (hwid, user_uuid);
ALTER TABLE ONLY public.infra_billing_history
    ADD CONSTRAINT infra_billing_history_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.infra_billing_nodes
    ADD CONSTRAINT infra_billing_nodes_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.infra_providers
    ADD CONSTRAINT infra_providers_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.internal_squad_host_exclusions
    ADD CONSTRAINT internal_squad_host_exclusions_pkey PRIMARY KEY (host_uuid, squad_uuid);
ALTER TABLE ONLY public.internal_squad_inbounds
    ADD CONSTRAINT internal_squad_inbounds_pkey PRIMARY KEY (internal_squad_uuid, inbound_uuid);
ALTER TABLE ONLY public.internal_squad_members
    ADD CONSTRAINT internal_squad_members_pkey PRIMARY KEY (internal_squad_uuid, user_id);
ALTER TABLE ONLY public.internal_squads
    ADD CONSTRAINT internal_squads_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.keygen
    ADD CONSTRAINT keygen_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.modules_settings
    ADD CONSTRAINT modules_settings_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.nodes
    ADD CONSTRAINT nodes_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.nodes_traffic_usage_history
    ADD CONSTRAINT nodes_traffic_usage_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.nodes_usage_history
    ADD CONSTRAINT nodes_usage_history_pkey PRIMARY KEY (node_uuid, created_at);
ALTER TABLE ONLY public.nodes_user_usage_history
    ADD CONSTRAINT nodes_user_usage_history_pkey PRIMARY KEY (node_id, created_at, user_id);
ALTER TABLE ONLY public.passkeys
    ADD CONSTRAINT passkeys_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.srs_lists
    ADD CONSTRAINT srs_lists_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.sub_nodes
    ADD CONSTRAINT sub_nodes_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.sub_nodes_to_subscription_page_config
    ADD CONSTRAINT sub_nodes_to_subscription_page_config_pkey PRIMARY KEY (node_uuid, subpage_config_uuid);
ALTER TABLE ONLY public.subscription_page_config
    ADD CONSTRAINT subscription_page_config_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.subscription_settings
    ADD CONSTRAINT subscription_settings_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.subscription_templates
    ADD CONSTRAINT subscription_templates_pkey PRIMARY KEY (uuid);
ALTER TABLE ONLY public.user_subscription_request_history
    ADD CONSTRAINT user_subscription_request_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.user_traffic
    ADD CONSTRAINT user_traffic_pkey PRIMARY KEY (t_id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (t_id);
ALTER TABLE ONLY public.notification_delivery_queue
    ADD CONSTRAINT notification_delivery_queue_pkey PRIMARY KEY (id);
CREATE INDEX admin_sessions_admin_uuid_idx ON public.admin_sessions USING btree (admin_uuid);
CREATE INDEX admin_sessions_expires_at_idx ON public.admin_sessions USING btree (expires_at);
CREATE UNIQUE INDEX admin_username_key ON public.admin USING btree (username);
CREATE UNIQUE INDEX api_tokens_token_key ON public.api_tokens USING btree (token);
CREATE INDEX config_profile_inbounds_profile_uuid_uuid_idx ON public.config_profile_inbounds USING btree (profile_uuid, uuid);
CREATE UNIQUE INDEX config_profile_inbounds_tag_key ON public.config_profile_inbounds USING btree (tag);
CREATE UNIQUE INDEX config_profiles_name_key ON public.config_profiles USING btree (name);
CREATE UNIQUE INDEX external_squads_name_key ON public.external_squads USING btree (name);
CREATE INDEX infra_billing_nodes_next_billing_at_idx ON public.infra_billing_nodes USING btree (next_billing_at);
CREATE UNIQUE INDEX infra_billing_nodes_node_uuid_provider_uuid_key ON public.infra_billing_nodes USING btree (node_uuid, provider_uuid);
CREATE UNIQUE INDEX infra_providers_name_key ON public.infra_providers USING btree (name);
CREATE INDEX internal_squad_members_internal_squad_uuid_idx ON public.internal_squad_members USING btree (internal_squad_uuid);
CREATE INDEX internal_squad_members_user_id_idx ON public.internal_squad_members USING btree (user_id);
CREATE UNIQUE INDEX internal_squads_name_key ON public.internal_squads USING btree (name);
CREATE UNIQUE INDEX nodes_address_key ON public.nodes USING btree (address);
CREATE INDEX nodes_id_idx ON public.nodes USING btree (id);
CREATE UNIQUE INDEX nodes_id_key ON public.nodes USING btree (id);
CREATE UNIQUE INDEX nodes_name_key ON public.nodes USING btree (name);
CREATE INDEX nodes_usage_history_node_uuid_created_at_idx ON public.nodes_usage_history USING btree (node_uuid, created_at DESC);
CREATE INDEX notification_delivery_queue_due_idx ON public.notification_delivery_queue USING btree (failed_at, next_attempt_at, id);
CREATE INDEX passkeys_admin_uuid_idx ON public.passkeys USING btree (admin_uuid);
CREATE INDEX passkeys_id_idx ON public.passkeys USING btree (id);
CREATE UNIQUE INDEX srs_lists_file_name_key ON public.srs_lists USING btree (file_name);
CREATE UNIQUE INDEX srs_lists_tag_key ON public.srs_lists USING btree (tag);
CREATE UNIQUE INDEX srs_lists_url_key ON public.srs_lists USING btree (url);
CREATE UNIQUE INDEX sub_nodes_address_port_api_path_key ON public.sub_nodes USING btree (address, port, api_path);
CREATE INDEX sub_nodes_id_idx ON public.sub_nodes USING btree (id);
CREATE UNIQUE INDEX sub_nodes_id_key ON public.sub_nodes USING btree (id);
CREATE UNIQUE INDEX sub_nodes_name_key ON public.sub_nodes USING btree (name);
CREATE UNIQUE INDEX sub_nodes_to_subscription_page_config_node_uuid_key ON public.sub_nodes_to_subscription_page_config USING btree (node_uuid);
CREATE UNIQUE INDEX subscription_page_config_name_key ON public.subscription_page_config USING btree (name);
CREATE UNIQUE INDEX subscription_templates_template_type_name_key ON public.subscription_templates USING btree (template_type, name);
CREATE INDEX user_subscription_request_history_request_at_idx ON public.user_subscription_request_history USING btree (request_at);
CREATE INDEX user_subscription_request_history_user_uuid_idx ON public.user_subscription_request_history USING btree (user_uuid);
CREATE INDEX users_created_at_idx ON public.users USING btree (created_at);
CREATE UNIQUE INDEX users_short_uuid_key ON public.users USING btree (short_uuid);
CREATE UNIQUE INDEX users_username_key ON public.users USING btree (username);
CREATE UNIQUE INDEX users_uuid_key ON public.users USING btree (uuid);
ALTER TABLE ONLY public.admin_sessions
    ADD CONSTRAINT admin_sessions_admin_uuid_fkey FOREIGN KEY (admin_uuid) REFERENCES public.admin(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.config_profile_inbounds
    ADD CONSTRAINT config_profile_inbounds_profile_uuid_fkey FOREIGN KEY (profile_uuid) REFERENCES public.config_profiles(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.config_profile_inbounds_to_nodes
    ADD CONSTRAINT config_profile_inbounds_to_nodes_config_profile_inbound_uuid_fk FOREIGN KEY (config_profile_inbound_uuid) REFERENCES public.config_profile_inbounds(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.config_profile_inbounds_to_nodes
    ADD CONSTRAINT config_profile_inbounds_to_nodes_node_uuid_fkey FOREIGN KEY (node_uuid) REFERENCES public.nodes(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.external_squads
    ADD CONSTRAINT external_squads_subpage_config_uuid_fkey FOREIGN KEY (subpage_config_uuid) REFERENCES public.subscription_page_config(uuid) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE ONLY public.external_squads_templates
    ADD CONSTRAINT external_squads_templates_external_squad_uuid_fkey FOREIGN KEY (external_squad_uuid) REFERENCES public.external_squads(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.external_squads_templates
    ADD CONSTRAINT external_squads_templates_template_uuid_fkey FOREIGN KEY (template_uuid) REFERENCES public.subscription_templates(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.hosts
    ADD CONSTRAINT hosts_config_profile_inbound_uuid_fkey FOREIGN KEY (config_profile_inbound_uuid) REFERENCES public.config_profile_inbounds(uuid) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE ONLY public.hosts
    ADD CONSTRAINT hosts_config_profile_uuid_fkey FOREIGN KEY (config_profile_uuid) REFERENCES public.config_profiles(uuid) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE ONLY public.hosts_to_nodes
    ADD CONSTRAINT hosts_to_nodes_host_uuid_fkey FOREIGN KEY (host_uuid) REFERENCES public.hosts(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.hosts_to_nodes
    ADD CONSTRAINT hosts_to_nodes_node_uuid_fkey FOREIGN KEY (node_uuid) REFERENCES public.nodes(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.hosts
    ADD CONSTRAINT hosts_xray_json_template_uuid_fkey FOREIGN KEY (xray_json_template_uuid) REFERENCES public.subscription_templates(uuid) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE ONLY public.hwid_user_devices
    ADD CONSTRAINT hwid_user_devices_user_uuid_fkey FOREIGN KEY (user_uuid) REFERENCES public.users(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.infra_billing_history
    ADD CONSTRAINT infra_billing_history_provider_uuid_fkey FOREIGN KEY (provider_uuid) REFERENCES public.infra_providers(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.infra_billing_nodes
    ADD CONSTRAINT infra_billing_nodes_node_uuid_fkey FOREIGN KEY (node_uuid) REFERENCES public.nodes(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.infra_billing_nodes
    ADD CONSTRAINT infra_billing_nodes_provider_uuid_fkey FOREIGN KEY (provider_uuid) REFERENCES public.infra_providers(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.internal_squad_host_exclusions
    ADD CONSTRAINT internal_squad_host_exclusions_host_uuid_fkey FOREIGN KEY (host_uuid) REFERENCES public.hosts(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.internal_squad_host_exclusions
    ADD CONSTRAINT internal_squad_host_exclusions_squad_uuid_fkey FOREIGN KEY (squad_uuid) REFERENCES public.internal_squads(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.internal_squad_inbounds
    ADD CONSTRAINT internal_squad_inbounds_inbound_uuid_fkey FOREIGN KEY (inbound_uuid) REFERENCES public.config_profile_inbounds(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.internal_squad_inbounds
    ADD CONSTRAINT internal_squad_inbounds_internal_squad_uuid_fkey FOREIGN KEY (internal_squad_uuid) REFERENCES public.internal_squads(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.internal_squad_members
    ADD CONSTRAINT internal_squad_members_internal_squad_uuid_fkey FOREIGN KEY (internal_squad_uuid) REFERENCES public.internal_squads(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.internal_squad_members
    ADD CONSTRAINT internal_squad_members_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(t_id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.nodes
    ADD CONSTRAINT nodes_active_config_profile_uuid_fkey FOREIGN KEY (active_config_profile_uuid) REFERENCES public.config_profiles(uuid) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE ONLY public.nodes
    ADD CONSTRAINT nodes_provider_uuid_fkey FOREIGN KEY (provider_uuid) REFERENCES public.infra_providers(uuid) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE ONLY public.nodes_traffic_usage_history
    ADD CONSTRAINT nodes_traffic_usage_history_node_uuid_fkey FOREIGN KEY (node_uuid) REFERENCES public.nodes(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.nodes_usage_history
    ADD CONSTRAINT nodes_usage_history_node_uuid_fkey FOREIGN KEY (node_uuid) REFERENCES public.nodes(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.nodes_user_usage_history
    ADD CONSTRAINT nodes_user_usage_history_node_id_fkey FOREIGN KEY (node_id) REFERENCES public.nodes(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.nodes_user_usage_history
    ADD CONSTRAINT nodes_user_usage_history_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(t_id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.passkeys
    ADD CONSTRAINT passkeys_admin_uuid_fkey FOREIGN KEY (admin_uuid) REFERENCES public.admin(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.sub_nodes
    ADD CONSTRAINT sub_nodes_provider_uuid_fkey FOREIGN KEY (provider_uuid) REFERENCES public.infra_providers(uuid) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE ONLY public.sub_nodes_to_subscription_page_config
    ADD CONSTRAINT sub_nodes_to_subscription_page_config_node_uuid_fkey FOREIGN KEY (node_uuid) REFERENCES public.sub_nodes(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.sub_nodes_to_subscription_page_config
    ADD CONSTRAINT sub_nodes_to_subscription_page_config_subpage_config_uuid_fkey FOREIGN KEY (subpage_config_uuid) REFERENCES public.subscription_page_config(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.user_subscription_request_history
    ADD CONSTRAINT user_subscription_request_history_user_uuid_fkey FOREIGN KEY (user_uuid) REFERENCES public.users(uuid) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.user_traffic
    ADD CONSTRAINT user_traffic_last_connected_node_uuid_fkey FOREIGN KEY (last_connected_node_uuid) REFERENCES public.nodes(uuid) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE ONLY public.user_traffic
    ADD CONSTRAINT user_traffic_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.users(t_id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_external_squad_uuid_fkey FOREIGN KEY (external_squad_uuid) REFERENCES public.external_squads(uuid) ON UPDATE CASCADE ON DELETE SET NULL;
