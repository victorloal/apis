CREATE SCHEMA IF NOT EXISTS "public";

-- CURRENT_TIMESTAMP
CREATE TABLE "public"."election_authorities" (
    "id" serial NOT NULL,
    "election" uuid NOT NULL,
    "cc" int,
    "name" varchar(100) NOT NULL,
    "email" varchar(255) NOT NULL UNIQUE,
    "password" varchar(255) NOT NULL,
    -- secret_key_share
    "s_key" bytea NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    "updated_at" timestamp with time zone NOT NULL,
    "delete_at" timestamp with time zone,
    CONSTRAINT "pk_election_authorities_id" PRIMARY KEY ("id")
);
COMMENT ON TABLE "public"."election_authorities" IS 'CURRENT_TIMESTAMP';
COMMENT ON COLUMN "public"."election_authorities"."s_key" IS 'secret_key_share';

-- CURRENT_TIMESTAMP
CREATE TABLE "public"."voters" (
    "id" SERIAL,
    "elections" uuid,
    "token" varchar(200) NOT NULL,
    "vote_status" boolean NOT NULL DEFAULT false,
    "verification_hash" bytea,
    "is_active" boolean,
    "updated_at" timestamp with time zone NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    "delete_at" timestamp with time zone,
    CONSTRAINT "pk_voters_id" PRIMARY KEY ("id")
);
COMMENT ON TABLE "public"."voters" IS 'CURRENT_TIMESTAMP';

CREATE TABLE "public"."homomorphic_keys" (
    "id" serial NOT NULL,
    "elections" uuid NOT NULL,
    -- Homomorphic public key
    "p_key" bytea,
    -- encryption_parameters
    "Params" jsonb,
    "updated_at" timestamp with time zone NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    "delete_at" timestamp with time zone,
    CONSTRAINT "pk_homomorphic_keys_id" PRIMARY KEY ("id")
);
COMMENT ON COLUMN "public"."homomorphic_keys"."p_key" IS 'Homomorphic public key';
COMMENT ON COLUMN "public"."homomorphic_keys"."Params" IS 'encryption_parameters';

CREATE TABLE "public"."elections" (
    "id" uuid NOT NULL,
    "status" int,
    "encrypted" boolean NOT NULL,
    "name" varchar(100) NOT NULL,
    "description" text,
    "start_date" timestamp NOT NULL,
    "end_date" timestamp NOT NULL,
    "updated_at" timestamp with time zone NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    CONSTRAINT "pk_elections_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."status" (
    "id" int NOT NULL,
    "name" varchar(20) NOT NULL UNIQUE,
    "updated_at" timestamp with time zone NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    "delete_at" timestamp with time zone,
    CONSTRAINT "pk_status_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."candidates" (
    "id" bigint NOT NULL,
    "elections" uuid NOT NULL,
    "name" varchar(100) NOT NULL,
    "description" text,
    "photo_url" varchar(500) NOT NULL UNIQUE,
    "candidate_order" int NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    "update_at" timestamp with time zone NOT NULL,
    "delete_at" timestamp with time zone,
    CONSTRAINT "pk_candidates_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."ballots" (
    "id" varchar(200) NOT NULL,
    "elections" uuid NOT NULL,
    "voter" integer NOT NULL,
    "vote" jsonb NOT NULL,
    "voting_device_fingerprint" varchar(225),
    "ip_address" inet,
    "updated_at" timestamp with time zone NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    "delete_at" timestamp with time zone,
    CONSTRAINT "pk_table_7_id" PRIMARY KEY ("id", "elections", "voter")
);

CREATE TABLE "tally_results" (
    "id" serial NOT NULL,
    "election" uuid NOT NULL UNIQUE,
    -- e.g. {"candidate_id": count, ...} or encrypted aggregate
    "results" jsonb NOT NULL,
    "total_votes" int NOT NULL,
    "computed_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- authority id or system process
    "computed_by" varchar(255),
    -- optional cryptographic proofs / receipts
    "proof" jsonb,
    "updated_at" timestamp with time zone NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    CONSTRAINT "pk_tally_results_id" PRIMARY KEY ("id")
);
COMMENT ON COLUMN "tally_results"."results" IS 'e.g. {"candidate_id": count, ...} or encrypted aggregate';
COMMENT ON COLUMN "tally_results"."computed_by" IS 'authority id or system process';
COMMENT ON COLUMN "tally_results"."proof" IS 'optional cryptographic proofs / receipts';

CREATE TABLE "election_audit_config" (
    "id" serial NOT NULL,
    "election" uuid NOT NULL UNIQUE,
    "enable_ballot_audit" boolean NOT NULL DEFAULT true,
    "enable_access_logs" boolean NOT NULL DEFAULT true,
    "updated_at" timestamp with time zone NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    CONSTRAINT "pk_election_audit_config_id" PRIMARY KEY ("id")
);

-- Foreign key constraints
-- Schema: public
ALTER TABLE "public"."status" ADD CONSTRAINT "fk_status_id_elections_status" FOREIGN KEY("id") REFERENCES "public"."elections"("status");
ALTER TABLE "public"."election_authorities" ADD CONSTRAINT "fk_election_authorities_election_elections_id" FOREIGN KEY("election") REFERENCES "public"."elections"("id");
ALTER TABLE "public"."homomorphic_keys" ADD CONSTRAINT "fk_homomorphic_keys_elections_elections_id" FOREIGN KEY("elections") REFERENCES "public"."elections"("id");
ALTER TABLE "public"."candidates" ADD CONSTRAINT "fk_candidates_elections_elections_id" FOREIGN KEY("elections") REFERENCES "public"."elections"("id");
ALTER TABLE "public"."ballots" ADD CONSTRAINT "fk_ballots_voter_voters_id" FOREIGN KEY("voter") REFERENCES "public"."voters"("id");
ALTER TABLE "public"."voters" ADD CONSTRAINT "fk_voters_elections_elections_id" FOREIGN KEY("elections") REFERENCES "public"."elections"("id");
ALTER TABLE "public"."elections" ADD CONSTRAINT "fk_elections_id_ballots_elections" FOREIGN KEY("id") REFERENCES "public"."ballots"("elections");
ALTER TABLE "public"."elections" ADD CONSTRAINT "fk_elections_id_tally_results_election" FOREIGN KEY("id") REFERENCES "tally_results"("election");
ALTER TABLE "public"."elections" ADD CONSTRAINT "fk_elections_id_election_audit_config_election" FOREIGN KEY("id") REFERENCES "election_audit_config"("election");