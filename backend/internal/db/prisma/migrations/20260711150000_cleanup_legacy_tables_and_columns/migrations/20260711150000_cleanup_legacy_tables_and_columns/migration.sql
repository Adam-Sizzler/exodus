-- Drop Table admin_sessions
DROP TABLE IF EXISTS public.admin_sessions CASCADE;

-- Drop Column grpc_auth_token from keygen table
ALTER TABLE public.keygen DROP COLUMN IF EXISTS grpc_auth_token;