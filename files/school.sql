CREATE TABLE admin (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	username STRING NOT NULL,
	password STRING NOT NULL,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	is_active BOOL NOT NULL DEFAULT true,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	FAMILY "primary" (id, username, password, created_by, created_at, updated_by, updated_at, is_active)
);

CREATE TABLE faculty (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	code INT8 NOT NULL,
	abbreviation STRING NOT NULL,
	name STRING NOT NULL,
	description STRING NOT NULL,
	is_delete BOOL NOT NULL DEFAULT false,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	FAMILY "primary" (id, code, abbreviation, name, description, is_delete, created_by, created_at, updated_by, updated_at)
);

CREATE TABLE program (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	faculty_id UUID NOT NULL,
	name STRING NOT NULL,
	code INT8 NOT NULL,
	description STRING NOT NULL,
	is_delete BOOL NOT NULL DEFAULT false,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	INDEX program_faculty_id_idx (faculty_id ASC),
	FAMILY "primary" (id, faculty_id, name, code, description, is_delete, created_by, created_at, updated_by, updated_at)
);

CREATE TABLE student (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	name STRING NOT NULL,
	address STRING NULL,
	date_of_birth TIMESTAMP NOT NULL,
	gender INT8 NOT NULL,
	email STRING NOT NULL,
	phone_no STRING NOT NULL,
	is_active BOOL NOT NULL DEFAULT true,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	password STRING NOT NULL DEFAULT '':::STRING,
	student_code STRING NOT NULL DEFAULT '':::STRING,
	program_id UUID NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	INDEX student_program_id_idx (program_id ASC),
	INDEX student_auto_index_student_fk (program_id ASC),
	FAMILY "primary" (id, name, address, date_of_birth, gender, email, phone_no, is_active, created_by, created_at, updated_by, updated_at, password, student_code, program_id)
);

CREATE TABLE subject (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	name STRING NOT NULL,
	description STRING NOT NULL,
	duration INT8 NOT NULL DEFAULT 0:::INT8,
	is_delete BOOL NOT NULL DEFAULT false,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	FAMILY "primary" (id, name, description, duration, is_delete, created_by, created_at, updated_by, updated_at)
);

CREATE TABLE lecturer (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	name STRING NOT NULL,
	address STRING NULL,
	phone_no STRING NOT NULL,
	is_active BOOL NOT NULL DEFAULT true,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	email STRING NULL,
	password STRING NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	FAMILY "primary" (id, name, address, phone_no, is_active, created_by, created_at, updated_by, updated_at, email, password)
);

CREATE TABLE classroom (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	faculty_id UUID NOT NULL,
	floor INT8 NOT NULL,
	room_no INT8 NOT NULL,
	code STRING NOT NULL,
	is_delete BOOL NOT NULL DEFAULT false,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	INDEX classroom_faculty_id_idx (faculty_id ASC),
	INDEX classroom_auto_index_classroom_fk (faculty_id ASC),
	FAMILY "primary" (id, faculty_id, floor, room_no, code, is_delete, created_by, created_at, updated_by, updated_at)
);

CREATE TABLE intake (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	year STRING NOT NULL,
	month INT8 NOT NULL,
	is_delete BOOL NOT NULL DEFAULT false,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	trimester INT8 NOT NULL,
	start_date TIMESTAMPTZ NOT NULL,
	end_date TIMESTAMPTZ NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	FAMILY "primary" (id, year, month, is_delete, created_by, created_at, updated_by, updated_at, trimester, start_date, end_date)
);

CREATE TABLE session (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	subject_id UUID NOT NULL,
	lecturer_id UUID NOT NULL,
	intake_id UUID NOT NULL,
	is_delete BOOL NOT NULL DEFAULT false,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	classroom_id UUID NOT NULL,
	program_id UUID NOT NULL,
	day INT8 NOT NULL,
	start_time TIME NOT NULL,
	end_time TIME NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	INDEX timetable_subject_id_idx (subject_id ASC, lecturer_id ASC, intake_id ASC, classroom_id ASC, program_id ASC),
	INDEX timetable_auto_index_timetable_fk (subject_id ASC),
	INDEX timetable_auto_index_timetable_fk_1 (lecturer_id ASC),
	INDEX timetable_auto_index_timetable_fk_2 (program_id ASC),
	INDEX timetable_auto_index_timetable_fk_3 (classroom_id ASC),
	INDEX timetable_auto_index_timetable_fk_4 (intake_id ASC),
	FAMILY "primary" (id, subject_id, lecturer_id, intake_id, is_delete, created_by, created_at, updated_by, updated_at, classroom_id, program_id, day, start_time, end_time)
);

CREATE TABLE student_enroll (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	"session_ID" UUID NOT NULL,
	student_id UUID NOT NULL,
	is_delete BOOL NOT NULL DEFAULT false,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	INDEX student_enroll_timetable_id_idx ("session_ID" ASC, student_id ASC),
	INDEX student_enroll_auto_index_student_enroll_fk (student_id ASC),
	INDEX student_enroll_auto_index_student_enroll_fk_1 ("session_ID" ASC),
	FAMILY "primary" (id, "session_ID", student_id, is_delete, created_by, created_at, updated_by, updated_at)
);

CREATE TABLE attendance (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	student_id UUID NOT NULL,
	student_enroll_id UUID NOT NULL,
	is_attend BOOL NOT NULL DEFAULT false,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	INDEX attendance_student_id_idx (student_id ASC, student_enroll_id ASC),
	INDEX attendance_auto_index_attendance_fk_1 (student_enroll_id ASC),
	FAMILY "primary" (id, student_id, student_enroll_id, is_attend, created_by, created_at, updated_by, updated_at)
);

CREATE TABLE result (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	student_enroll_id UUID NOT NULL,
	grade VARCHAR(1) NOT NULL,
	marks INT8 NOT NULL,
	is_delete BOOL NOT NULL DEFAULT false,
	created_by UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_by UUID NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	INDEX result_student_enroll_id_idx (student_enroll_id ASC),
	INDEX result_auto_index_result_fk (student_enroll_id ASC),
	FAMILY "primary" (id, student_enroll_id, grade, marks, is_delete, created_by, created_at, updated_by, updated_at)
);

INSERT INTO admin (id, username, password, created_by, created_at, updated_by, updated_at, is_active) VALUES
	('6517ea2d-8d78-4a3f-83a1-4877d82ced59', 'admin', '$2y$12$4HSCAV33u.2lmzbdGh/CP.5F2VJjM.9NqrnuvNIsY4.mvd.rFMuVW', '2e1910ed-1951-42a2-841d-e6e9527f0448', '2020-03-23 07:39:52.574076+00:00', NULL, NULL, true);

ALTER TABLE program ADD CONSTRAINT program_fk FOREIGN KEY (faculty_id) REFERENCES faculty(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE student ADD CONSTRAINT student_fk FOREIGN KEY (program_id) REFERENCES program(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE classroom ADD CONSTRAINT classroom_fk FOREIGN KEY (faculty_id) REFERENCES faculty(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE session ADD CONSTRAINT timetable_fk FOREIGN KEY (subject_id) REFERENCES subject(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE session ADD CONSTRAINT timetable_fk_1 FOREIGN KEY (lecturer_id) REFERENCES lecturer(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE session ADD CONSTRAINT timetable_fk_2 FOREIGN KEY (program_id) REFERENCES program(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE session ADD CONSTRAINT timetable_fk_3 FOREIGN KEY (classroom_id) REFERENCES classroom(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE session ADD CONSTRAINT timetable_fk_4 FOREIGN KEY (intake_id) REFERENCES intake(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE student_enroll ADD CONSTRAINT student_enroll_fk FOREIGN KEY (student_id) REFERENCES student(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE student_enroll ADD CONSTRAINT student_enroll_fk_1 FOREIGN KEY ("session_ID") REFERENCES session(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE attendance ADD CONSTRAINT attendance_fk FOREIGN KEY (student_id) REFERENCES student(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE attendance ADD CONSTRAINT attendance_fk_1 FOREIGN KEY (student_enroll_id) REFERENCES student_enroll(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE result ADD CONSTRAINT result_fk FOREIGN KEY (student_enroll_id) REFERENCES student_enroll(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- Validate foreign key constraints. These can fail if there was unvalidated data during the dump.
ALTER TABLE program VALIDATE CONSTRAINT program_fk;
ALTER TABLE student VALIDATE CONSTRAINT student_fk;
ALTER TABLE classroom VALIDATE CONSTRAINT classroom_fk;
ALTER TABLE session VALIDATE CONSTRAINT timetable_fk;
ALTER TABLE session VALIDATE CONSTRAINT timetable_fk_1;
ALTER TABLE session VALIDATE CONSTRAINT timetable_fk_2;
ALTER TABLE session VALIDATE CONSTRAINT timetable_fk_3;
ALTER TABLE session VALIDATE CONSTRAINT timetable_fk_4;
ALTER TABLE student_enroll VALIDATE CONSTRAINT student_enroll_fk;
ALTER TABLE student_enroll VALIDATE CONSTRAINT student_enroll_fk_1;
ALTER TABLE attendance VALIDATE CONSTRAINT attendance_fk;
ALTER TABLE attendance VALIDATE CONSTRAINT attendance_fk_1;
ALTER TABLE result VALIDATE CONSTRAINT result_fk;
