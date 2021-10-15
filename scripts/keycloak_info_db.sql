\c tks;
CREATE TABLE keycloak_infos
(
    id uuid primary key,
    cluster_id uuid,
    realm character varying(100) COLLATE pg_catalog."default",
    client_id character varying(100) COLLATE pg_catalog."default",
    secret character varying(1000) COLLATE pg_catalog."default",
    private_key character varying(1000) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    created_at timestamp with time zone
	CONSTRAINT keycloak_infos_ukey UNIQUE (cluster_id, realm, secret)
);
