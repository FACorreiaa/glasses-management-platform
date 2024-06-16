create table glasses
(
  glasses_id         uuid primary key     default uuid_generate_v4(),
  reference          text        not null,
  brand              text        not null,
  color              text        not null,
  left_eye_strength  numeric     not null,
  right_eye_strength numeric     not null,
  created_at         timestamptz not null default now(),
  updated_at         timestamptz
);

create trigger set_updated_at
  before update
  on glasses
  for each row
execute procedure set_updated_at();

insert into "glasses" (reference, brand, color, left_eye_strength, right_eye_strength)
values ('Glasses 1', 'Brand 1', 'red', 1.5, 1.5),
       ('Glasses 2', 'Brand 2', 'yellow', 2.5, 2.5),
       ('Glasses 3', 'Brand 3', 'orange', 3.5, 3.5),
       ('Glasses 4', 'Brand 4', 'orange', 4.5, 4.5),
       ('Glasses 5', 'Brand 5', 'blue', 5.5, 5.5),
       ('Glasses 6', 'Brand 6', 'black', 6.5, 6.5),
       ('Glasses 7', 'Brand 7', 'black', 7.5, 7.5),
       ('Glasses 8', 'Brand 8', 'black', 8.5, 8.5),
       ('Glasses 9', 'Brand 9', 'black', 9.5, 9.5),
       ('Glasses 10', 'Brand 10', 'black', 10.5, 10.5)
on conflict do nothing;
