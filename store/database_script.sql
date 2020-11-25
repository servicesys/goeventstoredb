

CREATE SCHEMA eventstore;

CREATE TABLE eventstore.event (
	event_id uuid NOT NULL,
	domain_tenant varchar(40) NOT NULL,
	event_type varchar(40) NULL,
	event_version varchar(40) NULL,
	payload jsonb NULL,
	meta_data jsonb NULL,
	created_at timestamptz NULL,
	user_id varchar(40) NULL,
	aggregate_id int8 NOT NULL,
	aggregate_type varchar(40) NULL
);


