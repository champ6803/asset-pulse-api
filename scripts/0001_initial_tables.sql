-- ======================================================================
-- Asset Pulse – Core Schema (PostgreSQL)
-- VARCHAR(n) hybrid + no FK on company_code/department_code
-- ======================================================================

SET search_path TO public;

-- ======================
-- 1) Organization layer
-- ======================

CREATE TABLE orgs (
  id           BIGSERIAL PRIMARY KEY,
  code         VARCHAR(20) UNIQUE NOT NULL,
  name         VARCHAR(255) NOT NULL,
  description  TEXT,
  created_at   TIMESTAMPTZ DEFAULT now(),
  created_by   BIGINT,
  updated_at   TIMESTAMPTZ,
  updated_by   BIGINT
);

CREATE TABLE companies (
  id           BIGSERIAL PRIMARY KEY,
  org_id       BIGINT REFERENCES orgs(id) ON UPDATE CASCADE,
  code         VARCHAR(20) UNIQUE NOT NULL,         -- เช่น SCBX, DATAX
  name         VARCHAR(255) NOT NULL,
  description  TEXT,
  created_at   TIMESTAMPTZ DEFAULT now(),
  created_by   BIGINT,
  updated_at   TIMESTAMPTZ,
  updated_by   BIGINT
);

CREATE TABLE departments (
  id               BIGSERIAL PRIMARY KEY,
  company_id       BIGINT REFERENCES companies(id) ON UPDATE CASCADE,
  code             VARCHAR(20) UNIQUE NOT NULL,     -- dept code
  name             VARCHAR(255) NOT NULL,
  description      TEXT,
  manager_user_id  BIGINT,                          -- (optional) FK -> users.id หลังมี users แล้วค่อยเติม
  created_at       TIMESTAMPTZ DEFAULT now(),
  created_by       BIGINT,
  updated_at       TIMESTAMPTZ,
  updated_by       BIGINT
);

CREATE TABLE roles (
  id           BIGSERIAL PRIMARY KEY,
  company_code VARCHAR(20),                          -- no FK
  key          VARCHAR(50) UNIQUE NOT NULL,          -- เช่น Subsidiary.CTO
  name         VARCHAR(100) NOT NULL,
  description  TEXT,
  created_at   TIMESTAMPTZ DEFAULT now(),
  created_by   BIGINT,
  updated_at   TIMESTAMPTZ,
  updated_by   BIGINT
);

CREATE TABLE users (
  id               BIGSERIAL PRIMARY KEY,
  company_code     VARCHAR(20),                      -- no FK
  department_code  VARCHAR(20),                      -- no FK
  entra_id         UUID UNIQUE,
  email            VARCHAR(100) UNIQUE,
  display_name     VARCHAR(255),
  title            VARCHAR(100),
  employee_id      VARCHAR(50),
  status           VARCHAR(30),                      -- active/inactive/...
  created_at       TIMESTAMPTZ DEFAULT now(),
  created_by       BIGINT,
  updated_at       TIMESTAMPTZ,
  updated_by       BIGINT
);

CREATE TABLE user_roles (
  id            BIGSERIAL PRIMARY KEY,
  user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id       BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  scope_level   VARCHAR(30) NOT NULL,                -- group/subsidiary/department/app/user
  scope_ref_id  BIGINT,
  assigned_at   TIMESTAMPTZ DEFAULT now(),
  assigned_by   BIGINT,
  revoked_at    TIMESTAMPTZ,
  revoked_by    BIGINT,
  created_at    TIMESTAMPTZ DEFAULT now(),
  created_by    BIGINT,
  updated_at    TIMESTAMPTZ,
  updated_by    BIGINT
);
CREATE INDEX ix_user_roles_user_scope ON user_roles(user_id, scope_level, scope_ref_id);

-- ===========================
-- 2) Application & similarity
-- ===========================

CREATE TABLE apps (
  id                    BIGSERIAL PRIMARY KEY,
  company_code          VARCHAR(20),                  -- no FK
  key                   VARCHAR(50) UNIQUE,
  name                  VARCHAR(255) NOT NULL,
  alias                 VARCHAR(100),
  status                VARCHAR(30),                  -- Active/Inactive
  description           TEXT,
  category              VARCHAR(100),
  data_classification   VARCHAR(50),
  application_tier      VARCHAR(50),
  asvl                  VARCHAR(50),
  service_hour          VARCHAR(50),
  service_life_years    INT,
  owner_dept            VARCHAR(100),
  owner_business        VARCHAR(100),
  owner_it              VARCHAR(100),
  user_group            VARCHAR(100),
  lifecycle_status      VARCHAR(50),
  created_at            TIMESTAMPTZ DEFAULT now(),
  created_by            BIGINT,
  updated_at            TIMESTAMPTZ,
  updated_by            BIGINT
);

CREATE TABLE app_features (
  id            BIGSERIAL PRIMARY KEY,
  app_id        BIGINT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
  feature_key   VARCHAR(50),
  feature_name  VARCHAR(255),
  description   TEXT
);
CREATE INDEX ix_app_features_app ON app_features(app_id);

-- NOTE: ถ้าจะใช้ pgvector: CREATE EXTENSION IF NOT EXISTS vector; แล้วเปลี่ยน BYTEA -> VECTOR(1536) ตามโมเดล
CREATE TABLE app_feature_embeddings (
  id               BIGSERIAL PRIMARY KEY,
  app_id           BIGINT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
  embedding_vector BYTEA
);

CREATE TABLE app_similarity (
  app_id          BIGINT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
  similar_app_id  BIGINT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
  score           NUMERIC NOT NULL,
  method          VARCHAR(50),                        -- tfidf/embedding/...
  reason_json     JSONB,
  feature_names   TEXT,
  description     TEXT,
  PRIMARY KEY (app_id, similar_app_id)
);

-- =========================
-- 3) Vendor / Contracts etc
-- =========================

CREATE TABLE vendors (
  id          BIGSERIAL PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  tax_id      VARCHAR(50),
  country     VARCHAR(50),
  created_at  TIMESTAMPTZ DEFAULT now(),
  created_by  BIGINT,
  updated_at  TIMESTAMPTZ,
  updated_by  BIGINT
);

CREATE TABLE contracts (
  id                 BIGSERIAL PRIMARY KEY,
  company_code       VARCHAR(20),                      -- no FK
  vendor_id          BIGINT REFERENCES vendors(id) ON UPDATE CASCADE,
  contract_number    VARCHAR(50),
  title              VARCHAR(255),
  start_date         DATE,
  end_date           DATE,
  status             VARCHAR(30),                      -- active/expired/pending
  parent_contract_id BIGINT REFERENCES contracts(id),
  currency           VARCHAR(10),                      -- THB/USD/...
  created_at         TIMESTAMPTZ DEFAULT now(),
  created_by         BIGINT,
  updated_at         TIMESTAMPTZ,
  updated_by         BIGINT
);
CREATE INDEX ix_contracts_vendor  ON contracts(vendor_id);
CREATE INDEX ix_contracts_company ON contracts(company_code);

CREATE TABLE contract_terms (
  id                BIGSERIAL PRIMARY KEY,
  company_code      VARCHAR(20),                      -- no FK
  contract_id       BIGINT NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
  app_id            BIGINT REFERENCES apps(id) ON UPDATE CASCADE,
  license_tier      VARCHAR(50),                      -- Basic/Pro/Ent
  seat_committed    INT,
  price_per_seat    NUMERIC,
  unit_price        NUMERIC,
  uom               VARCHAR(20),                      -- Lic/User/Event...
  overage_rule_json JSONB
);
CREATE INDEX ix_contract_terms_contract_app ON contract_terms(contract_id, app_id);

CREATE TABLE contract_scopes (
  id           BIGSERIAL PRIMARY KEY,
  contract_id  BIGINT NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
  scope_level  VARCHAR(30) NOT NULL,                  -- group/subsidiary/department
  scope_ref_id BIGINT
);

CREATE TABLE group_consolidation_opps (
  id                   BIGSERIAL PRIMARY KEY,
  company_code         VARCHAR(20),                    -- no FK
  app_id               BIGINT REFERENCES apps(id) ON UPDATE CASCADE,
  current_contract_ids BIGINT[],
  potential_saving_amt NUMERIC,
  rationale            TEXT,
  snapshot_json        JSONB,
  status               VARCHAR(30) DEFAULT 'pending',  -- pending/accepted/rejected/applied
  created_at           TIMESTAMPTZ DEFAULT now(),
  created_by           BIGINT,
  updated_at           TIMESTAMPTZ,
  updated_by           BIGINT
);

-- ==========================
-- 4) License / Usage / Price
-- ==========================

CREATE TABLE license_inventories (
  id             BIGSERIAL PRIMARY KEY,
  company_code   VARCHAR(20),                           -- no FK
  app_id         BIGINT REFERENCES apps(id) ON UPDATE CASCADE,
  contract_id    BIGINT REFERENCES contracts(id) ON UPDATE CASCADE,
  license_tier   VARCHAR(50),
  total_seats    INT,
  reserved_seats INT DEFAULT 0,
  effective_date DATE,
  expire_date    DATE,
  created_at     TIMESTAMPTZ DEFAULT now(),
  created_by     BIGINT,
  updated_at     TIMESTAMPTZ,
  updated_by     BIGINT
);
CREATE INDEX ix_license_inv_app      ON license_inventories(app_id);
CREATE INDEX ix_license_inv_contract ON license_inventories(contract_id);

CREATE TABLE license_assignments (
  id                BIGSERIAL PRIMARY KEY,
  company_code      VARCHAR(20),                         -- no FK
  user_id           BIGINT REFERENCES users(id) ON UPDATE CASCADE,
  app_id            BIGINT REFERENCES apps(id) ON UPDATE CASCADE,
  license_tier      VARCHAR(50),
  assignment_source VARCHAR(30),                         -- manual/template/auto
  assigned_at       TIMESTAMPTZ DEFAULT now(),
  revoked_at        TIMESTAMPTZ,
  reason            TEXT,
  created_at        TIMESTAMPTZ DEFAULT now(),
  created_by        BIGINT,
  updated_at        TIMESTAMPTZ,
  updated_by        BIGINT
);
CREATE INDEX ix_license_assign_active ON license_assignments(app_id, user_id) WHERE revoked_at IS NULL;
CREATE INDEX ix_license_user          ON license_assignments(user_id, app_id);

CREATE TABLE usage_events (
  id            BIGSERIAL PRIMARY KEY,
  company_code  VARCHAR(20),                             -- no FK
  app_id        BIGINT REFERENCES apps(id) ON UPDATE CASCADE,
  user_id       BIGINT REFERENCES users(id) ON UPDATE CASCADE,
  event_at      TIMESTAMPTZ NOT NULL,
  source        VARCHAR(30),                             -- AAD/API/etc.
  event_type    VARCHAR(50),                             -- signin/edit/call...
  value_numeric NUMERIC,
  meta_json     JSONB,
  created_at    TIMESTAMPTZ DEFAULT now(),
  created_by    BIGINT,
  updated_at    TIMESTAMPTZ,
  updated_by    BIGINT
);
CREATE INDEX ix_usage_app_user_time ON usage_events(app_id, user_id, event_at);
CREATE INDEX ix_usage_company_time  ON usage_events(company_code, event_at);

CREATE TABLE price_books (
  id         BIGSERIAL PRIMARY KEY,
  app_id     BIGINT REFERENCES apps(id) ON UPDATE CASCADE,
  tier       VARCHAR(50),
  unit       VARCHAR(20),                                 -- seat/event/minute/gb
  list_price NUMERIC,
  currency   VARCHAR(10),
  valid_from DATE,
  valid_to   DATE,
  created_at TIMESTAMPTZ DEFAULT now(),
  created_by BIGINT,
  updated_at TIMESTAMPTZ,
  updated_by BIGINT
);

-- ==============================================
-- 5) Recommendation / Templates / Onboarding
-- ==============================================

CREATE TABLE job_profiles (
  id           BIGSERIAL PRIMARY KEY,
  company_code VARCHAR(20),                                -- no FK
  code         VARCHAR(50) UNIQUE NOT NULL,
  name         VARCHAR(255) NOT NULL,
  description  TEXT,
  created_at   TIMESTAMPTZ DEFAULT now(),
  created_by   BIGINT,
  updated_at   TIMESTAMPTZ,
  updated_by   BIGINT
);

CREATE TABLE user_job_profiles (
  id             BIGSERIAL PRIMARY KEY,
  user_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  job_profile_id BIGINT NOT NULL REFERENCES job_profiles(id) ON DELETE CASCADE,
  assigned_at    TIMESTAMPTZ DEFAULT now(),
  created_at     TIMESTAMPTZ DEFAULT now(),
  created_by     BIGINT,
  updated_at     TIMESTAMPTZ,
  updated_by     BIGINT,
  UNIQUE(user_id, job_profile_id)
);

CREATE TABLE recommendation_rules (
  id            BIGSERIAL PRIMARY KEY,
  company_code  VARCHAR(20),                               -- no FK
  target        VARCHAR(30) NOT NULL,                      -- job_profile/department/subsidiary
  target_ref_id BIGINT,
  app_id        BIGINT REFERENCES apps(id) ON UPDATE CASCADE,
  license_tier  VARCHAR(50),
  weight        NUMERIC,
  rationale     TEXT,
  created_at    TIMESTAMPTZ DEFAULT now(),
  created_by    BIGINT,
  updated_at    TIMESTAMPTZ,
  updated_by    BIGINT
);

CREATE TABLE recommendations (
  id                   BIGSERIAL PRIMARY KEY,
  company_code         VARCHAR(20),                        -- no FK
  type                 VARCHAR(30) NOT NULL,               -- seat_opt/payg_opt/onboarding/...
  target_level         VARCHAR(30) NOT NULL,               -- user/department/subsidiary/app
  target_ref_id        BIGINT,
  app_id               BIGINT REFERENCES apps(id) ON UPDATE CASCADE,
  action               VARCHAR(30),                        -- keep/revoke/switch/consolidate
  impact_saving_amt    NUMERIC,
  priority             INT DEFAULT 0,
  reason_json          JSONB,
  generated_at         TIMESTAMPTZ DEFAULT now(),
  status               VARCHAR(30) DEFAULT 'draft',
  is_accepted          BOOLEAN DEFAULT FALSE,
  applied_at           TIMESTAMPTZ,
  saving_verified_amt  NUMERIC,
  created_at           TIMESTAMPTZ DEFAULT now(),
  created_by           BIGINT,
  updated_at           TIMESTAMPTZ,
  updated_by           BIGINT
);
CREATE INDEX ix_reco_target ON recommendations(target_level, target_ref_id, status);

CREATE TABLE purchase_templates (
  id           BIGSERIAL PRIMARY KEY,
  company_code VARCHAR(20),                                -- no FK
  name         VARCHAR(255) NOT NULL,
  target       VARCHAR(30),                                -- department/onboarding
  scope_level  VARCHAR(30),
  scope_ref_id BIGINT,
  description  TEXT,
  version      INT DEFAULT 1,
  is_active    BOOLEAN DEFAULT TRUE,
  created_at   TIMESTAMPTZ DEFAULT now(),
  created_by   BIGINT,
  updated_at   TIMESTAMPTZ,
  updated_by   BIGINT
);

CREATE TABLE purchase_template_items (
  id           BIGSERIAL PRIMARY KEY,
  template_id  BIGINT NOT NULL REFERENCES purchase_templates(id) ON DELETE CASCADE,
  app_id       BIGINT REFERENCES apps(id) ON UPDATE CASCADE,
  license_tier VARCHAR(50),
  quantity     INT,
  editable     BOOLEAN DEFAULT TRUE,
  created_at   TIMESTAMPTZ DEFAULT now(),
  created_by   BIGINT,
  updated_at   TIMESTAMPTZ,
  updated_by   BIGINT
);

-- ===========================
-- 6) Workflow / Memo / Audit
-- ===========================

CREATE TABLE requests (
  id                BIGSERIAL PRIMARY KEY,
  company_code      VARCHAR(20),                           -- no FK
  type              VARCHAR(30) NOT NULL,                  -- purchase/memo/...
  requester_user_id BIGINT REFERENCES users(id),
  scope_level       VARCHAR(30),
  scope_ref_id      BIGINT,
  payload_json      JSONB,
  status            VARCHAR(30),                           -- draft/pending/approved/rejected
  created_at        TIMESTAMPTZ DEFAULT now(),
  created_by        BIGINT,
  updated_at        TIMESTAMPTZ,
  updated_by        BIGINT
);

CREATE TABLE request_steps (
  id               BIGSERIAL PRIMARY KEY,
  request_id       BIGINT NOT NULL REFERENCES requests(id) ON DELETE CASCADE,
  step_no          INT NOT NULL,
  approver_role    VARCHAR(100),
  approver_user_id BIGINT REFERENCES users(id),
  status           VARCHAR(30),                            -- pending/approved/rejected
  comment          TEXT,
  acted_at         TIMESTAMPTZ,
  sla_due_at       TIMESTAMPTZ,
  created_at       TIMESTAMPTZ DEFAULT now(),
  created_by       BIGINT,
  updated_at       TIMESTAMPTZ,
  updated_by       BIGINT,
  UNIQUE(request_id, step_no)
);

CREATE TABLE memos (
  id                        BIGSERIAL PRIMARY KEY,
  company_code              VARCHAR(20),                    -- no FK
  title                     VARCHAR(255) NOT NULL,
  body_md                   TEXT,
  linked_recommendation_ids BIGINT[],
  visibility_scope_level    VARCHAR(30),
  ref_id                    BIGINT,
  created_at                TIMESTAMPTZ DEFAULT now(),
  created_by                BIGINT,
  updated_at                TIMESTAMPTZ,
  updated_by                BIGINT
);

CREATE TABLE notifications (
  id          BIGSERIAL PRIMARY KEY,
  company_code VARCHAR(20),                                 -- no FK
  user_id     BIGINT REFERENCES users(id),
  type        VARCHAR(30),
  subject     VARCHAR(255),
  body        TEXT,
  ref_table   VARCHAR(100),
  ref_id      BIGINT,
  created_at  TIMESTAMPTZ DEFAULT now(),
  created_by  BIGINT,
  read_at     TIMESTAMPTZ
);

CREATE TABLE audit_logs (
  id            BIGSERIAL PRIMARY KEY,
  company_code  VARCHAR(20),                                -- no FK
  actor_user_id BIGINT REFERENCES users(id),
  action        VARCHAR(100) NOT NULL,
  entity        VARCHAR(100) NOT NULL,
  entity_id     BIGINT,
  diff_json     JSONB,
  at            TIMESTAMPTZ DEFAULT now()
);

-- ===========================
-- Helpful composite indexes
-- ===========================
CREATE INDEX ix_requests_scope   ON requests(scope_level, scope_ref_id, status);
CREATE INDEX ix_templates_scope  ON purchase_templates(scope_level, scope_ref_id) WHERE is_active;
