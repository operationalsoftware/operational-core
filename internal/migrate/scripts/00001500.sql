ALTER TABLE pdf_print_log
    DROP CONSTRAINT IF EXISTS pdf_print_log_pdf_generation_log_id_fkey;

ALTER TABLE pdf_print_log
    ADD CONSTRAINT pdf_print_log_pdf_generation_log_id_fkey
    FOREIGN KEY (pdf_generation_log_id)
    REFERENCES pdf_generation_log(pdf_generation_log_id);

ALTER TABLE pdf_print_log
    DROP COLUMN IF EXISTS printer_name;

ALTER TABLE pdf_print_log
    DROP COLUMN IF EXISTS requirement_name;

ALTER TABLE pdf_print_log
    DROP COLUMN IF EXISTS status;

ALTER TABLE pdf_print_log
    ADD COLUMN print_requirement_id INT NOT NULL REFERENCES print_requirement(print_requirement_id);

ALTER TABLE pdf_print_log
    ALTER COLUMN printnode_job_id TYPE BIGINT;

ALTER TABLE pdf_generation_log
    ADD COLUMN print_node_options JSONB NOT NULL;

ALTER TABLE pdf_generation_log
    ALTER COLUMN input_data TYPE JSONB USING input_data::jsonb;

ALTER TABLE pdf_generation_log
    ALTER COLUMN pdf_title SET NOT NULL;

ALTER TABLE pdf_generation_log
    ALTER COLUMN created_by SET NOT NULL;

ALTER TABLE pdf_print_log
    ALTER COLUMN created_by SET NOT NULL;

ALTER TABLE print_requirement
    ALTER COLUMN assigned_by SET NOT NULL;

CREATE OR REPLACE VIEW pdf_generation_log_view AS
SELECT
    gl.pdf_generation_log_id,
    gl.template_name,
    gl.input_data,
    gl.file_id::text AS file_id,
    gl.pdf_title,
    gl.print_node_options,
    u.username AS created_by_username,
    gl.created_at
FROM pdf_generation_log gl
INNER JOIN app_user u ON gl.created_by = u.user_id;

CREATE OR REPLACE VIEW pdf_print_log_view AS
SELECT
    pl.pdf_print_log_id,
    pl.pdf_generation_log_id,
    pl.template_name,
    pr.requirement_name,
    pl.printnode_job_id,
    pl.error_message,
    u.username AS created_by_username,
    pl.created_at,
    gl.file_id::text AS file_id,
    gl.pdf_title,
    gl.input_data
FROM pdf_print_log pl
LEFT JOIN pdf_generation_log gl ON gl.pdf_generation_log_id = pl.pdf_generation_log_id
LEFT JOIN print_requirement pr ON pr.print_requirement_id = pl.print_requirement_id
INNER JOIN app_user u ON pl.created_by = u.user_id;
