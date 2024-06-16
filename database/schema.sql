CREATE TYPE "link" AS ENUM (
  'primary',
  'secondary'
);

CREATE TABLE "contact" (
  "id" serial PRIMARY KEY,
  "phone_number" varchar,
  "email" varchar,
  "linked_id" int,
  "link_precedence" link NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "deleted_at" timestamp
);

CREATE  FUNCTION update_time()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_on_change
    BEFORE UPDATE
    ON
       contact
    FOR EACH ROW
EXECUTE PROCEDURE update_time();