DO $$ DECLARE
    r RECORD;
BEGIN
    FOR r IN (
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'public' AND tablename NOT IN (
            SELECT c.relname
            FROM pg_depend d
            JOIN pg_extension e ON d.refobjid = e.oid
            JOIN pg_class c ON d.objid = c.oid
            WHERE d.deptype = 'e'
        )
    ) LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
END $$;