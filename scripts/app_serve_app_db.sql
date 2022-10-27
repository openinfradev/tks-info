\c tks;
CREATE TABLE app_serve_apps
(
    id uuid primary key,
    name character varying(50) COLLATE pg_catalog."default",
    contract_id character varying(10) COLLATE pg_catalog."default",
    task_type character varying(10) COLLATE pg_catalog."default",
    endpoint_url character varying(300) COLLATE pg_catalog."default",
    target_cluster_id character varying(10) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    created_at timestamp with time zone
);
CREATE TABLE app_serve_app_tasks
(
    id uuid primary key,
    app_serve_app_id uuid,
    version character varying(20) COLLATE pg_catalog."default",
    status character varying(20) COLLATE pg_catalog."default",
    output character varying(10000) COLLATE pg_catalog."default",
    artifact_url character varying(300) COLLATE pg_catalog."default",
    image_url character varying(300) COLLATE pg_catalog."default",
    profile character varying(20) COLLATE pg_catalog."default",
    helm_revision integer,
    updated_at timestamp with time zone,
    created_at timestamp with time zone,
    FOREIGN KEY (app_serve_app_id)
    REFERENCES app_serve_apps(ID) ON UPDATE CASCADE ON DELETE RESTRICT
);
