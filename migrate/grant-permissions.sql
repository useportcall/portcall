-- Connect to main_portcall_db
\c main_portcall_db

-- Grant all privileges on the database
GRANT ALL PRIVILEGES ON DATABASE main_portcall_db TO portcall;

-- Grant usage on schema
GRANT USAGE ON SCHEMA public TO portcall;

-- Grant all privileges on schema public
GRANT ALL PRIVILEGES ON SCHEMA public TO portcall;

-- Grant all on all tables (existing)
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO portcall;

-- Grant all on all sequences (existing)
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO portcall;

-- Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO portcall;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO portcall;

-- Connect to keycloak database
\c keycloak

-- Grant all privileges on the database
GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak_user;

-- Grant usage on schema
GRANT USAGE ON SCHEMA public TO keycloak_user;

-- Grant all privileges on schema public
GRANT ALL PRIVILEGES ON SCHEMA public TO keycloak_user;

-- Grant all on all tables (existing)
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO keycloak_user;

-- Grant all on all sequences (existing)
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO keycloak_user;

-- Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO keycloak_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO keycloak_user;
