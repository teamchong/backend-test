CREATE SCHEMA IF NOT EXISTS public AUTHORIZATION postgres;
CREATE TABLE IF NOT EXISTS public.random_user (
    id bigint NOT NULL,
    first_name text NULL,
    last_name text NULL,
    date_of_birth text NULL,
    city text NULL,
    street_name text NULL,
    street_address text NULL,
    zip_code text NULL,
    state text NULL,
    country text NULL,
    lat NUMERIC(28, 15) NULL,
    lng NUMERIC(28, 15) NULL,
    CONSTRAINT pk_id PRIMARY KEY (id)
);
ALTER TABLE public.random_user OWNER TO postgres;
GRANT ALL ON TABLE public.random_user TO postgres;
