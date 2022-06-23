CREATE TABLE IF NOT EXISTS records (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    surname text NOT NULL,
    middle_name text,
    iin text NOT NULL,
    dom int,
    kv int,
    city text,
    street text,
    cadastr_num text NOT NULL,
    area float DEFAULT 0,
    lat float,
    long float,
    comment text,
    status text NOT NULL DEFAULT 'ACTIVE_STATUS'
);
