
CREATE TABLE "user" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "name" VARCHAR(36) NOT NULL,
  "email" VARCHAR(100) NOT NULL,
  "password" VARCHAR(1024) NOT NULL,
  "created_at" TIMESTAMP DEFAULT now()
)

