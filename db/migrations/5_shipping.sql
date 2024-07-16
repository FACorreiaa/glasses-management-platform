create table shipping_details
(
  shipping_id   uuid primary key     default uuid_generate_v4(),
  glasses_id    uuid        not null references glasses (glasses_id) on delete cascade,
  customer_id   uuid        not null references customer (customer_id) on delete cascade,
  shipped_by    uuid        not null references collaborator (user_id) on delete cascade,
  shipping_date timestamptz not null default now(),
  created_at    timestamptz not null default now(),
  updated_at    timestamptz
);

create trigger set_updated_at
  before update
  on shipping_details
  for each row
execute procedure set_updated_at();
