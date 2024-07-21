create table customer
(
  customer_id     uuid primary key     default uuid_generate_v4(),
  user_id         uuid        not null references "user" (user_id) on delete cascade,
  glasses_id      uuid        not null references glasses (glasses_id) on delete cascade,
  name            text        not null,
  card_id_number  text unique not null,
  address         text,
  address_details text,
  city            text,
  postal_code     text,
  country         text,
  continent       text,
  phone_number    text,
  email           citext unique,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz
);

create trigger set_updated_at
  before update
  on customer
  for each row
execute procedure set_updated_at();
