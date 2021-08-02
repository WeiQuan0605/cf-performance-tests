-- FUNC DEF:
CREATE FUNCTION create_test_orgs(
    orgs INTEGER
) RETURNS void AS
$$
DECLARE
org_guid text;
BEGIN
FOR i IN 1..orgs
        LOOP
            org_guid := gen_random_uuid();
INSERT INTO organizations (guid, name, quota_definition_id)
VALUES (org_guid, 'perf-test-org-' || org_guid, 1);
END LOOP;
END;
$$ LANGUAGE plpgsql;

-- ============================================================= --

-- FUNC DEF:
CREATE FUNCTION create_shared_domains(
    shared_domains INTEGER
) RETURNS void AS
$$
DECLARE
shared_domain_guid text;
BEGIN
FOR i IN 1..shared_domains
            LOOP
                shared_domain_guid := gen_random_uuid();
INSERT INTO domains (guid, name)
VALUES (shared_domain_guid, 'perf-test-shared-domain-' || shared_domain_guid);
END LOOP;
END;
$$ LANGUAGE plpgsql;

-- ============================================================= --

-- FUNC DEF:
CREATE FUNCTION create_private_domains(
    private_domains INTEGER
) RETURNS void AS
$$
DECLARE
private_domain_guid text;
BEGIN
FOR i IN 1..private_domains
            LOOP
                private_domain_guid := gen_random_uuid();
INSERT INTO domains (guid, name, owning_organization_id)
SELECT private_domain_guid, 'perf-test-private-domain-' || private_domain_guid, id
FROM organizations WHERE name != 'default' ORDER BY random() LIMIT 1;
END LOOP;
END;
$$ LANGUAGE plpgsql;