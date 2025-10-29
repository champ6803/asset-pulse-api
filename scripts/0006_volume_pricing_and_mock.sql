-- ===========================================================
-- Volume Pricing DDL + Mock Data
-- - app_volume_pricing_tiers: rule per app (priority)
-- - vendor_volume_pricing_tiers: rule per vendor+feature_cluster (fallback)
-- - switching_cost_policy: policy per feature_cluster
-- - mock inserts seeded from existing apps, vendors, price_books
-- ===========================================================
SET search_path TO public;

-- =====================
-- 1) Create DDL (idempotent)
-- =====================

DO $DDL$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='app_volume_pricing_tiers') THEN
    EXECUTE $SQL$
      CREATE TABLE app_volume_pricing_tiers (
        id              BIGSERIAL PRIMARY KEY,
        app_id          BIGINT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
        threshold_qty   INT NOT NULL,
        unit_price      NUMERIC NOT NULL,
        currency        VARCHAR(10) NOT NULL,
        billing_period  VARCHAR(20) NOT NULL,          -- monthly/yearly
        pricing_mode    VARCHAR(20),                   -- piecewise/progressive (nullable)
        is_active       BOOLEAN DEFAULT TRUE,
        effective_from  DATE,
        effective_to    DATE,
        created_at      TIMESTAMPTZ DEFAULT now(),
        created_by      BIGINT,
        updated_at      TIMESTAMPTZ,
        updated_by      BIGINT
      );
      CREATE INDEX IF NOT EXISTS ix_app_vol_tier_app ON app_volume_pricing_tiers(app_id, is_active);
    $SQL$;
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='vendor_volume_pricing_tiers') THEN
    EXECUTE $SQL$
      CREATE TABLE vendor_volume_pricing_tiers (
        id                 BIGSERIAL PRIMARY KEY,
        vendor_id          BIGINT NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
        feature_cluster_key VARCHAR(100),              -- ใช้ apps.category เป็น proxy
        threshold_qty      INT NOT NULL,
        unit_price         NUMERIC NOT NULL,
        currency           VARCHAR(10) NOT NULL,
        billing_period     VARCHAR(20) NOT NULL,
        pricing_mode       VARCHAR(20),               -- piecewise/progressive (nullable)
        is_active          BOOLEAN DEFAULT TRUE,
        effective_from     DATE,
        effective_to       DATE,
        created_at         TIMESTAMPTZ DEFAULT now(),
        created_by         BIGINT,
        updated_at         TIMESTAMPTZ,
        updated_by         BIGINT
      );
      CREATE INDEX IF NOT EXISTS ix_vendor_vol_tier ON vendor_volume_pricing_tiers(vendor_id, feature_cluster_key, is_active);
    $SQL$;
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname='switching_cost_policy') THEN
    EXECUTE $SQL$
      CREATE TABLE switching_cost_policy (
        id                           BIGSERIAL PRIMARY KEY,
        feature_cluster_key          VARCHAR(100) NOT NULL,
        training_cost_per_user       NUMERIC DEFAULT 0,
        migration_flat_cost          NUMERIC DEFAULT 0,
        early_termination_penalty_rate NUMERIC DEFAULT 0,  -- 0..1
        created_at                   TIMESTAMPTZ DEFAULT now(),
        created_by                   BIGINT,
        updated_at                   TIMESTAMPTZ,
        updated_by                   BIGINT,
        UNIQUE(feature_cluster_key)
      );
    $SQL$;
  END IF;
END
$DDL$;

-- Unique constraint equivalent via expression index (allow NULL effective_from)
CREATE UNIQUE INDEX IF NOT EXISTS uq_app_vol_tier
  ON app_volume_pricing_tiers(
    app_id,
    threshold_qty,
    currency,
    billing_period,
    (COALESCE(effective_from, DATE '1900-01-01'))
  );

-- =====================
-- 2) Seed Mock Volume Pricing (App-level)
-- =====================

-- Strategy:
-- - For each app, derive a base monthly price from price_books (prefer 'Pro' > 'Standard' > any latest)
-- - Create tiers at thresholds 1, 50, 200 with 0%, 20%, 40% discounts respectively (piecewise)

WITH base_price AS (
  SELECT a.id AS app_id,
         -- prefer Pro seat/month, else Standard seat/month, else any recent list_price
         COALESCE(
           (
             SELECT pb.list_price::numeric
             FROM price_books pb
             WHERE pb.app_id=a.id AND lower(pb.tier)='pro' AND pb.unit LIKE 'seat%'
             ORDER BY pb.valid_to DESC NULLS LAST, pb.valid_from DESC NULLS LAST
             LIMIT 1
           ),
           (
             SELECT pb.list_price::numeric
             FROM price_books pb
             WHERE pb.app_id=a.id AND lower(pb.tier)='standard' AND pb.unit LIKE 'seat%'
             ORDER BY pb.valid_to DESC NULLS LAST, pb.valid_from DESC NULLS LAST
             LIMIT 1
           ),
           (
             SELECT pb.list_price::numeric
             FROM price_books pb
             WHERE pb.app_id=a.id AND pb.unit LIKE 'seat%'
             ORDER BY pb.valid_to DESC NULLS LAST, pb.valid_from DESC NULLS LAST
             LIMIT 1
           )
         ) AS base_list_thb
  FROM apps a
)
INSERT INTO app_volume_pricing_tiers(app_id, threshold_qty, unit_price, currency, billing_period, pricing_mode, is_active, effective_from)
SELECT app_id, threshold_qty,
       GREATEST(0.0::numeric, ROUND((base_list_thb * (1::numeric - discount_rate::numeric))::numeric, 2)) AS unit_price,
       'THB'::varchar(10) AS currency,
       'monthly'::varchar(20) AS billing_period,
       'piecewise'::varchar(20) AS pricing_mode,
       TRUE AS is_active,
       current_date - 30 AS effective_from
FROM base_price bp
CROSS JOIN (
  VALUES (1, 0.00::numeric), (50, 0.20::numeric), (200, 0.40::numeric)
) t(threshold_qty, discount_rate)
WHERE bp.base_list_thb IS NOT NULL
ON CONFLICT DO NOTHING;

-- =====================
-- 3) Seed Mock Volume Pricing (Vendor-level fallback)
-- =====================

-- Strategy:
-- - For each vendor x category (as feature_cluster_key), create tiers with slightly less aggressive discounts
-- - Use median of app base price in that category as reference

WITH app_cat AS (
  SELECT a.id AS app_id, a.category AS feature_cluster_key
  FROM apps a
),
app_base AS (
  SELECT ac.feature_cluster_key, COALESCE(
    (
      SELECT pb.list_price::numeric
      FROM price_books pb
      WHERE pb.app_id=ac.app_id AND lower(pb.tier)='pro' AND pb.unit LIKE 'seat%'
      ORDER BY pb.valid_to DESC NULLS LAST, pb.valid_from DESC NULLS LAST
      LIMIT 1
    ),
    (
      SELECT pb.list_price::numeric
      FROM price_books pb
      WHERE pb.app_id=ac.app_id AND lower(pb.tier)='standard' AND pb.unit LIKE 'seat%'
      ORDER BY pb.valid_to DESC NULLS LAST, pb.valid_from DESC NULLS LAST
      LIMIT 1
    ),
    (
      SELECT pb.list_price::numeric
      FROM price_books pb
      WHERE pb.app_id=ac.app_id AND pb.unit LIKE 'seat%'
      ORDER BY pb.valid_to DESC NULLS LAST, pb.valid_from DESC NULLS LAST
      LIMIT 1
    )
  ) AS list_price
  FROM app_cat ac
),
cat_ref AS (
  SELECT feature_cluster_key,
         PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY list_price) AS median_price
  FROM app_base
  WHERE list_price IS NOT NULL
  GROUP BY feature_cluster_key
),
vendor_x_cat AS (
  SELECT v.id AS vendor_id, c.feature_cluster_key, c.median_price
  FROM vendors v
  CROSS JOIN cat_ref c
)
INSERT INTO vendor_volume_pricing_tiers(
  vendor_id, feature_cluster_key, threshold_qty, unit_price, currency, billing_period, pricing_mode, is_active, effective_from
)
SELECT vendor_id, feature_cluster_key, threshold_qty,
       GREATEST(0.0::numeric, ROUND((median_price * (1::numeric - discount_rate::numeric))::numeric, 2)) AS unit_price,
       'THB', 'monthly', 'piecewise', TRUE, current_date - 30
FROM vendor_x_cat x
CROSS JOIN (
  VALUES (1, 0.05::numeric), (100, 0.15::numeric), (300, 0.30::numeric)
) t(threshold_qty, discount_rate)
WHERE x.median_price IS NOT NULL
ON CONFLICT DO NOTHING;

-- =====================
-- 4) Seed Switching Cost Policy per cluster
-- =====================

INSERT INTO switching_cost_policy(feature_cluster_key, training_cost_per_user, migration_flat_cost, early_termination_penalty_rate)
SELECT c.feature_cluster_key,
       CASE c.feature_cluster_key
         WHEN 'DevOps' THEN 1200
         WHEN 'Security' THEN 1500
         WHEN 'ITSM' THEN 1400
         WHEN 'CRM' THEN 1300
         WHEN 'Analytics' THEN 1000
         ELSE 900
       END::numeric AS training_cost_per_user,
       CASE c.feature_cluster_key
         WHEN 'DevOps' THEN 200000
         WHEN 'Security' THEN 250000
         WHEN 'ITSM' THEN 220000
         WHEN 'CRM' THEN 180000
         WHEN 'Analytics' THEN 160000
         ELSE 120000
       END::numeric AS migration_flat_cost,
       0.15::numeric AS early_termination_penalty_rate
FROM (
  SELECT DISTINCT category AS feature_cluster_key FROM apps WHERE category IS NOT NULL
) c
ON CONFLICT (feature_cluster_key) DO NOTHING;

-- =====================
-- 5) Quick sanity checks (optional selects)
-- =====================

-- SELECT 'app tiers' AS t, COUNT(*) FROM app_volume_pricing_tiers;
-- SELECT 'vendor tiers' AS t, COUNT(*) FROM vendor_volume_pricing_tiers;
-- SELECT 'switching policies' AS t, COUNT(*) FROM switching_cost_policy;


