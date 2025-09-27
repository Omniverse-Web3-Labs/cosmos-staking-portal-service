drop table if exists "public"."staking_reward_daily";

CREATE TABLE "public"."staking_reward_daily" (
  "id" bigserial NOT NULL,
  "delegator" varchar(255) NOT NULL,
  "reward_amount" numeric(20, 10) NOT NULL,
  "staking_amount" numeric(20, 10) NOT NULL,
  "day_time" int8 NOT NULL,
  "created_at" int8 NOT NULL,
  CONSTRAINT "staking_reward_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "uidx_delegator_day" UNIQUE ("delegator", "day_time")
);
CREATE INDEX "staking_reward_daily_idx_day" ON "public"."staking_reward_daily" USING btree (
  "day_time"  ASC NULLS LAST
);

CREATE TABLE "public"."user" (
  "id" bigserial NOT NULL,
  "address" varchar(42) COLLATE "pg_catalog"."default" NOT NULL,
  "challenge" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6) NOT NULL,
  "updated_at" timestamptz(6) NOT NULL,
  "wallet_uuid" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "wallet_name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "announcement" int2 NOT NULL DEFAULT 0,
  CONSTRAINT "user_pk" PRIMARY KEY ("id"),
  CONSTRAINT "unique_address" UNIQUE ("address")
)
;

ALTER TABLE "public"."user" 
  OWNER TO "postgres";