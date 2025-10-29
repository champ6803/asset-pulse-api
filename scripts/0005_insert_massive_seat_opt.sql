-- ===========================================================
-- Massive Seat Optimization Mock Data
-- สร้าง recommendations เยอะๆ หลายแบบ (revoke, reallocate, downgrade)
-- เพื่อให้หน้า seat-optimization แสดงข้อมูลและการประหยัดเงินเยอะๆ
-- ===========================================================
SET search_path TO public;

-- ===========================================================
-- 0) ลบข้อมูลเก่า (รองรับการรันซ้ำ)
-- ===========================================================

-- ลบ recommendations เก่าที่เราสร้างไว้
DELETE FROM recommendations WHERE type = 'seat_opt';
DELETE FROM license_assignments WHERE reason IN ('seat_opt_mock_inactive', 'massive_seat_opt_revoke');
DELETE FROM users WHERE username LIKE 'inactive_%' OR username LIKE 'user_inactive_%';

-- ===========================================================
-- 1) สร้าง Active Users ที่ inactive (ไม่ได้ใช้ license 90+ วัน)
-- ===========================================================

-- สร้าง users ที่ active แต่มี licenses inactive
INSERT INTO users(company_code, department_code, username, email, display_name, title, employee_id, status)
SELECT 
  c.code,
  d.code,
  'user_inactive_' || c.code || '_' || d.code || '_' || generate_series,
  'user.' || c.code || '.' || d.code || '.' || generate_series || '@example.com',
  'User ' || generate_series,
  'Employee',
  'USR-' || c.code || '-' || generate_series,
  'active'
FROM companies c
CROSS JOIN departments d
CROSS JOIN generate_series(1, 8)  -- 8 users per department per company
WHERE d.company_id = c.id
  AND random() < 0.7;  -- เพิ่มให้ 70% ของ departments

-- ===========================================================
-- 2) Assign licenses แพงๆ ให้ inactive users (Revoke Opportunities)
-- ===========================================================

DROP TABLE IF EXISTS tmp_massive_revoke;
CREATE TEMP TABLE tmp_massive_revoke AS
SELECT 
  u.id AS user_id,
  u.company_code,
  u.department_code,
  i.id AS license_id,
  i.app_id,
  i.license_tier,
  a.name AS app_name,
  a.category AS app_category
FROM users u
JOIN license_inventories i ON i.company_code = u.company_code
JOIN apps a ON a.id = i.app_id
WHERE u.username LIKE 'user_inactive_%'
  AND i.license_tier IN ('Pro', 'Enterprise')
  AND random() < 0.7;  -- 70% of inactive users

-- Insert assignments (assigned 90-180 days ago, no usage)
INSERT INTO license_assignments(
  company_code, user_id, app_id, license_id, license_tier,
  assignment_source, assigned_at, reason
)
SELECT 
  company_code, user_id, app_id, license_id, license_tier,
  'manual',
  current_timestamp - ((90 + floor(random() * 90)) || ' days')::interval,
  'massive_seat_opt_revoke'
FROM tmp_massive_revoke;

-- ===========================================================
-- 3) Insert Revoke Recommendations แบบจำนวนเยอะ
-- ===========================================================

INSERT INTO recommendations(
  company_code, type, target_level, target_ref_id, app_id,
  action, impact_saving_amt, priority, reason_json, status
)
SELECT 
  m.company_code,
  'seat_opt',
  'user',
  m.user_id,
  m.app_id,
  'revoke',
  CASE 
    WHEN m.license_tier = 'Enterprise' THEN 3000 + floor(random() * 1000)  -- 3000-4000
    WHEN m.license_tier = 'Pro' THEN 1500 + floor(random() * 500)          -- 1500-2000
    ELSE 800 + floor(random() * 200)                                       -- 800-1000
  END,
  CASE 
    WHEN m.license_tier = 'Enterprise' THEN 8  -- High priority
    ELSE 6
  END,
  jsonb_build_object(
    'reason', 'inactive_user_90days',
    'app_name', m.app_name,
    'app_category', m.app_category,
    'license_tier', m.license_tier,
    'last_used_days', 90 + floor(random() * 90),
    'department_code', m.department_code
  ),
  'draft'
FROM tmp_massive_revoke m;

-- ===========================================================
-- 4) สร้าง Reallocate Opportunities (จำนวนมาก) - FIXED
-- ===========================================================

DROP TABLE IF EXISTS tmp_reallocate_massive;
CREATE TEMP TABLE tmp_reallocate_massive AS
WITH dept_usage AS (
  SELECT 
    u1.company_code,
    u1.department_code AS from_dept,
    d2.code           AS to_dept,            -- <-- ใช้ code แทน department_code
    m.app_id,
    m.license_tier,
    a.name            AS app_name,
    COUNT(DISTINCT u1.id)        AS inactive_count,
    floor(random() * 15 + 5)     AS pending_requests  -- simulate pending requests 5-20
  FROM tmp_massive_revoke m
  JOIN users u1 ON u1.id = m.user_id
  JOIN apps a  ON a.id = m.app_id
  JOIN companies c ON c.code = u1.company_code
  JOIN departments d2 
       ON d2.company_id = c.id               -- <-- map บริษัทให้ถูก
      AND d2.code <> u1.department_code      -- ข้ามแผนกเดิม
  GROUP BY u1.company_code, u1.department_code, d2.code, 
           m.app_id, m.license_tier, a.name
  HAVING COUNT(DISTINCT u1.id) >= 2
     AND random() < 0.4                      -- 40% chance
)
SELECT * FROM dept_usage;

-- Insert reallocate recommendations (ผูก target_ref_id ให้ถูกบริษัท)
INSERT INTO recommendations(
  company_code, type, target_level, target_ref_id, app_id,
  action, impact_saving_amt, priority, reason_json, status
)
SELECT 
  r.company_code,
  'seat_opt',
  'department',
  d_from.id AS target_ref_id,                 -- department ต้นทางของบริษัทเดียวกัน
  r.app_id,
  'reallocate',
  (r.inactive_count * 1500),                  -- ประหยัด 1500 ต่อคน
  7,
  jsonb_build_object(
    'reason', 'reallocate_pending_requests',
    'app_name', r.app_name,
    'license_tier', r.license_tier,
    'from_department', r.from_dept,
    'to_department', r.to_dept,
    'inactive_users', r.inactive_count,
    'pending_requests', r.pending_requests
  ),
  'pending'
FROM tmp_reallocate_massive r
JOIN companies c2 ON c2.code = r.company_code
JOIN departments d_from 
  ON d_from.company_id = c2.id 
 AND d_from.code = r.from_dept;               -- ระบุบริษัท + code ให้ชัด

-- ===========================================================
-- 5) สร้าง Downgrade Opportunities แบบเยอะ
-- ===========================================================

DROP TABLE IF EXISTS tmp_downgrade_massive;
CREATE TEMP TABLE tmp_downgrade_massive AS
SELECT 
  la.company_code,
  la.user_id,
  la.app_id,
  la.license_tier,
  a.name AS app_name,
  a.category AS app_category,
  u.department_code,
  COUNT(ue.id) AS usage_count,
  CASE 
    WHEN la.license_tier = 'Enterprise' THEN 'Pro'
    WHEN la.license_tier = 'Pro' THEN 'Basic'
    ELSE 'Basic'
  END AS suggested_tier
FROM license_assignments la
JOIN apps a ON a.id = la.app_id
JOIN users u ON u.id = la.user_id
LEFT JOIN usage_events ue ON ue.user_id = la.user_id 
  AND ue.app_id = la.app_id
  AND ue.event_at >= current_timestamp - interval '90 days'
WHERE la.license_tier IN ('Pro', 'Enterprise')
  AND la.revoked_at IS NULL
  AND u.status = 'active'
  AND u.username NOT LIKE 'inactive_%'
  AND random() < 0.5  -- 50% are downgrade candidates
GROUP BY la.company_code, la.user_id, la.app_id, la.license_tier, 
         a.name, a.category, u.department_code
HAVING COUNT(ue.id) < 15;  -- ใช้ features น้อย (mostly just sign-in)

-- Insert downgrade recommendations
INSERT INTO recommendations(
  company_code, type, target_level, target_ref_id, app_id,
  action, impact_saving_amt, priority, reason_json, status
)
SELECT 
  dc.company_code,
  'seat_opt',
  'user',
  dc.user_id,
  dc.app_id,
  'downgrade',
  CASE 
    WHEN dc.license_tier = 'Enterprise' THEN 2000 + floor(random() * 500)  -- 2000-2500
    WHEN dc.license_tier = 'Pro' THEN 800 + floor(random() * 200)          -- 800-1000
    ELSE 0
  END,
  6,
  jsonb_build_object(
    'reason', 'low_feature_usage',
    'app_name', dc.app_name,
    'app_category', dc.app_category,
    'current_tier', dc.license_tier,
    'suggested_tier', dc.suggested_tier,
    'usage_count_90d', dc.usage_count,
    'department_code', dc.department_code
  ),
  'draft'
FROM tmp_downgrade_massive dc;

-- ===========================================================
-- 6) Summary Statistics
-- ===========================================================

SELECT 
  '=== Massive Seat Optimization Summary ===' AS summary,
  (SELECT COUNT(*) FROM tmp_massive_revoke) AS revoke_opportunities,
  (SELECT COUNT(*) FROM tmp_reallocate_massive) AS reallocate_opportunities,
  (SELECT COUNT(*) FROM tmp_downgrade_massive) AS downgrade_opportunities,
  (SELECT COUNT(*) FROM recommendations WHERE type = 'seat_opt' AND action = 'revoke') AS total_revoke_recommendations,
  (SELECT COUNT(*) FROM recommendations WHERE type = 'seat_opt' AND action = 'reallocate') AS total_reallocate_recommendations,
  (SELECT COUNT(*) FROM recommendations WHERE type = 'seat_opt' AND action = 'downgrade') AS total_downgrade_recommendations,
  (SELECT SUM(impact_saving_amt) FROM recommendations WHERE type = 'seat_opt' AND action = 'revoke') AS total_revoke_savings,
  (SELECT SUM(impact_saving_amt) FROM recommendations WHERE type = 'seat_opt' AND action = 'reallocate') AS total_reallocate_savings,
  (SELECT SUM(impact_saving_amt) FROM recommendations WHERE type = 'seat_opt' AND action = 'downgrade') AS total_downgrade_savings,
  (SELECT SUM(impact_saving_amt) FROM recommendations WHERE type = 'seat_opt') AS total_potential_savings;
