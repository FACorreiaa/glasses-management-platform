create table glasses
(
  glasses_id         uuid primary key     default uuid_generate_v4(),
  user_id            uuid        not null references "user" (user_id),
  reference          text        not null,
  brand              text,
  color              text,
  left_eye_strength  numeric     not null,
  right_eye_strength numeric     not null,
  type               text CHECK (
    type IN (
             'adult',
             'children'
      )
    ),
  is_in_stock        bool                 default true,
  features           text,
  created_at         timestamptz not null default now(),
  updated_at         timestamptz
);

create trigger set_updated_at
  before update
  on glasses
  for each row
execute procedure set_updated_at();
