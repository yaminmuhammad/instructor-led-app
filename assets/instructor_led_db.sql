CREATE DATABASE instructor_led_db;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_type AS ENUM ('admin', 'participant', 'trainer');

CREATE TYPE absent_type AS ENUM('Present', 'Not Present');

CREATE TYPE question_status AS ENUM ('Finished', 'Process');

CREATE TYPE participant_type AS ENUM ('Basic', 'Advance');

CREATE TABLE users (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) UNIQUE,
  username VARCHAR(100) UNIQUE,
  address TEXT,
  hash_password VARCHAR(100),
  role user_type,
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE trainers (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  phone_number VARCHAR(13) UNIQUE,
  user_id uuid UNIQUE,
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE TABLE specializations (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  trainer_id uuid NOT NULL,
  name VARCHAR(100),
  FOREIGN KEY ("trainer_id") REFERENCES "trainers" ("id")
);

CREATE TABLE participants (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  date_of_birth DATE,
  place_of_birth VARCHAR(100),
  last_education VARCHAR(100),
  user_id uuid UNIQUE,
  role participant_type,
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE TABLE schedules (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  activity VARCHAR(45),
  date DATE DEFAULT CURRENT_DATE,
  trainer_id uuid NOT NULL,
  participant_id uuid NOT NULL,
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("trainer_id") REFERENCES "trainers" ("id"),
  FOREIGN KEY ("participant_id") REFERENCES "participants" ("id")
);

CREATE TABLE absences (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  date DATE DEFAULT CURRENT_DATE,
  information TEXT,
  absence_status absent_type,
  absence_time TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  participant_id uuid NOT NULL,
  trainer_id uuid NOT NULL,
  schedule_id uuid NOT NULL,
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("participant_id") REFERENCES "participants" ("id"),
  FOREIGN KEY ("trainer_id") REFERENCES "trainers" ("id"),
  FOREIGN KEY ("schedule_id") REFERENCES "schedules" ("id")
);

CREATE TABLE questions (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  question TEXT,
  answer TEXT,
  STATUS question_status,
  participant_id uuid NOT NULL,
  trainer_id uuid NOT NULL,
  schedule_id uuid NOT NULL,
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("participant_id") REFERENCES "participants" ("id"),
  FOREIGN KEY ("trainer_id") REFERENCES "trainers" ("id"),
  FOREIGN KEY ("schedule_id") REFERENCES "schedules" ("id")
);

CREATE TABLE schedule_images (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  schedule_id uuid NOT NULL,
  file_name VARCHAR(100),
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP(0),
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP(0),
  FOREIGN KEY ("schedule_id") REFERENCES "schedules" ("id")
);

INSERT INTO
  users(name, email, username, address, hash_password, role)
VALUES
  (
    'Iqi Tes',
    'iqi@mail.com',
    'iqi',
    'Cirebon',
    'password',
    'participant'
  );

INSERT INTO
  participants(
    date_of_birth,
    place_of_birth,
    last_education,
    user_id,
    role
  )
VALUES
  (
    '1999-10-10',
    'Jakarta',
    'Universitas Gadjah Mada',
    '37881fb4-5ca4-4939-a1ad-d8a36fbb23a6',
    'Advance'
  );

INSERT INTO
  schedules(
    activity,
    date,
    trainer_id,
    participant_id
  )
VALUES
  (
    'Training',
    '2023-12-23',
    'e0da3c80-17b4-41df-8087-5e2fbc19f654',
    'c112bc5c-29c8-4921-a85a-505b21b97b1d'
  );
