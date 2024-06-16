create table customer
(
  customer_id    uuid primary key     default uuid_generate_v4(),
  name           text,
  address        text,
  email          citext unique,
  phone_number   text,
  card_id_number text unique not null,
  created_at     timestamptz not null default now(),
  updated_at     timestamptz
);

create trigger set_updated_at
  before update
  on customer
  for each row
execute procedure set_updated_at();
