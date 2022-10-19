CREATE
SCHEMA eventstore;

CREATE TABLE eventstore.event
(
    id             uuid         NOT NULL,
    domain_tenant  varchar(40)  NOT NULL,
    app_name       varchar(40)  NOT NULL,
    transaction_id varchar(40)  NOT NULL,
    event_type     varchar(120) NOT NULL,
    event_version  varchar(40)  NULL,
    event_data     jsonb        NULL,
    created_at     timestamptz  NULL,
    user_id        varchar(40)  NULL,
    aggregate_id   bigint       NOT NULL,
    aggregate_type varchar(40)  NULL
);


CREATE TABLE eventstore.event_type
(
    id         char(120)   NOT NULL,
    meta_data  jsonb       NOT NULL,
    created_at timestamptz NULL
);

ALTER TABLE eventstore."event"
    ADD CONSTRAINT event_pk PRIMARY KEY (id);


ALTER TABLE eventstore."event_type"
    ADD CONSTRAINT event_type_pk PRIMARY KEY (id);



ALTER TABLE eventstore."event"
    ADD CONSTRAINT event_type_fk FOREIGN KEY(event_type) REFERENCES eventstore."event_type"("id");


CREATE
    UNIQUE INDEX event_transaction_id_idx ON eventstore."event" (transaction_id);

-- Write your migrate up statements here

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
