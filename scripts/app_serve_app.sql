\c tks;
CREATE TABLE app_serve_apps
(
    id uuid primary key,
    name character varying(50) COLLATE pg_catalog."default",
    contract_id character varying(10) COLLATE pg_catalog."default",
		version character varying(20) COLLATE pg_catalog."default",
    task_type character varying(10) COLLATE pg_catalog."default",
    status character varying(20) COLLATE pg_catalog."default",
    output character varying(1000) COLLATE pg_catalog."default",
    artifact_url character varying(50) COLLATE pg_catalog."default",
    image_url character varying(50) COLLATE pg_catalog."default",
    endpoint_url character varying(50) COLLATE pg_catalog."default",
    target_cluster character varying(10) COLLATE pg_catalog."default",
    profile character varying(10) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    created_at timestamp with time zone
);
