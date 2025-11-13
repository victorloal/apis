CREATE SCHEMA IF NOT EXISTS "public";


ALTER TABLE "public"."CrytografiaHomorfica" ADD CONSTRAINT "fk_Crytografia_Homorfica_id_Publico_HE" FOREIGN KEY("id") REFERENCES "public"."Publico"("HE");

CREATE TABLE "public"."AutoridadesElectorales" (
    "Id" uuid NOT NULL,
    "Name" varchar NOT NULL,
    "IdParty" uuid,
    "Mail" varchar NOT NULL,
    "password" varchar NOT NULL,
    "HE" uuid NOT NULL,
    CONSTRAINT "pk_AutoridadesElectorales_Id" PRIMARY KEY ("Id")
);

CREATE TABLE "public"."CrytografiaHomorfica" (
    "id" uuid NOT NULL,
    "T" integer NOT NULL,
    "N" integer NOT NULL,
    CONSTRAINT "pk_table_2_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."Publico" (
    "id" uuid NOT NULL,
    "PK" varchar,
    "HE" uuid,
    "Params" JSON,
    CONSTRAINT "pk_table_3_id" PRIMARY KEY ("id")
);

-- Foreign key constraints
-- Schema: public
ALTER TABLE "public"."Autoridades_Electorales" ADD CONSTRAINT "fk_Autoridades_Electorales_HE_Crytografia_Homorfica_id" FOREIGN KEY("HE") REFERENCES "public"."Crytografia Homorfica"("id");
ALTER TABLE "public"."Crytografia Homorfica" ADD CONSTRAINT "fk_Crytografia_Homorfica_id_Publico_HE" FOREIGN KEY("id") REFERENCES "public"."Publico"("HE");