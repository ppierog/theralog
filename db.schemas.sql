-- Create DB :
-- sqlite3 db_name < db.schemas.sql
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
   salt              CHAR(8)  NOT NULL,
   password          TEXT     NOT_NULL,
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

CREATE TABLE PatientManifest(
   -- rowid INT PRIMARY KEY NOT NULL, Added by sqlite3 by default
   patient_id     INT    NOT NULL,
   user_id        INT    NOT NULL,
   -- create 0x1, read 0x2, update 0x4, delate 0x8 - full access -> 0x1 | 0x2 | 0x4 | 0x8 = 0xF
   crud_mask      INT    NOT NULL,
   encrypted_aes  TEXT,
   FOREIGN KEY(patient_id) REFERENCES Patient(rowid),
   FOREIGN KEY(user_id) REFERENCES User(rowid)
);


-- SELECT rowid,* FROM Patient;
-- SELECT rowid, name, birth_year, city, telephone_number from Patient;

CREATE TABLE Note(
   -- rowid INT PRIMARY KEY NOT NULL, Added by sqlite3 by default
   name           TEXT,    -- description
   patient_id     INT,
   session_date   INT,     -- unix time
   note_date      INT,     -- unix time
   file_name      TEXT,    -- file with note
   is_crypted     BOOLEAN,
   CONSTRAINT fk_patients
      FOREIGN KEY (patient_id)
      REFERENCES Patient(rowid)
      ON DELETE CASCADE
);
