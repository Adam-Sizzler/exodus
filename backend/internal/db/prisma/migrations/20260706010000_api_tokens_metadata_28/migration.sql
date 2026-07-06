DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'api_tokens' AND column_name = 'token'
    ) THEN
        DROP INDEX IF EXISTS public.api_tokens_token_key;
        ALTER TABLE public.api_tokens DROP COLUMN token;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'api_tokens' AND column_name = 'token_name'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'api_tokens' AND column_name = 'name'
    ) THEN
        ALTER TABLE public.api_tokens RENAME COLUMN token_name TO name;
    END IF;
END $$;

ALTER TABLE public.api_tokens
    ADD COLUMN IF NOT EXISTS expire_at timestamp(3) without time zone DEFAULT (now() + interval '99999 days'),
    ADD COLUMN IF NOT EXISTS scopes text[] DEFAULT ARRAY['*']::text[];

UPDATE public.api_tokens
SET expire_at = COALESCE(expire_at, created_at + interval '99999 days'),
    scopes = COALESCE(scopes, ARRAY['*']::text[])
WHERE expire_at IS NULL OR scopes IS NULL;

ALTER TABLE public.api_tokens
    ALTER COLUMN expire_at SET NOT NULL,
    ALTER COLUMN scopes SET DEFAULT ARRAY['*']::text[],
    ALTER COLUMN scopes SET NOT NULL;

DO $$
DECLARE
    device_count bigint;
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'hwid_user_devices' AND column_name = 'user_uuid'
    ) THEN
        SELECT count(*) INTO device_count FROM public.hwid_user_devices;

        IF NOT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = 'hwid_user_devices' AND column_name = 'user_id'
        ) THEN
            IF device_count <= 500000 THEN
                ALTER TABLE public.hwid_user_devices ADD COLUMN user_id bigint;

                UPDATE public.hwid_user_devices h
                SET user_id = u.t_id
                FROM public.users u
                WHERE u.uuid = h.user_uuid;

                DELETE FROM public.hwid_user_devices WHERE user_id IS NULL;
                ALTER TABLE public.hwid_user_devices ALTER COLUMN user_id SET NOT NULL;
            ELSE
                TRUNCATE TABLE public.hwid_user_devices;
                ALTER TABLE public.hwid_user_devices ADD COLUMN user_id bigint NOT NULL;
            END IF;
        END IF;

        ALTER TABLE public.hwid_user_devices DROP CONSTRAINT IF EXISTS hwid_user_devices_pkey;
        ALTER TABLE public.hwid_user_devices DROP CONSTRAINT IF EXISTS hwid_user_devices_user_uuid_fkey;
        ALTER TABLE public.hwid_user_devices DROP COLUMN user_uuid;
    END IF;
END $$;

ALTER TABLE public.hwid_user_devices
    ADD COLUMN IF NOT EXISTS request_ip text;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'hwid_user_devices_pkey'
          AND conrelid = 'public.hwid_user_devices'::regclass
    ) THEN
        ALTER TABLE public.hwid_user_devices
            ADD CONSTRAINT hwid_user_devices_pkey PRIMARY KEY (hwid, user_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'hwid_user_devices_user_id_fkey'
          AND conrelid = 'public.hwid_user_devices'::regclass
    ) THEN
        ALTER TABLE public.hwid_user_devices
            ADD CONSTRAINT hwid_user_devices_user_id_fkey
            FOREIGN KEY (user_id) REFERENCES public.users(t_id) ON UPDATE CASCADE ON DELETE CASCADE;
    END IF;
END $$;

TRUNCATE TABLE public.user_subscription_request_history RESTART IDENTITY;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'user_subscription_request_history' AND column_name = 'user_uuid'
    ) THEN
        ALTER TABLE public.user_subscription_request_history DROP CONSTRAINT IF EXISTS user_subscription_request_history_user_uuid_fkey;
        DROP INDEX IF EXISTS public.user_subscription_request_history_user_uuid_idx;
        ALTER TABLE public.user_subscription_request_history DROP COLUMN user_uuid;
    END IF;
END $$;

ALTER TABLE public.user_subscription_request_history
    ADD COLUMN IF NOT EXISTS user_id bigint NOT NULL;

CREATE INDEX IF NOT EXISTS user_subscription_request_history_user_id_idx
    ON public.user_subscription_request_history(user_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'user_subscription_request_history_user_id_fkey'
          AND conrelid = 'public.user_subscription_request_history'::regclass
    ) THEN
        ALTER TABLE public.user_subscription_request_history
            ADD CONSTRAINT user_subscription_request_history_user_id_fkey
            FOREIGN KEY (user_id) REFERENCES public.users(t_id) ON UPDATE CASCADE ON DELETE CASCADE;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS public.user_meta (
    user_id bigint PRIMARY KEY REFERENCES public.users(t_id) ON UPDATE CASCADE ON DELETE CASCADE,
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS public.node_meta (
    node_id bigint PRIMARY KEY REFERENCES public.nodes(id) ON UPDATE CASCADE ON DELETE CASCADE,
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb
);
