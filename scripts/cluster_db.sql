\c tks;
CREATE TABLE clusters
(
    name character varying(50) COLLATE pg_catalog."default",
    id uuid primary key,
    contract_id uuid,
    csp_id uuid,
    status integer,
    ssh_key_name character varying(50) COLLATE pg_catalog."default",
    region character varying(50) COLLATE pg_catalog."default",
    num_of_az integer,
    machine_type character varying(50) COLLATE pg_catalog."default",
    min_size_per_az integer,
    max_size_per_az integer,
    kubeconfig character varying(1000) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    created_at timestamp with time zone
);
