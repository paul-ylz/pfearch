CREATE TABLE cat_foods
(
    id          serial PRIMARY KEY,
    ref         varchar(32),
    category    varchar(64) not null,
    name        text        not null,
    ingredients text        not null,
    language    varchar(32) not null default 'english',
    ts_idx_col  tsvector,
    created_at  timestamp            default current_timestamp,
    updated_at  timestamp            default current_timestamp
);


CREATE INDEX cat_food_ts_idx_col_idx ON cat_foods
    USING GIN (ts_idx_col);

CREATE INDEX cat_food_category_idx ON cat_foods (category);

