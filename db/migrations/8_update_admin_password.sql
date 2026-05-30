-- Update the admin account password.
-- 2_user.sql already ran on prod, so the seed insert cannot be edited
-- (immutable applied migration). This forward migration sets the new hash.
update "user"
set password_hash = '$2b$10$/99i3EjpM1lpIMqO9yMlaO5Zk8o9teshHaKHEA3lN.g4xwweBRGY'
where email = 'admin@email.com';
