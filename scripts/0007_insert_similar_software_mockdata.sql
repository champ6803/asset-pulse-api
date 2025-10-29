    -- ===========================================================
    -- Similar Software Mock Data
    -- สำหรับ Similar Software Detection feature
    -- สร้างข้อมูล mock สำหรับ group_consolidation_opps และเชื่อมกับ license_assignments
    -- ===========================================================
    SET search_path TO public;

    -- Clean up existing similar software mock data
    DELETE FROM group_consolidation_opps WHERE rationale LIKE '%Similar Software%' OR rationale LIKE '%consolidation%';

    -- ===========================================================
    -- Step 1: Ensure apps and vendors exist
    -- Note: Vendors are created but not directly linked to apps
    -- (vendors are linked to apps through contracts/contract_terms)
    -- ===========================================================

    -- Get or create vendors
    DO $$
    DECLARE
        v_atlassian_id BIGINT;
        v_asana_id BIGINT;
        v_figma_id BIGINT;
        v_adobe_id BIGINT;
        v_tableau_id BIGINT;
        v_microsoft_id BIGINT;
    BEGIN
        -- Atlassian
        SELECT id INTO v_atlassian_id FROM vendors WHERE name = 'Atlassian' LIMIT 1;
        IF v_atlassian_id IS NULL THEN
            INSERT INTO vendors (name) VALUES ('Atlassian') RETURNING id INTO v_atlassian_id;
        END IF;

        -- Asana
        SELECT id INTO v_asana_id FROM vendors WHERE name = 'Asana' LIMIT 1;
        IF v_asana_id IS NULL THEN
            INSERT INTO vendors (name) VALUES ('Asana') RETURNING id INTO v_asana_id;
        END IF;

        -- Figma Inc.
        SELECT id INTO v_figma_id FROM vendors WHERE name = 'Figma Inc.' LIMIT 1;
        IF v_figma_id IS NULL THEN
            INSERT INTO vendors (name) VALUES ('Figma Inc.') RETURNING id INTO v_figma_id;
        END IF;

        -- Adobe
        SELECT id INTO v_adobe_id FROM vendors WHERE name = 'Adobe' LIMIT 1;
        IF v_adobe_id IS NULL THEN
            INSERT INTO vendors (name) VALUES ('Adobe') RETURNING id INTO v_adobe_id;
        END IF;

        -- Tableau
        SELECT id INTO v_tableau_id FROM vendors WHERE name = 'Tableau' LIMIT 1;
        IF v_tableau_id IS NULL THEN
            INSERT INTO vendors (name) VALUES ('Tableau') RETURNING id INTO v_tableau_id;
        END IF;

        -- Microsoft
        SELECT id INTO v_microsoft_id FROM vendors WHERE name = 'Microsoft' LIMIT 1;
        IF v_microsoft_id IS NULL THEN
            INSERT INTO vendors (name) VALUES ('Microsoft') RETURNING id INTO v_microsoft_id;
        END IF;

        -- Create apps if they don't exist
        -- Jira Software
        IF NOT EXISTS (SELECT 1 FROM apps WHERE name = 'Jira Software') THEN
            INSERT INTO apps (company_code, key, name, alias, status, category, application_tier)
            VALUES (NULL, 'jira-software', 'Jira Software', 'Jira', 'Active', 'Project Management', 'Enterprise')
            ON CONFLICT (key) DO NOTHING;
        END IF;

        -- Asana
        IF NOT EXISTS (SELECT 1 FROM apps WHERE name = 'Asana') THEN
            INSERT INTO apps (company_code, key, name, alias, status, category, application_tier)
            VALUES (NULL, 'asana', 'Asana', 'Asana', 'Active', 'Project Management', 'Enterprise')
            ON CONFLICT (key) DO NOTHING;
        END IF;

        -- Figma
        IF NOT EXISTS (SELECT 1 FROM apps WHERE name LIKE '%Figma%' OR name = 'Figma') THEN
            INSERT INTO apps (company_code, key, name, alias, status, category, application_tier)
            VALUES (NULL, 'figma', 'Figma', 'Figma', 'Active', 'Design', 'Enterprise')
            ON CONFLICT (key) DO NOTHING;
        END IF;

        -- Adobe XD
        IF NOT EXISTS (SELECT 1 FROM apps WHERE name LIKE '%Adobe XD%' OR name LIKE '%Adobe%XD%') THEN
            INSERT INTO apps (company_code, key, name, alias, status, category, application_tier)
            VALUES (NULL, 'adobe-xd', 'Adobe XD', 'Adobe XD', 'Active', 'Design', 'Enterprise')
            ON CONFLICT (key) DO NOTHING;
        END IF;

        -- Tableau
        IF NOT EXISTS (SELECT 1 FROM apps WHERE name = 'Tableau') THEN
            INSERT INTO apps (company_code, key, name, alias, status, category, application_tier)
            VALUES (NULL, 'tableau', 'Tableau', 'Tableau', 'Active', 'Analytics', 'Enterprise')
            ON CONFLICT (key) DO NOTHING;
        END IF;

        -- Power BI
        IF NOT EXISTS (SELECT 1 FROM apps WHERE name LIKE '%Power BI%' OR name = 'Power BI') THEN
            INSERT INTO apps (company_code, key, name, alias, status, category, application_tier)
            VALUES (NULL, 'power-bi', 'Power BI', 'Power BI', 'Active', 'Analytics', 'Enterprise')
            ON CONFLICT (key) DO NOTHING;
        END IF;

    END $$;

    -- ===========================================================
    -- Step 2: Create license assignments for these apps across multiple subsidiaries
    -- This simulates real usage across different companies
    -- ===========================================================

    -- Get company codes from existing companies
    WITH company_codes AS (
        SELECT DISTINCT code FROM companies WHERE code IN ('SCBX', 'SCB', 'DATAX', 'SCBAM', 'INVX') LIMIT 5
    ),
    target_apps AS (
        SELECT id, name FROM apps WHERE name IN ('Jira Software', 'Asana', 'Figma', 'Adobe XD', 'Tableau', 'Power BI')
    )
    INSERT INTO license_assignments (company_code, user_id, app_id, license_tier, assignment_source, assigned_at)
    SELECT 
        cc.code AS company_code,
        u.id AS user_id,
        ta.id AS app_id,
        'Enterprise' AS license_tier,
        'mock_similar_software' AS assignment_source,
        NOW() - (random() * INTERVAL '365 days') AS assigned_at
    FROM company_codes cc
    CROSS JOIN target_apps ta
    CROSS JOIN LATERAL (
        SELECT id FROM users 
        WHERE company_code = cc.code 
        AND status = 'active'
        ORDER BY random()
        LIMIT (CASE 
            WHEN ta.name = 'Jira Software' THEN 180
            WHEN ta.name = 'Asana' THEN 60
            WHEN ta.name = 'Figma' THEN 70
            WHEN ta.name = 'Adobe XD' THEN 30
            WHEN ta.name = 'Tableau' THEN 120
            WHEN ta.name = 'Power BI' THEN 90
            ELSE 50
        END / (SELECT COUNT(*)::int FROM company_codes))
    ) u
    WHERE NOT EXISTS (
        SELECT 1 FROM license_assignments la
        WHERE la.company_code = cc.code 
        AND la.app_id = ta.id 
        AND la.user_id = u.id
        AND la.revoked_at IS NULL
    );

    -- ===========================================================
    -- Step 3: Insert consolidation opportunities
    -- ===========================================================

    INSERT INTO group_consolidation_opps (company_code, app_id, potential_saving_amt, rationale, status)
    SELECT 
        NULL AS company_code,
        a.id AS app_id,
        CASE 
            WHEN a.name = 'Jira Software' THEN 46000
            WHEN a.name = 'Asana' THEN 32000
            WHEN a.name = 'Figma' THEN 55000
            WHEN a.name = 'Adobe XD' THEN 42000
            WHEN a.name = 'Tableau' THEN 78000
            WHEN a.name = 'Power BI' THEN 65000
            ELSE 30000
        END AS potential_saving_amt,
        CASE 
            WHEN a.name = 'Jira Software' THEN 'Multiple subsidiaries use Jira Software. Consolidation would allow volume pricing and reduce maintenance overhead.'
            WHEN a.name = 'Asana' THEN 'Asana is used across multiple subsidiaries. Consolidating to a single vendor (Atlassian) would provide better pricing tiers.'
            WHEN a.name = 'Figma' THEN 'Figma usage is widespread across design teams. Consolidation would enable enterprise pricing and shared team libraries.'
            WHEN a.name = 'Adobe XD' THEN 'Adobe XD and Figma serve similar purposes. Consider consolidating to one design tool for consistency.'
            WHEN a.name = 'Tableau' THEN 'Tableau is deployed across several subsidiaries. Group licensing would significantly reduce costs per seat.'
            WHEN a.name = 'Power BI' THEN 'Power BI usage overlaps with Tableau. Consider consolidating BI tools for better licensing terms.'
            ELSE 'Similar software consolidation opportunity'
        END AS rationale,
        'pending' AS status
    FROM apps a
    WHERE a.name IN ('Jira Software', 'Asana', 'Figma', 'Adobe XD', 'Tableau', 'Power BI')
    AND NOT EXISTS (
        SELECT 1 FROM group_consolidation_opps gco 
        WHERE gco.app_id = a.id AND gco.company_code IS NULL
    );

    -- ===========================================================
    -- Verification Query
    -- ===========================================================
    SELECT 
        gco.id,
        gco.company_code,
        a.name AS app_name,
        a.category,
        COUNT(DISTINCT la.company_code) AS subsidiaries_count,
        COUNT(DISTINCT la.user_id) AS total_users,
        gco.potential_saving_amt,
        gco.status
    FROM group_consolidation_opps gco
    LEFT JOIN apps a ON a.id = gco.app_id
    LEFT JOIN license_assignments la ON la.app_id = a.id AND la.revoked_at IS NULL
    WHERE gco.company_code IS NULL
    GROUP BY gco.id, gco.company_code, a.name, a.category, gco.potential_saving_amt, gco.status
    ORDER BY gco.potential_saving_amt DESC NULLS LAST;

