-- Remove the type CHECK constraint from glasses.
-- Originally defined in 3_glasses.sql as CHECK (type IN ('adult','children')).
-- Applied migrations are immutable, so the relaxation ships as a new migration.
alter table glasses drop constraint if exists glasses_type_check;
