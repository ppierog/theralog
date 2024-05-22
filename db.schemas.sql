-- Create DB :
-- sqlite3 db_name < db.schema.sql
-- Cmds :
-- sqlite3>
-- .tables
-- .schema
-- .quit

CREATE TABLE User(
   -- rowid INT PRIMARY KEY NOT NULL, Added by sqlite3 by default
   name              TEXT     NOT NULL,
   last_name         TEXT     NOT NULL,
   email             TEXT,
   telephone_number  CHAR(12),
   password_salt     CHAR(8)  NOT NULL,
   password_sha256   TEXT     NOT_NULL,
   pub_key           TEXT
);

CREATE TABLE Patient(
   -- rowid INT PRIMARY KEY NOT NULL, Added by sqlite3 by default
   name           TEXT    NOT NULL,
   occupation     TEXT,
   birth_year     INT,
   city           TEXT,
   telephone_number CHAR(12)
);
-- SELECT rowid,* FROM Patient;
-- SELECT rowid, name, birth_year, city, telephone_number from Patient;

CREATE TABLE Note(
   -- rowid INT PRIMARY KEY NOT NULL, Added by sqlite3 by default
   patient_id     INT,
   description    TEXT,
   session_date   INT,  -- unix time
   note_date      INT,  -- unix time
   file_name      TEXT, -- file with note
   FOREIGN KEY(patient_id) REFERENCES Patient(rowid)
);