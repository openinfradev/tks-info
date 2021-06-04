CREATE TABLE csp_infos
(
    id uuid primary key,
    contract_id uuid,
    name character varying(50) COLLATE pg_catalog."default",
    auth character varying(200) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    created_at timestamp with time zone
);
