CREATE TABLE IF NOT EXISTS public.config_profile_snippets (
    name character varying(255) NOT NULL,
    snippet jsonb NOT NULL,
    created_at timestamp(3) without time zone DEFAULT now() NOT NULL
);

ALTER TABLE public.config_profile_snippets
    ADD COLUMN IF NOT EXISTS name character varying(255),
    ADD COLUMN IF NOT EXISTS snippet jsonb,
    ADD COLUMN IF NOT EXISTS created_at timestamp(3) without time zone DEFAULT now();

UPDATE public.config_profile_snippets
SET created_at = now()
WHERE created_at IS NULL;

ALTER TABLE public.config_profile_snippets
    ALTER COLUMN name TYPE character varying(255),
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN snippet SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT now(),
    ALTER COLUMN created_at SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'config_profile_snippets_pkey'
          AND conrelid = 'public.config_profile_snippets'::regclass
    ) THEN
        ALTER TABLE public.config_profile_snippets
            ADD CONSTRAINT config_profile_snippets_pkey PRIMARY KEY (name);
    END IF;
END $$;
