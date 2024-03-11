CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "change_password_at" timestamptz DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_curency_key" UNIQUE ("owner","curency");
