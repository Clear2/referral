DROP TABLE "credit_transactions";
DROP TABLE "referrals";
-- reverse: set comment to table: "users"
COMMENT ON TABLE "users" IS '';
-- reverse: create index "users_referral_code_key" to table: "users"
DROP INDEX "users_referral_code_key";
-- reverse: create index "users_email_key" to table: "users"
DROP INDEX "users_email_key";
-- reverse: create index "user_name_email" to table: "users"
DROP INDEX "user_name_email";
-- reverse: create "users" table
DROP TABLE "users";
