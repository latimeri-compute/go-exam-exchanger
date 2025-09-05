CREATE TABLE valutes(
    id integer GENERATED ALWAYS AS IDENTITY NOT NULL,
    code varchar(6) NOT NULL,
    full_name varchar(255) NOT NULL,
    PRIMARY KEY(id)
);
CREATE UNIQUE INDEX unique_code ON valutes USING btree (code);

INSERT INTO valutes (code, full_name) VALUES
('rub', 'Russian Ruble'),
('usd', 'United States Dollar'),
('eur', 'Euro')
ON CONFLICT DO NOTHING;

CREATE TABLE exchanges(
    from_valute_id integer NOT NULL,
    to_valute_id integer NOT NULL,
    rate bigint NOT NULL,
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    rate_id integer GENERATED ALWAYS AS IDENTITY NOT NULL,
    PRIMARY KEY(rate_id),
    CONSTRAINT exchange_from_valute_id_fkey FOREIGN key(from_valute_id) REFERENCES valutes(id),
    CONSTRAINT exchange_to_valute_id_fkey FOREIGN key(to_valute_id) REFERENCES valutes(id)
);
CREATE UNIQUE INDEX unique_combo ON exchanges USING btree (from_valute_id, to_valute_id);
INSERT INTO exchanges(from_valute_id, to_valute_id, rate) VALUES
(2, 1, 803466),
(3, 1, 935604),
(1, 2, 124),
(1, 3, 107),
(2, 3, 8588),
(3, 2, 11645)
ON CONFLICT DO NOTHING;