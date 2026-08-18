-- Password authentication now uses password_hash. Some deployed databases
-- still have the pre-migration password column, while fresh databases do not.
DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'password'
  ) THEN
    ALTER TABLE "users" ALTER COLUMN "password" DROP NOT NULL;
  END IF;
END
$$;
