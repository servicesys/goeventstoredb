CREATE SCHEMA eventstore;

CREATE TABLE eventstore.event
(
    event_id       uuid        NOT NULL,
    domain_tenant  varchar(40) NOT NULL,
    app_name       varchar(40) NOT NULL,
    transaction_id varchar(40) NOT NULL,
    event_type     varchar(40) NULL,
    event_version  varchar(40) NULL,
    payload        jsonb NULL,
    meta_data      jsonb NULL,
    created_at     timestamptz NULL,
    user_id        varchar(40) NULL,
    aggregate_id   bigint      NOT NULL,
    aggregate_type varchar(40) NULL
);

ALTER TABLE eventstore."event"
    ADD CONSTRAINT event_pk PRIMARY KEY (event_id);

CREATE
UNIQUE INDEX event_transaction_id_idx ON eventstore."event" (transaction_id);





