CREATE TYPE "role" AS ENUM (
    'super',
    'owner',
    'employee',
    'customer'
    );

CREATE TABLE "users" (
                         "id" SERIAL PRIMARY KEY,
                         "merchant_id" int,
                         "outlets" text,
                         "name" varchar(100) NOT NULL,
                         "email" varchar(100) UNIQUE NOT NULL,
                         "password" varchar(100) NOT NULL,
                         "created_at" bigint,
                         "updated_at" bigint,
                         "role" role NOT NULL
);

CREATE TABLE "merchant" (
                            "id" SERIAL PRIMARY KEY,
                            "merchant_name" varchar(255) NOT NULL,
                            "created_at" bigint,
                            "updated_at" bigint
);

CREATE TABLE "outlets" (
                           "id" int PRIMARY KEY,
                           "merchant_id" int NOT NULL,
                           "outlet_name" varchar(255) NOT NULL,
                           "address" text NOT NULL,
                           "created_at" bigint,
                           "updated_at" bigint
);

CREATE TABLE "products" (
                            "id" int PRIMARY KEY,
                            "merchant_id" int,
                            "code" varchar(100) NOT NULL,
                            "name" varchar(255) NOT NULL,
                            "def_buy_price" int NOT NULL,
                            "def_sell_price" int NOT NULL,
                            "image" text NOT NULL DEFAULT '',
                            "created_at" bigint,
                            "updated_at" bigint
);

CREATE TABLE "product_price" (
                                 "id" varchar(100) PRIMARY KEY,
                                 "product_id" int NOT NULL,
                                 "outlet_id" int NOT NULL,
                                 "buy_price" int NOT NULL,
                                 "sell_price" int NOT NULL,
                                 "image" text NOT NULL DEFAULT '',
                                 "updated_at" bigint
);

ALTER TABLE "users" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchant" ("id");

ALTER TABLE "outlets" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchant" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchant" ("id");

ALTER TABLE "product_price" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_price" ADD FOREIGN KEY ("outlet_id") REFERENCES "outlets" ("id");

CREATE INDEX "u_product_id" ON "users" ("merchant_id");

CREATE INDEX "o_product_id" ON "outlets" ("merchant_id");

CREATE INDEX "p_product_id" ON "products" ("merchant_id");

CREATE INDEX "pp_product_id" ON "product_price" ("product_id");

CREATE INDEX "pp_outlet_id" ON "product_price" ("outlet_id");

COMMENT ON COLUMN "users"."outlets" IS 'slice of outlets restriction';
