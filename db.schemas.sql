-- Create DB :
-- sqlite3 db_name < db.schema.sql
-- Cmds :
-- sqlite3>
-- .tables
-- .schema
-- .quit

CREATE TABLE Patient(
   -- rowid INT PRIMARY KEY NOT NULL, Added by sqlite3 by default 
   name           TEXT    NOT NULL,
   occupation     TEXT,
   birth_year     INT,
   city           TEXT,
   telephone_number CHAR(12)
);

-- SELECT rowid,* FROM Patient;
-- SELECT rowid, name, birth_year, city from Patient;
-- DROP Patient;
-- DELETE FROM Patient;
