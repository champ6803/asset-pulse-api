-- ===========================================================
-- Asset Pulse - FULL MOCK DATA for REAL SCBX subsidiaries (schema: public)
-- (FIXED to run end-to-end on PostgreSQL)
-- ===========================================================
SET search_path TO public;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
SELECT setseed(0.777);

--------------------------------------------------------------
-- 0) CLEANUP
--------------------------------------------------------------
DO $$
DECLARE t text;
BEGIN
  FOR t IN
    SELECT tablename FROM pg_tables
    WHERE schemaname='public'
  LOOP
    EXECUTE format('TRUNCATE TABLE %I RESTART IDENTITY CASCADE;', t);
  END LOOP;
END$$;

--------------------------------------------------------------
-- 1) ORG / COMPANIES / DEPARTMENTS / ROLES
--------------------------------------------------------------
INSERT INTO orgs(code, name, description)
VALUES ('SCBX','SCB X Group','Real subsidiaries mock for Asset Pulse');

WITH companies_input(code, name) AS (
  VALUES
   ('SCBX','SCB X'),
   ('SCB','The Siam Commercial Bank'),
   ('SCBAM','SCB Asset Management'),
   ('INVX','InnovestX'),
   ('PROTECT','SCB Protect'),
   ('ABACUS','Abacus Digital'),
   ('SCBJB','SCB Julius Baer'),
   ('SCB10X','SCB 10X'),
   ('MONIX','Monix'),
   ('TOKENX','Token X'),
   ('TECHX','SCB TechX'),
   ('ALPHAX','Alpha X'),
   ('RUTCHAM','Rutchayothin Asset Management'),
   ('SCBPLUS','SCB Plus'),
   ('MAHISORN','Mahisorn'),
   ('SCBTC','SCB Training Center'),
   ('AUTOX','AUTO X'),
   ('CARDX','Card X'),
   ('CARDXAMC','Card X AMC'),
   ('ALPHAXP','Alpha X Plus'),
   ('DATAX','SCB Data X'),
   ('AKULAKUX','Akulaku X'),
   ('POINTX','Point X'),
   ('CCB','Cambodian Commercial Bank'),
   ('SCBMM','Siam Commercial Bank Myanmar')
)
INSERT INTO companies(org_id, code, name, description)
SELECT 1, code, name, 'Real subsidiary'
FROM companies_input;

WITH c AS (SELECT id, code FROM companies),
dc AS (
  SELECT c.id AS company_id, (c.code||'_'||d) AS dept_code, d AS dept_name
  FROM c, (VALUES('FIN'),('HR'),('IT'),('OPS'),('MKT')) v(d)
)
INSERT INTO departments(company_id, code, name, description)
SELECT company_id, dept_code, dept_name, 'Dept '||dept_name FROM dc;

INSERT INTO roles(company_code, key, name, description)
VALUES
 (NULL,'Employee','Employee','Regular employee'),
 (NULL,'HR','HR','Human Resources'),
 (NULL,'Finance.Manager','Finance Manager','Finance approver'),
 (NULL,'Subsidiary.CTO','Subsidiary CTO','Subsidiary technology owner'),
 (NULL,'Group.CTO','Group CTO','Group technology owner'),
 (NULL,'App.Owner','App Owner','Application owner')
ON CONFLICT DO NOTHING;

--------------------------------------------------------------
-- 2) USERS + ROLE ASSIGNMENTS  (FIXED: ไม่ได้ชื่อซ้ำทั้งตับ)
--------------------------------------------------------------
WITH first_names(fn) AS (
  VALUES
  ('Somchai'),('Somsak'),('Somporn'),('Warinthra'),('Chanida'),
  ('Kittipong'),('Supakorn'),('Nattapong'),('Pimchanok'),('Jirawat'),
  ('Thanyarat'),('Chanon'),('Suthida'),('Pattarapon'),('Thidarat'),
  ('Noppadon'),('Tanyaporn'),('Chinnachod'),('Warin'),('Kornkanok'),
  ('Pimnapa'),('Akkarapong'),('Anan'),('Sirapob'),('Ratchanon'),
  ('Kamonchanok'),('Natcha'),('Phannapa'),('Nirutt'),('Metinee')
),
last_names(ln) AS (
  VALUES
  ('Tanyakul'),('Suksawat'),('Sriwong'),('Wongchai'),('Intarasiri'),
  ('Prasertsuk'),('Thongchai'),('Manop'),('Buranaporn'),('Khetkam'),
  ('Chantarangsu'),('Boonyarit'),('Aksorn'),('Siriporn'),('Jitkawin'),
  ('Chaiprasit'),('Khemsri'),('Kittikorn'),('Saengchan'),('Phanurat'),
  ('Rattanagul'),('Sae-Lim'),('Pumiporn'),('Jaroensri'),('Kantajit')
),
-- อ้างบริษัท/แผนก
c AS (SELECT id, code FROM companies),
d AS (SELECT id, code, company_id FROM departments),

-- สร้างแถวผู้ใช้ต่อบริษัท/แผนก
user_seed AS (
  SELECT c.id AS company_id, c.code AS company_code,
         d.id AS department_id, d.code AS department_code,
         gs AS seq
  FROM c
  JOIN d ON d.company_id = c.id
  JOIN generate_series(1,60) gs ON true
),

ranked AS (
  SELECT *, ROW_NUMBER() OVER (PARTITION BY company_id ORDER BY md5(company_id::text||department_id::text||seq::text)) AS rn
  FROM user_seed
),
picked AS (
  SELECT * FROM ranked WHERE rn <= (20 + floor(random()*31))
),

-- ทำลิสต์ชื่อ/นามสกุลแบบมี index
fn_src AS (SELECT row_number() OVER (ORDER BY fn) AS rn, fn FROM first_names),
ln_src AS (SELECT row_number() OVER (ORDER BY ln) AS rn, ln FROM last_names),
fn_cnt AS (SELECT COUNT(*)::int AS n FROM fn_src),
ln_cnt AS (SELECT COUNT(*)::int AS n FROM ln_src),

-- จับคู่ชื่อ-นามสกุลแบบ deterministic กระจายด้วย hash (ไม่ใช้ random() แล้ว)
named AS (
  SELECT
    p.*,
    1 + (abs( ('x'||substr(md5(p.company_code||p.department_code||p.seq::text),1,8))::bit(32)::int ) % (SELECT n FROM fn_cnt)) AS fn_rn,
    1 + (abs( ('x'||substr(md5(p.company_code||p.department_code||(p.seq+7)::text),1,8))::bit(32)::int ) % (SELECT n FROM ln_cnt)) AS ln_rn
  FROM picked p
),
named2 AS (
  SELECT p.*, f.fn, l.ln
  FROM named p
  JOIN fn_src f ON f.rn = p.fn_rn
  JOIN ln_src l ON l.rn = p.ln_rn
),

users_ins AS (
  INSERT INTO users(company_code, department_code, entra_id, email, display_name, title, employee_id, status)
  SELECT
    n.company_code,
    n.department_code,
    gen_random_uuid(),
    lower(
      concat(
        n.fn, '.', n.ln, '.', substr(md5(n.company_code||n.department_code||n.seq::text),1,6),
        '@', lower(n.company_code), '.com'
      )
    ) AS email,
    concat(n.fn,' ',n.ln) AS display_name,
    CASE
      WHEN n.department_code LIKE '%HR'  THEN 'HR Officer'
      WHEN n.department_code LIKE '%FIN' THEN 'Financial Analyst'
      WHEN n.department_code LIKE '%IT'  THEN 'Software Engineer'
      WHEN n.department_code LIKE '%OPS' THEN 'Operations Specialist'
      ELSE 'Marketing Executive'
    END AS title,
    concat(n.company_code,'-',n.seq) AS employee_id,
    'active'
  FROM named2 n
  RETURNING id, company_code, department_code, email
),

sub_groups AS (
  SELECT u.company_code, array_agg(u.id ORDER BY md5(u.email)) AS user_ids
  FROM users_ins u
  GROUP BY u.company_code
),
r AS (SELECT key, id AS role_id FROM roles)

INSERT INTO user_roles(user_id, role_id, scope_level, scope_ref_id, assigned_by)
SELECT
  u.id,
  CASE
    WHEN u.id = (SELECT (user_ids)[1] FROM sub_groups sg WHERE sg.company_code=u.company_code)
      THEN (SELECT role_id FROM r WHERE key='Subsidiary.CTO')
    WHEN u.id = (SELECT (user_ids)[2] FROM sub_groups sg WHERE sg.company_code=u.company_code)
      THEN (SELECT role_id FROM r WHERE key='Finance.Manager')
    WHEN u.id IN (
      (SELECT (user_ids)[3] FROM sub_groups sg WHERE sg.company_code=u.company_code),
      (SELECT (user_ids)[4] FROM sub_groups sg WHERE sg.company_code=u.company_code)
    ) THEN (SELECT role_id FROM r WHERE key='HR')
    ELSE (SELECT role_id FROM r WHERE key='Employee')
  END,
  'subsidiary', CAST(NULL AS BIGINT), 0
FROM users_ins u;

WITH scbx AS (
  SELECT array_agg(id ORDER BY md5(email)) AS uids
  FROM users WHERE company_code='SCBX'
)
INSERT INTO user_roles(user_id, role_id, scope_level, scope_ref_id, assigned_by)
SELECT (SELECT uids[1] FROM scbx), (SELECT id FROM roles WHERE key='Group.CTO'), 'group', CAST(NULL AS BIGINT), 0
UNION ALL
SELECT (SELECT uids[2] FROM scbx), (SELECT id FROM roles WHERE key='Group.CTO'), 'group', CAST(NULL AS BIGINT), 0;

/* ==========================================================
   3) APPS (realistic catalog) + FEATURES + PRICE BOOKS
   ========================================================== */

-- ใช้เป็น Global Catalog (company_code = NULL)
-- ถ้าในฐานะมึงมีข้อมูล apps เดิมอยู่แล้วและอยากแทนที่:
-- TRUNCATE apps CASCADE;  -- (ปลดคอมเมนต์ถ้าต้องการรีเซ็ต)

DROP TABLE IF EXISTS real_apps;
CREATE TEMP TABLE real_apps(
  key text,
  name text,
  alias text,
  category text,
  application_tier text
);

INSERT INTO real_apps(key, name, alias, category, application_tier) VALUES
-- Collaboration / Productivity
('JIRA',            'Atlassian Jira Software', 'Jira',          'DevOps',       'Core'),
('CONFLUENCE',      'Atlassian Confluence',    'Confluence',    'Collaboration','Supporting'),
('TRELLO',          'Trello',                   'Trello',        'Collaboration','Supporting'),
('MONDAY',          'monday.com Work OS',       'Monday',        'Collaboration','Supporting'),
('ASANA',           'Asana',                    'Asana',         'Collaboration','Supporting'),
('SLACK',           'Slack',                    'Slack',         'Collaboration','Core'),
('TEAMS',           'Microsoft Teams',          'Teams',         'Collaboration','Core'),
('O365',            'Microsoft 365',            'M365',          'Productivity', 'Core'),
('GWS',             'Google Workspace',         'GWS',           'Productivity', 'Core'),
('NOTION',          'Notion',                   'Notion',        'Collaboration','Supporting'),
('ZOOM',            'Zoom',                     'Zoom',          'Collaboration','Supporting'),
('MIRO',            'Miro',                     'Miro',          'Collaboration','Supporting'),
('FIGMA',           'Figma',                    'Figma',         'Design',       'Supporting'),

-- Dev / DevOps
('GITHUB',          'GitHub',                   'GitHub',        'DevOps',       'Core'),
('GITLAB',          'GitLab',                   'GitLab',        'DevOps',       'Core'),
('BITBUCKET',       'Bitbucket',                'Bitbucket',     'DevOps',       'Supporting'),
('AZDO',            'Azure DevOps',             'Azure DevOps',  'DevOps',       'Core'),
('JENKINS',         'Jenkins',                  'Jenkins',       'DevOps',       'Supporting'),
('CIRCLECI',        'CircleCI',                 'CircleCI',      'DevOps',       'Supporting'),
('POSTMAN',         'Postman',                  'Postman',       'DevOps',       'Supporting'),

-- Observability / Sec / ITSM
('DATADOG',         'Datadog',                  'Datadog',       'Security',     'Supporting'),
('NEWRELIC',        'New Relic',                'New Relic',     'Security',     'Supporting'),
('SENTRY',          'Sentry',                   'Sentry',        'Security',     'Supporting'),
('PAGERDUTY',       'PagerDuty',                'PagerDuty',     'Security',     'Supporting'),
('OKTA',            'Okta Workforce Identity',  'Okta',          'Security',     'Core'),
('CLOUDFLARE',      'Cloudflare',               'Cloudflare',    'Security',     'Supporting'),
('SERVICENOW',      'ServiceNow',               'ServiceNow',    'ITSM',         'Core'),
('ZENDESK',         'Zendesk',                  'Zendesk',       'ITSM',         'Supporting'),
('FRESHDESK',       'Freshdesk',                'Freshdesk',     'ITSM',         'Supporting'),

-- CRM / Analytics
('SALESFORCE',      'Salesforce Sales Cloud',   'Salesforce',    'CRM',          'Core'),
('HUBSPOT',         'HubSpot',                  'HubSpot',       'CRM',          'Supporting'),
('TABLEAU',         'Tableau',                  'Tableau',       'Analytics',    'Supporting'),
('POWERBI',         'Microsoft Power BI',       'Power BI',      'Analytics',    'Supporting');

-- เติมลงตาราง apps
WITH numbered AS (
  SELECT
    NULL::text AS company_code,
    key,
    name,
    alias,
    'Active'::text AS status,
    category,
    'External SaaS'::text AS data_classification,
    application_tier,
    ROW_NUMBER() OVER () AS rn
  FROM real_apps
)
INSERT INTO apps(company_code, key, name, alias, status, category, data_classification, application_tier, asvl, service_hour, service_life_years, lifecycle_status)
SELECT
  company_code,
  key,
  name,
  alias,
  status,
  category,
  data_classification,
  application_tier,
  'ASVL-' || (1000 + rn)::text,
  '24x7',
  CASE WHEN category IN ('DevOps','Security','ITSM') THEN 5 ELSE 3 END,
  'Production'
FROM numbered
ON CONFLICT DO NOTHING;

-- ===== FEATURES (realistic highlights) =====
DROP TABLE IF EXISTS app_feature_defs;
CREATE TEMP TABLE app_feature_defs(app_key text, feature_key text, feature_name text, description text);

INSERT INTO app_feature_defs(app_key, feature_key, feature_name, description) VALUES
('JIRA','jira.projects','Projects','Project & issue hierarchy'),
('JIRA','jira.boards','Agile Boards','Scrum/Kanban boards & sprints'),
('CONFLUENCE','confluence.pages','Pages','Team wiki & knowledge base'),
('CONFLUENCE','confluence.macros','Macros','Embeds/integrations'),
('MONDAY','monday.boards','Boards','Visual boards & workflows'),
('MONDAY','monday.automations','Automations','No-code automations'),
('ASANA','asana.tasks','Tasks','Tasks & timelines'),
('ASANA','asana.portfolio','Portfolio','Portfolio tracking'),
('SLACK','slack.channels','Channels','Team channels & DMs'),
('SLACK','slack.huddles','Huddles','Lightweight calls'),
('O365','m365.officeapps','Office Apps','Word/Excel/PowerPoint/Outlook'),
('O365','m365.onedrive','OneDrive','Cloud storage & sharing'),
('GWS','gws.workspace','Workspace Apps','Gmail/Docs/Sheets/Meet'),
('GWS','gws.drive','Drive','File storage & sharing'),
('GITHUB','github.repos','Repos','Git hosting & PRs'),
('GITHUB','github.actions','Actions','CI/CD pipelines'),
('GITLAB','gitlab.repos','Repos','Git hosting & merge requests'),
('GITLAB','gitlab.ci','CI','Pipelines & runners'),
('AZDO','azdo.boards','Boards','Boards/Backlogs/Sprints'),
('AZDO','azdo.pipelines','Pipelines','CI/CD pipelines'),
('POSTMAN','postman.collections','Collections','API collections & tests'),
('POSTMAN','postman.mock','Mock Servers','API mock & monitor'),
('DATADOG','datadog.apm','APM','Tracing & metrics'),
('DATADOG','datadog.log','Log Management','Centralized logs'),
('SENTRY','sentry.issues','Issues','Error tracking'),
('PAGERDUTY','pagerduty.oncall','On-call','On-call & escalation'),
('OKTA','okta.sso','SSO','SSO & MFA'),
('SERVICENOW','snow.incident','Incident','Incident & change mgmt'),
('ZENDESK','zendesk.support','Support','Ticketing & help center'),
('SALESFORCE','sf.leads','Leads','Lead/opportunity mgmt'),
('HUBSPOT','hubspot.crm','CRM','Marketing & sales hub'),
('TABLEAU','tableau.viz','Visualize','Interactive dashboards'),
('POWERBI','pbi.datasets','Datasets','Semantic models & reports');

-- ใส่ฟีเจอร์เข้าตาราง app_features
INSERT INTO app_features(app_id, feature_key, feature_name, description)
SELECT a.id, f.feature_key, f.feature_name, f.description
FROM app_feature_defs f
JOIN apps a ON a.key = f.app_key
ON CONFLICT DO NOTHING;

-- ===== PRICE BOOKS (mock แบบ realistic ตาม tier) =====
-- เพิ่ม tier หลัก ๆ: Free/Standard/Pro/Enterprise (หน่วย seat/เดือน, THB)
-- ปล. ถ้าตารางมี unique constraint ซ้ำ ให้พึ่ง ON CONFLICT DO NOTHING ตามนี้ได้เลย
INSERT INTO price_books(app_id, tier, unit, list_price, currency, valid_from, valid_to)
SELECT id, 'Free','seat/month',0,'THB',(current_date - 150),(current_date + 365) FROM apps
ON CONFLICT DO NOTHING;

INSERT INTO price_books(app_id, tier, unit, list_price, currency, valid_from, valid_to)
SELECT id, 'Standard','seat/month',
  CASE category
    WHEN 'DevOps' THEN 450
    WHEN 'Security' THEN 650
    WHEN 'ITSM' THEN 600
    WHEN 'CRM' THEN 700
    WHEN 'Analytics' THEN 500
    ELSE 400
  END,
  'THB',(current_date - 150),(current_date + 365) FROM apps
ON CONFLICT DO NOTHING;

INSERT INTO price_books(app_id, tier, unit, list_price, currency, valid_from, valid_to)
SELECT id, 'Pro','seat/month',
  CASE category
    WHEN 'DevOps' THEN 850
    WHEN 'Security' THEN 1200
    WHEN 'ITSM' THEN 1100
    WHEN 'CRM' THEN 1400
    WHEN 'Analytics' THEN 1000
    ELSE 800
  END,
  'THB',(current_date - 150),(current_date + 365) FROM apps
ON CONFLICT DO NOTHING;

INSERT INTO price_books(app_id, tier, unit, list_price, currency, valid_from, valid_to)
SELECT id, 'Enterprise','seat/month',
  CASE category
    WHEN 'DevOps' THEN 1400
    WHEN 'Security' THEN 2200
    WHEN 'ITSM' THEN 2000
    WHEN 'CRM' THEN 2600
    WHEN 'Analytics' THEN 1900
    ELSE 1600
  END,
  'THB',(current_date - 150),(current_date + 365) FROM apps
ON CONFLICT DO NOTHING;

/* ==========================================================
   4) VENDORS / CONTRACTS / TERMS / SCOPES (fixed)
   ========================================================== */

INSERT INTO vendors(name, tax_id, country) VALUES
 ('Microsoft','TAX-Msft','US'),
 ('Atlassian','TAX-Atlas','AU'),
 ('Google','TAX-Google','US'),
 ('OpenAI','TAX-OAI','US'),
 ('HashiCorp','TAX-Hashi','US')
ON CONFLICT DO NOTHING;

/* pick apps per company (~35%) */
DROP TABLE IF EXISTS tmp_pick_apps;
CREATE TEMP TABLE tmp_pick_apps AS
SELECT c.code AS company_code,
       (SELECT id FROM vendors ORDER BY random() LIMIT 1) AS vendor_id,
       a.id AS app_id
FROM companies c
CROSS JOIN apps a
WHERE random() < 0.35;

/* one or more contracts per company (bundle style) */
DROP TABLE IF EXISTS tmp_contracts;
CREATE TEMP TABLE tmp_contracts AS
SELECT DISTINCT
  p.company_code,
  p.vendor_id,
  concat(p.company_code,'-CN-',p.vendor_id,'-',floor(random()*1000))::varchar(50) AS contract_number,
  'Contract bundle for selected apps'::varchar(255) AS title,
  (current_date - (90 + floor(random()*120))::int) AS start_date,
  (current_date + (180 + floor(random()*365))::int) AS end_date,
  'active'::varchar(30) AS status,
  'THB'::varchar(10) AS currency
FROM tmp_pick_apps p;

-- insert contracts and keep (contract_id, company_code)
DROP TABLE IF EXISTS tmp_contract_ids2;
CREATE TEMP TABLE tmp_contract_ids2 AS
WITH ins AS (
  INSERT INTO contracts(company_code, vendor_id, contract_number, title, start_date, end_date, status, currency)
  SELECT c.company_code, c.vendor_id, c.contract_number, c.title, c.start_date, c.end_date, c.status, c.currency
  FROM tmp_contracts c
  RETURNING id, company_code
)
SELECT id AS contract_id, company_code
FROM ins;

-- create terms: map each (company, app) to a random contract of the same company
DROP TABLE IF EXISTS tmp_terms_src;
CREATE TEMP TABLE tmp_terms_src AS
SELECT p.company_code,
       p.app_id,
       (SELECT contract_id
        FROM tmp_contract_ids2 i
        WHERE i.company_code = p.company_code
        ORDER BY random() LIMIT 1) AS contract_id,
       CASE WHEN random()<0.3 THEN 'Enterprise'
            WHEN random()<0.7 THEN 'Pro'
            ELSE 'Basic' END::varchar(50) AS license_tier,
       (60 + floor(random()*180))::int AS seat_committed,
       (400 + floor(random()*600))::int AS price_per_seat
FROM tmp_pick_apps p;

-- insert terms and keep (term_id, contract_id, app_id)
DROP TABLE IF EXISTS tmp_term_ids2;
CREATE TEMP TABLE tmp_term_ids2 AS
WITH ins AS (
  INSERT INTO contract_terms(company_code, contract_id, app_id, license_tier, seat_committed, price_per_seat, unit_price, uom, overage_rule_json)
  SELECT t.company_code, t.contract_id, t.app_id, t.license_tier, t.seat_committed, t.price_per_seat, NULL, 'Lic',
         jsonb_build_object('payg', false)
  FROM tmp_terms_src t
  RETURNING id, contract_id, app_id
)
SELECT id AS term_id, contract_id, app_id
FROM ins;

-- scopes (subsidiary-level, no specific ref id)
INSERT INTO contract_scopes(contract_id, scope_level, scope_ref_id)
SELECT DISTINCT contract_id, 'subsidiary', CAST(NULL AS BIGINT)
FROM tmp_term_ids2;

/* ==========================================================
   5) LICENSE INVENTORIES / ASSIGNMENTS (fixed)
   ========================================================== */

DROP TABLE IF EXISTS tmp_inv_src;
CREATE TEMP TABLE tmp_inv_src AS
SELECT t.company_code, t.app_id, t.contract_id, t.license_tier,
       t.seat_committed AS total_seats,
       GREATEST(0,(t.seat_committed/12)::int) AS reserved_seats,
       (current_date - 45) AS effective_date,
       (current_date + 365) AS expire_date
FROM contract_terms t
JOIN contracts c ON c.id = t.contract_id;

-- insert inventories and keep (inv_id, company_code, app_id, total_seats)
DROP TABLE IF EXISTS tmp_inv_ids2;
CREATE TEMP TABLE tmp_inv_ids2 AS
WITH ins AS (
  INSERT INTO license_inventories(company_code, app_id, contract_id, license_tier, total_seats, reserved_seats, effective_date, expire_date)
  SELECT s.company_code, s.app_id, s.contract_id, s.license_tier, s.total_seats, s.reserved_seats, s.effective_date, s.expire_date
  FROM tmp_inv_src s
  RETURNING id, company_code, app_id, total_seats
)
SELECT id AS inv_id, company_code, app_id, total_seats
FROM ins;

-- users per company for assignment
DROP TABLE IF EXISTS tmp_users_by_comp;
CREATE TEMP TABLE tmp_users_by_comp AS
SELECT u.id AS user_id, u.company_code, u.department_code
FROM users u;

-- explode to assignments up to min(total_seats, 80)
DROP TABLE IF EXISTS tmp_assign_plan;
CREATE TEMP TABLE tmp_assign_plan AS
SELECT i.inv_id, i.company_code, i.app_id,
       u.user_id,
       ROW_NUMBER() OVER (PARTITION BY i.inv_id ORDER BY random()) AS rn,
       i.total_seats
FROM tmp_inv_ids2 i
JOIN tmp_users_by_comp u ON u.company_code = i.company_code;

INSERT INTO license_assignments(company_code, user_id, app_id, license_tier, assignment_source, assigned_at, reason)
SELECT
  p.company_code, p.user_id, p.app_id,
  (SELECT license_tier FROM license_inventories li WHERE li.id = p.inv_id),
  CASE WHEN random()<0.25 THEN 'template' ELSE 'manual' END,
  (current_timestamp - (floor(random()*75)||' days')::interval),
  'bulk assignment'
FROM tmp_assign_plan p
WHERE p.rn <= LEAST(p.total_seats, 80);

/* ==========================================================
   6) PRICE BOOKS
   ========================================================== */

INSERT INTO price_books(app_id, tier, unit, list_price, currency, valid_from, valid_to)
SELECT id, 'Basic','seat',300,'THB',(current_date - 150),(current_date + 365) FROM apps
ON CONFLICT DO NOTHING;

INSERT INTO price_books(app_id, tier, unit, list_price, currency, valid_from, valid_to)
SELECT id, 'Pro','seat',700,'THB',(current_date - 150),(current_date + 365) FROM apps
ON CONFLICT DO NOTHING;

/* ==========================================================
   7) USAGE EVENTS (6 months back, sparse)
   ========================================================== */

-- build day grid per assignment
DROP TABLE IF EXISTS tmp_usage_days;
CREATE TEMP TABLE tmp_usage_days AS
SELECT la.company_code, la.user_id, la.app_id,
       d::date AS day
FROM license_assignments la,
LATERAL generate_series(
  (current_date - interval '6 months')::date,
  current_date::date,
  interval '1 day'
) AS d
WHERE random() < 0.12;

INSERT INTO usage_events(company_code, app_id, user_id, event_at, source, event_type, value_numeric, meta_json)
SELECT
  company_code, app_id, user_id,
  (day + (floor(random()*24)||' hours')::interval + (floor(random()*60)||' minutes')::interval),
  CASE WHEN random()<0.82 THEN 'AAD' ELSE 'API' END,
  'signin',
  NULL,
  jsonb_build_object('ip', concat('10.',floor(random()*255)::int,'.',floor(random()*255)::int,'.',floor(random()*255)::int))
FROM tmp_usage_days;

/* ==========================================================
   8) JOB PROFILES / USER_JD / RECO RULES
   ========================================================== */

INSERT INTO job_profiles(company_code, code, name, description) VALUES
 (NULL,'ENG','Engineer','Software Engineering'),
 (NULL,'HR','HR','Human Resources'),
 (NULL,'FIN','Finance','Finance & Accounting'),
 (NULL,'MKT','Marketing','Marketing & Growth'),
 (NULL,'OPS','Operations','Operations')
ON CONFLICT DO NOTHING;

-- map users → JD by department suffix
INSERT INTO user_job_profiles(user_id, job_profile_id, assigned_at)
SELECT u.id,
  CASE
    WHEN u.department_code LIKE '%IT'  THEN (SELECT id FROM job_profiles WHERE code='ENG')
    WHEN u.department_code LIKE '%HR'  THEN (SELECT id FROM job_profiles WHERE code='HR')
    WHEN u.department_code LIKE '%FIN' THEN (SELECT id FROM job_profiles WHERE code='FIN')
    WHEN u.department_code LIKE '%OPS' THEN (SELECT id FROM job_profiles WHERE code='OPS')
    ELSE (SELECT id FROM job_profiles WHERE code='MKT')
  END,
  (current_timestamp - interval '12 days')
FROM users u
LEFT JOIN user_job_profiles uj ON uj.user_id = u.id
WHERE uj.user_id IS NULL;

-- reco rules by JD → subset of apps
DROP TABLE IF EXISTS tmp_rule_targets;
CREATE TEMP TABLE tmp_rule_targets AS
SELECT a.id AS app_id, a.category
FROM apps a;

-- ENG
INSERT INTO recommendation_rules(company_code, target, target_ref_id, app_id, license_tier, weight, rationale)
SELECT NULL,'job_profile',
  (SELECT id FROM job_profiles WHERE code='ENG'),
  t.app_id, 'Pro', 0.9, 'Dev tools for engineers'
FROM tmp_rule_targets t
WHERE t.category IN ('DevOps','Collaboration') AND random()<0.35;

-- HR
INSERT INTO recommendation_rules(company_code, target, target_ref_id, app_id, license_tier, weight, rationale)
SELECT NULL,'job_profile',
  (SELECT id FROM job_profiles WHERE code='HR'),
  t.app_id, 'Basic', 0.8, 'HR workflows'
FROM tmp_rule_targets t
WHERE t.category IN ('Collaboration','Productivity') AND random()<0.35;

-- FIN
INSERT INTO recommendation_rules(company_code, target, target_ref_id, app_id, license_tier, weight, rationale)
SELECT NULL,'job_profile',
  (SELECT id FROM job_profiles WHERE code='FIN'),
  t.app_id, 'Pro', 0.85, 'Finance reporting'
FROM tmp_rule_targets t
WHERE t.category IN ('Productivity','Security') AND random()<0.28;

/* ==========================================================
   9) RECOMMENDATIONS (seat optimization)
   ========================================================== */

DROP TABLE IF EXISTS tmp_last_seen;
CREATE TEMP TABLE tmp_last_seen AS
SELECT app_id, user_id, MAX(event_at) AS last_seen
FROM usage_events
WHERE event_at >= (current_timestamp - interval '6 months')
GROUP BY app_id, user_id;

DROP TABLE IF EXISTS tmp_inactive;
CREATE TEMP TABLE tmp_inactive AS
SELECT la.company_code, la.app_id, la.user_id,
       ls.last_seen
FROM license_assignments la
LEFT JOIN tmp_last_seen ls ON ls.app_id=la.app_id AND ls.user_id=la.user_id
WHERE (ls.last_seen IS NULL OR ls.last_seen < (current_timestamp - interval '30 days'));

INSERT INTO recommendations(company_code, type, target_level, target_ref_id, app_id, action, impact_saving_amt, priority, reason_json, status)
SELECT
  i.company_code,
  'seat_opt',
  'user',
  i.user_id,
  i.app_id,
  CASE WHEN random()<0.65 THEN 'revoke' ELSE 'downgrade' END,
  (250 + floor(random()*600))::int,
  5,
  jsonb_build_object('reason','inactive_30d'),
  'draft'
FROM tmp_inactive i
WHERE random()<0.18;

/* ==========================================================
   10) PURCHASE TEMPLATES / ITEMS / REQUESTS / STEPS (fixed)
   ========================================================== */

-- templates per department
DROP TABLE IF EXISTS tmp_templates;
CREATE TEMP TABLE tmp_templates AS
SELECT
  split_part(d.code,'_',1) AS company_code,
  'Template '||d.code AS name,
  'department'::varchar(30) AS target,
  'department'::varchar(30) AS scope_level,
  d.id::bigint AS scope_ref_id,
  'Baseline toolset for '||d.code AS description,
  1 AS version,
  true AS is_active
FROM departments d;

-- insert templates and collect ids
DROP TABLE IF EXISTS tmp_template_ids2;
CREATE TEMP TABLE tmp_template_ids2 AS
WITH ins AS (
  INSERT INTO purchase_templates(company_code, name, target, scope_level, scope_ref_id, description, version, is_active)
  SELECT t.company_code, t.name, t.target, t.scope_level, t.scope_ref_id, t.description, t.version, t.is_active
  FROM tmp_templates t
  RETURNING id, company_code, scope_ref_id
)
SELECT id AS template_id, company_code, scope_ref_id
FROM ins;

-- each template picks 3 random apps
INSERT INTO purchase_template_items(template_id, app_id, license_tier, quantity, editable)
SELECT tt.template_id,
       a.id,
       CASE WHEN random()<0.55 THEN 'Basic' ELSE 'Pro' END,
       5 + floor(random()*10)::int,
       true
FROM tmp_template_ids2 tt
JOIN LATERAL (SELECT id FROM apps ORDER BY random() LIMIT 3) a ON true;

-- make requests from a subset of templates (keep request_id, company_code)
DROP TABLE IF EXISTS tmp_req_ids2;
CREATE TEMP TABLE tmp_req_ids2 AS
WITH src AS (
  SELECT t.template_id, t.company_code, t.scope_ref_id
  FROM tmp_template_ids2 t
  WHERE random() < 0.22
),
ins AS (
  INSERT INTO requests(company_code, type, requester_user_id, scope_level, scope_ref_id, payload_json, status)
  SELECT
    r.company_code,
    'purchase',
    (SELECT id FROM users u WHERE u.company_code=r.company_code ORDER BY random() LIMIT 1),
    'department',
    r.scope_ref_id,
    jsonb_build_object('template_id', r.template_id),
    'pending'
  FROM src r
  RETURNING id, company_code
)
SELECT id AS request_id, company_code
FROM ins;

-- approval steps
INSERT INTO request_steps(request_id, step_no, approver_role, status, sla_due_at)
SELECT request_id, 1, 'Department.Manager', 'pending', (current_timestamp + interval '3 days')
FROM tmp_req_ids2;
INSERT INTO request_steps(request_id, step_no, approver_role, status, sla_due_at)
SELECT request_id, 2, 'Finance.Manager', 'pending', (current_timestamp + interval '5 days')
FROM tmp_req_ids2;

/* ==========================================================
   11) MEMOS / NOTIFICATIONS / AUDIT LOGS
   ========================================================== */

INSERT INTO memos(company_code, title, body_md, visibility_scope_level, ref_id)
VALUES ('SCBX','Seat Optimization This Month','Please review inactive seats across subsidiaries.','group', CAST(NULL AS BIGINT));

INSERT INTO notifications(company_code, user_id, type, subject, body, ref_table, ref_id)
SELECT 'SCBX',
       (SELECT id FROM users WHERE company_code='SCBX' ORDER BY random() LIMIT 1),
       'approval','Purchase Request','Please approve the purchase request','requests',
       (SELECT id FROM requests ORDER BY random() LIMIT 1);

INSERT INTO audit_logs(company_code, actor_user_id, action, entity, entity_id, diff_json)
SELECT 'SCBX',
       (SELECT id FROM users WHERE company_code='SCBX' ORDER BY random() LIMIT 1),
       'create_request','requests',(SELECT id FROM requests ORDER BY random() LIMIT 1),
       jsonb_build_object('status','pending','by','mock');

/* ==========================================================
   X) APP FEATURE EMBEDDINGS + APP SIMILARITY (REALISTIC MOCK)
   ========================================================== */

-- (ถ้ามี pgvector ในคลัสเตอร์อยู่แล้ว ใช้ได้เลย; ไม่มี ก็ยังรันส่วน similarity ได้)
-- CREATE EXTENSION IF NOT EXISTS vector;  -- ปลดคอมเมนต์ถ้ามีสิทธิ์ติดตั้ง

-- 1) TABLES (idempotent)
DO $DDL$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='app_feature_embeddings') THEN
    EXECUTE $SQL$
      CREATE TABLE app_feature_embeddings (
        id              BIGSERIAL PRIMARY KEY,
        app_id          BIGINT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
        feature_id      BIGINT NOT NULL REFERENCES app_features(id) ON DELETE CASCADE,
        method          VARCHAR(50) NOT NULL DEFAULT 'mock-v1',
        -- ถ้ามี pgvector ให้ใช้ vector(32); ถ้าไม่มีก็ใช้ float8[] ไปก่อน
        embedding_vector float8[] NOT NULL,
        created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
        created_by      BIGINT DEFAULT 0
      );
      CREATE UNIQUE INDEX app_feature_embeddings_uq
        ON app_feature_embeddings(feature_id, method);
    $SQL$;
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='app_similarity') THEN
    EXECUTE $SQL$
      CREATE TABLE app_similarity (
        id            BIGSERIAL PRIMARY KEY,
        app_a_id      BIGINT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
        app_b_id      BIGINT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
        method        VARCHAR(50) NOT NULL,
        score         NUMERIC(6,4) NOT NULL,        -- 0..1
        reason_json   JSONB,
        description   TEXT,
        created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
        created_by    BIGINT DEFAULT 0,
        CONSTRAINT app_similarity_pair_chk CHECK (app_a_id < app_b_id)
      );
      CREATE UNIQUE INDEX app_similarity_uq
        ON app_similarity(app_a_id, app_b_id, method);
      CREATE INDEX app_similarity_score_idx
        ON app_similarity(method, score DESC);
    $SQL$;
  END IF;
END
$DDL$;

-- 2) EMBEDDINGS (MOCK) —ผูกกับแอปฟีเจอร์จริงที่เธอเพิ่งใส่ (Jira, Monday, GitHub ฯลฯ)
--    ทำเวกเตอร์ 32 มิติแบบ deterministic-ish จาก hash ของ (app_id,feature_id)
--    เพื่อให้ rerun แล้วค่ายัง “ใกล้เคียงเดิม”
WITH feats AS (
  SELECT af.id AS feature_id, af.app_id,
         (af.feature_key || ':' || af.feature_name) AS sig
  FROM app_features af
),
expanded AS (
  -- แตก 32 มิติ สร้างตัวเลข 0..1 จาก md5 ที่ต่างกันในแต่ละมิติ
  SELECT
    f.app_id, f.feature_id,
    ARRAY(
      SELECT ((('x' || substr(md5(f.sig||'-'||gs::text),1,8))::bit(32))::int % 1000)::float8 / 1000.0
      FROM generate_series(1,32) gs
    ) AS vec
  FROM feats f
)
INSERT INTO app_feature_embeddings(app_id, feature_id, method, embedding_vector)
SELECT e.app_id, e.feature_id, 'mock-v1', e.vec
FROM expanded e
ON CONFLICT (feature_id, method) DO NOTHING;

-- 3) APP-TO-APP SIMILARITY (แนวทาง “ไม่พึ่งเวกเตอร์”) : Jaccard จากชุด feature_key
--    ข้อดี: รันได้ทุกที่, ให้ผล realistic (แอปสายเดียวกันคล้ายกัน)
DROP TABLE IF EXISTS tmp_app_feats;
CREATE TEMP TABLE tmp_app_feats AS
SELECT a.id AS app_id,
       array_agg(DISTINCT af.feature_key ORDER BY af.feature_key) AS fkeys
FROM apps a
LEFT JOIN app_features af ON af.app_id = a.id
GROUP BY a.id;

-- สร้างคู่ (a,b) แล้วคำนวณ Jaccard = |intersect| / |union|
WITH pairs AS (
  SELECT a.app_id AS a_id, b.app_id AS b_id,
         a.fkeys AS a_keys, b.fkeys AS b_keys
  FROM tmp_app_feats a
  JOIN tmp_app_feats b ON b.app_id > a.app_id
),
calc AS (
  SELECT
    a_id, b_id,
    (SELECT array_agg(x) FROM (
        SELECT UNNEST(a_keys) INTERSECT SELECT UNNEST(b_keys)
    ) s(x)) AS inter,
    (SELECT array_agg(x) FROM (
        SELECT UNNEST(a_keys) UNION SELECT UNNEST(b_keys)
    ) s(x)) AS uni
  FROM pairs
),
scored AS (
  SELECT
    a_id, b_id,
    CASE
      WHEN uni IS NULL OR array_length(uni,1)=0 THEN 0.0
      ELSE COALESCE(array_length(inter,1),0)::numeric / array_length(uni,1)
    END AS jaccard,
    jsonb_build_object(
      'method','feature-jaccard',
      'intersect', COALESCE(inter,'{}'::text[]),
      'union_count', COALESCE(array_length(uni,1),0)
    ) AS reason
  FROM calc
)
INSERT INTO app_similarity(app_a_id, app_b_id, method, score, reason_json, description)
SELECT a_id, b_id, 'feature-jaccard', ROUND(jaccard,4),
       reason,
       'Similarity based on feature-key overlap'
FROM scored
WHERE jaccard > 0   -- เก็บเฉพาะที่มีความคล้าย > 0
ON CONFLICT (app_a_id, app_b_id, method) DO UPDATE
  SET score = EXCLUDED.score,
      reason_json = EXCLUDED.reason_json,
      description = EXCLUDED.description;

-- (ทางเลือก) ถ้ามี pgvector และอยากคำนวณ cosine จาก “centroid per app”
-- ให้ใช้บล็อกนี้แทน/เพิ่มเติมได้
-- /*
-- DROP TABLE IF EXISTS tmp_app_centroid;
-- CREATE TEMP TABLE tmp_app_centroid AS
-- SELECT afe.app_id,
--        -- หาค่าเฉลี่ย 32 มิติจาก embedding_vector (float8[])
--        (SELECT ARRAY(
--           SELECT AVG(val)::float8
--           FROM (
--             SELECT i, (vals)[i] AS val
--             FROM (
--               SELECT generate_series(1, array_length(vals,1)) AS i, vals
--             ) t
--           ) z GROUP BY i ORDER BY i
--         )
--        )
--        AS centroid
-- FROM (
--   SELECT app_id, ARRAY_AGG(embedding_vector) AS all_vecs
--   FROM app_feature_embeddings
--   WHERE method='mock-v1'
--   GROUP BY app_id
-- ) s
-- CROSS JOIN LATERAL (
--   -- รวมเป็นเมทริกซ์แล้วบวกลดเป็นอาเรย์เดียว (ง่าย ๆ)
--   SELECT (SELECT ARRAY(
--            SELECT AVG(x)::float8
--            FROM UNNEST(ALL_VEC) AS x
--            )  -- (หมายเหตุ: โค้ดเฉลี่ยอาเรย์หลายแถวอาจต้องฟังก์ชันเสริมในบางเวอร์ชัน)
-- ) dummy
-- ;
--
-- -- cosine(a,b) = dot(a,b) / (||a||*||b||)
-- WITH pairs AS (
--   SELECT a.app_id AS a_id, b.app_id AS b_id, a.centroid AS a_vec, b.centroid AS b_vec
--   FROM tmp_app_centroid a
--   JOIN tmp_app_centroid b ON b.app_id > a.app_id
-- ),
-- cos AS (
--   SELECT a_id, b_id,
--     (
--       SELECT SUM(ax*bx) FROM UNNEST(a_vec, b_vec) AS t(ax,bx)
--     ) /
--     GREATEST(1e-9::float8,
--       sqrt((SELECT SUM(ax*ax) FROM UNNEST(a_vec) AS t(ax))) *
--       sqrt((SELECT SUM(bx*bx) FROM UNNEST(b_vec) AS t(bx)))
--     ) AS cosine
--   FROM pairs
-- )
-- INSERT INTO app_similarity(app_a_id, app_b_id, method, score, reason_json, description)
-- SELECT a_id, b_id, 'embedding-cosine', ROUND(cosine::numeric,4),
--        jsonb_build_object('method','embedding-cosine'),
--        'Cosine similarity of app centroids (mock embeddings)'
-- FROM cos
-- ON CONFLICT (app_a_id, app_b_id, method) DO UPDATE
--   SET score = EXCLUDED.score,
--       reason_json = EXCLUDED.reason_json,
--       description = EXCLUDED.description;
-- */
