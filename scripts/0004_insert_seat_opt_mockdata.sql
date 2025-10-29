-- ===========================================================
-- Seat Optimization Mock Data Generator
-- สร้าง mock data เยอะๆ สำหรับหน้า seat-optimization
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
-- 1) เพิ่ม Users เยอะๆ (inactive users สำหรับ revoke) - FIXED
-- ===========================================================
WITH gs AS (
  SELECT generate_series AS n FROM generate_series(1, 5)  -- 5 inactive users per department
)
INSERT INTO users (company_code, department_code, username, email, display_name, title, employee_id, status)
SELECT
  c.code AS company_code,
  d.code AS department_code,
  format('inactive_%s_%s_%s', c.code, d.code, gs.n) AS username,
  format('inactive.%s.%s.%s@%s.com',
         lower(c.code),
         lower(d.code),
         gs.n,
         lower(replace(c.name, ' ', ''))) AS email,
  format('Inactive User %s-%s-%s', c.code, d.code, gs.n) AS display_name,
  'Employee' AS title,
  format('EMP-%s-%s-%s', c.code, d.code, gs.n) AS employee_id,
  'inactive' AS status
FROM companies c
JOIN departments d ON d.company_id = c.id
CROSS JOIN gs
ON CONFLICT (email) DO NOTHING;  -- กันซ้ำรันสคริปต์หลายรอบ

-- ===========================================================
-- 2) เพิ่ม License Assignments สำหรับ inactive users
-- ===========================================================

-- Assign licenses แพงๆ ให้ inactive users (เพื่อเห็น savings เยอะๆ)
DROP TABLE IF EXISTS tmp_seat_opt_inactive;
CREATE TEMP TABLE tmp_seat_opt_inactive AS
SELECT 
  u.id AS user_id,
  u.company_code,
  u.department_code,
  i.id AS license_id,
  i.app_id,
  i.license_tier
FROM users u
JOIN license_inventories i ON i.company_code = u.company_code
WHERE u.status = 'inactive'
  AND i.license_tier IN ('Pro', 'Enterprise')
  AND random() < 0.6;  -- 60% chance per license

-- Insert inactive assignments (assigned 90+ days ago)
INSERT INTO license_assignments(
  company_code, user_id, app_id, license_id, license_tier, 
  assignment_source, assigned_at, reason
)
SELECT 
  company_code, user_id, app_id, license_id, license_tier,
  'manual',
  current_timestamp - ((90 + floor(random() * 180)) || ' days')::interval,  -- assigned 90-270 days ago
  'seat_opt_mock_inactive'
FROM tmp_seat_opt_inactive;

-- ===========================================================
-- 3) สร้าง Pending Requests (สำหรับ reallocate) - FIXED
-- ===========================================================

DROP TABLE IF EXISTS tmp_reallocate_opps;
CREATE TEMP TABLE tmp_reallocate_opps AS
WITH dept_pairs AS (
  SELECT 
    inactive.company_code,
    inactive.department_code AS from_department,
    d.code AS to_department,        -- อีกแผนกในบริษัทเดียวกัน
    inactive.app_id,
    inactive.license_tier,
    COUNT(*) AS inactive_count
  FROM tmp_seat_opt_inactive AS inactive
  JOIN companies c
    ON c.code = inactive.company_code
  JOIN departments d
    ON d.company_id = c.id
   AND d.code <> inactive.department_code   -- ต่าง department
  WHERE random() < 0.3                      -- 30% chance per pair
  GROUP BY inactive.company_code, inactive.department_code, d.code,
           inactive.app_id, inactive.license_tier
  HAVING COUNT(*) >= 2
)
SELECT * FROM dept_pairs;

-- สร้าง pending request records (simulate ด้วย recommendations) - FIXED
INSERT INTO recommendations(
  company_code, type, target_level, target_ref_id, app_id,
  action, impact_saving_amt, priority, reason_json, status
)
SELECT 
  tro.company_code,
  'seat_opt',
  'department',
  d2.id AS target_ref_id,                   -- department ต้นทาง (from_department) ของบริษัทเดียวกัน
  tro.app_id,
  'reallocate',
  (tro.inactive_count * 750),               -- ประหยัดต่อคน
  7,
  jsonb_build_object(
    'reason', 'reallocate_opportunity',
    'from_department', tro.from_department,
    'to_department', tro.to_department,
    'inactive_users', tro.inactive_count,
    'license_tier', tro.license_tier
  ),
  'pending'
FROM tmp_reallocate_opps AS tro
JOIN companies c2
  ON c2.code = tro.company_code
JOIN departments d2
  ON d2.company_id = c2.id
 AND d2.code = tro.from_department;

-- ===========================================================
-- 4) สร้าง Downgrade Opportunities (users ที่ใช้ tier แพงแต่ไม่จำเป็น)
-- ===========================================================

-- หา users ที่มี Pro/Enterprise licenses แต่ไม่ได้ใช้ features พิเศษ
DROP TABLE IF EXISTS tmp_downgrade_candidates;
CREATE TEMP TABLE tmp_downgrade_candidates AS
SELECT 
  la.company_code,
  la.user_id,
  la.app_id,
  la.license_tier,
  COUNT(ue.id) AS usage_count,
  CASE 
    WHEN la.license_tier = 'Enterprise' THEN 'Pro'
    WHEN la.license_tier = 'Pro' THEN 'Basic'
    ELSE 'Basic'
  END AS suggested_tier
FROM license_assignments la
LEFT JOIN usage_events ue ON ue.user_id = la.user_id 
  AND ue.app_id = la.app_id
  AND ue.event_at >= current_timestamp - interval '90 days'
WHERE la.license_tier IN ('Pro', 'Enterprise')
  AND la.revoked_at IS NULL
  AND random() < 0.4  -- 40% are candidates
GROUP BY la.company_code, la.user_id, la.app_id, la.license_tier
HAVING COUNT(ue.id) < 10;  -- ใช้ features น้อย (sign-in เท่านั้น ไม่ได้ใช้ features ลึก)

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
    WHEN dc.license_tier = 'Enterprise' THEN 1200  -- ประหยัด 1200/เดือน
    WHEN dc.license_tier = 'Pro' THEN 400          -- ประหยัด 400/เดือน
    ELSE 0
  END,
  6,
  jsonb_build_object(
    'reason', 'low_feature_usage',
    'current_tier', dc.license_tier,
    'suggested_tier', dc.suggested_tier,
    'usage_count_90d', dc.usage_count
  ),
  'draft'
FROM tmp_downgrade_candidates dc;

-- ===========================================================
-- 5) เพิ่ม Price Books แบบราคาสูง (เพื่อเห็น savings เยอะๆ)
-- ===========================================================

-- อัปเดต price books ให้มีราคาสูงพอเพื่อให้เห็น savings มาก
INSERT INTO price_books(app_id, tier, unit, list_price, currency, valid_from, valid_to)
SELECT 
  a.id,
  'Basic',
  'seat',
  800,  -- เพิ่มราคา Basic tier
  'THB',
  current_date - 150,
  current_date + 365
FROM apps a
WHERE a.category IN ('Productivity', 'Collaboration')
ON CONFLICT DO NOTHING;

INSERT INTO price_books(app_id, tier, unit, list_price, currency, valid_from, valid_to)
SELECT 
  a.id,
  'Pro',
  'seat',
  2000,  -- เพิ่มราคา Pro tier
  'THB',
  current_date - 150,
  current_date + 365
FROM apps a
WHERE a.category IN ('Productivity', 'Collaboration', 'DevOps')
ON CONFLICT DO NOTHING;

INSERT INTO price_books(app_id, tier, unit, list_price, currency, valid_from, valid_to)
SELECT 
  a.id,
  'Enterprise',
  'seat',
  3500,  -- เพิ่มราคา Enterprise tier
  'THB',
  current_date - 150,
  current_date + 365
FROM apps a
WHERE a.category IN ('Productivity', 'Collaboration', 'DevOps', 'Security')
ON CONFLICT DO NOTHING;

-- ===========================================================
-- 6) สร้าง Recommendations ทบท่วมใย (summary stats)
-- ===========================================================

-- สรุปสถานการณ์ให้เห็นภาพรวม
SELECT 
  '=== Seat Optimization Mock Data Summary ===' AS info,
  (SELECT COUNT(*) FROM users WHERE status = 'inactive') AS inactive_users_count,
  (SELECT COUNT(*) FROM tmp_seat_opt_inactive) AS inactive_assignments_count,
  (SELECT COUNT(*) FROM tmp_reallocate_opps) AS reallocate_opportunities_count,
  (SELECT COUNT(*) FROM tmp_downgrade_candidates) AS downgrade_opportunities_count,
  (SELECT COUNT(*) FROM recommendations WHERE type = 'seat_opt') AS total_recommendations_count,
  (SELECT SUM(impact_saving_amt) FROM recommendations WHERE type = 'seat_opt') AS total_potential_savings;
