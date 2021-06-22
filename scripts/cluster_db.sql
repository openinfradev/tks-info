CREATE TABLE clusters
(
    name character varying(50) COLLATE pg_catalog."default",
    id uuid primary key,
    contract_id uuid,
    csp_id uuid,
    status integer,
    master_flavor character varying(50) COLLATE pg_catalog."default",
    master_replicas integer,
    master_root_size bigint,
    worker_flavor character varying(50) COLLATE pg_catalog."default",
    worker_replicas integer,
    worker_root_size bigint,
    k8s_version character varying(50) COLLATE pg_catalog."default",
    kubeconfig character varying(1000) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    created_at timestamp with time zone
);
