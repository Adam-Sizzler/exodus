-- Performance indexes for hot paths

-- 1. Subscription request history FIFO cleanup (ORDER BY request_at DESC, id DESC LIMIT 24)
CREATE INDEX IF NOT EXISTS "idx_user_sub_req_history_user_date"
    ON public.user_subscription_request_history ("user_id", "request_at" DESC, "id" DESC);

-- 2. User watchdog 30s check (WHERE status IN ('ACTIVE', 'LIMITED') AND expire_at < CURRENT_TIMESTAMP)
CREATE INDEX IF NOT EXISTS "idx_users_active_expired"
    ON public.users ("expire_at")
    WHERE "status" IN ('ACTIVE', 'LIMITED');

-- 3. User watchdog 45s check (WHERE status = 'ACTIVE' AND traffic_limit_bytes > 0)
CREATE INDEX IF NOT EXISTS "idx_users_active_traffic_limit"
    ON public.users ("id")
    WHERE "status" = 'ACTIVE' AND "traffic_limit_bytes" > 0;
