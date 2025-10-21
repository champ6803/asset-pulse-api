/* ==========================================================
   VIEW: recommend_pack_for_user
   Purpose: แสดง Software Pack ที่เหมาะกับแต่ละ user
   รวมทั้งจาก JD (recommendation_rules) และจาก template แผนก
   ========================================================== */

CREATE OR REPLACE VIEW recommend_pack_for_user AS
WITH user_jd AS (
  SELECT uj.user_id, jp.id AS job_profile_id, jp.code AS job_code, jp.name AS job_name
  FROM user_job_profiles uj
  JOIN job_profiles jp ON jp.id = uj.job_profile_id
  WHERE uj.assigned_at = (
    SELECT MAX(assigned_at) FROM user_job_profiles uj2 WHERE uj2.user_id = uj.user_id
  )
),
rule_pack AS (
  SELECT
    u.user_id,
    a.id AS app_id,
    a.key AS app_key,
    a.name AS app_name,
    rr.license_tier,
    rr.weight,
    rr.rationale,
    'rule_job_profile'::text AS source
  FROM recommendation_rules rr
  JOIN user_jd u ON rr.target='job_profile' AND rr.target_ref_id=u.job_profile_id
  JOIN apps a ON a.id = rr.app_id
),
dept_pack AS (
  SELECT
    u.id AS user_id,
    a.id AS app_id,
    a.key AS app_key,
    a.name AS app_name,
    pti.license_tier,
    NULL::numeric AS weight,
    'department_template'::text AS rationale,
    'template_department'::text AS source
  FROM users u
  JOIN departments d ON d.code = u.department_code
  JOIN purchase_templates pt ON pt.target='department' AND pt.scope_ref_id=d.id AND pt.is_active = true
  JOIN purchase_template_items pti ON pti.template_id = pt.id
  JOIN apps a ON a.id = pti.app_id
),
unioned AS (
  SELECT * FROM rule_pack
  UNION ALL
  SELECT * FROM dept_pack
),
rank_tier AS (
  SELECT
    user_id, app_id, app_key, app_name, source, rationale, weight, license_tier,
    CASE UPPER(license_tier)
      WHEN 'ENTERPRISE' THEN 5
      WHEN 'PRO' THEN 4
      WHEN 'STANDARD' THEN 3
      WHEN 'BASIC' THEN 2
      WHEN 'FREE' THEN 1
      ELSE 0
    END AS tier_rank
  FROM unioned
),
picked AS (
  SELECT DISTINCT ON (user_id, app_id)
    user_id, app_id, app_key, app_name, license_tier, tier_rank
  FROM rank_tier
  ORDER BY user_id, app_id, tier_rank DESC
)
SELECT
  p.user_id,
  u.display_name AS user_name,
  u.department_code,
  jd.job_code,
  jd.job_name,
  p.app_id,
  p.app_key,
  p.app_name,
  p.license_tier,
  (
    SELECT jsonb_agg(
      jsonb_build_object(
        'source', r.source,
        'tier', r.license_tier,
        'weight', r.weight,
        'rationale', r.rationale
      )
      ORDER BY r.tier_rank DESC NULLS LAST
    )
    FROM rank_tier r
    WHERE r.user_id = p.user_id AND r.app_id = p.app_id
  ) AS sources_json
FROM picked p
JOIN users u ON u.id = p.user_id
LEFT JOIN user_jd jd ON jd.user_id = p.user_id
ORDER BY p.user_id, p.app_name;
