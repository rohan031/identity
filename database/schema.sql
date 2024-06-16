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
