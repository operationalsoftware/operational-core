ALTER TABLE pdf_print_log
    DROP COLUMN printer_id;

ALTER TABLE print_requirement
    DROP COLUMN printer_id;

DROP INDEX IF EXISTS print_requirement_unique_printer_id;

ALTER TABLE pdf_print_log
    ALTER COLUMN printer_name SET NOT NULL;

ALTER TABLE print_requirement
    ALTER COLUMN printer_name SET NOT NULL;

DROP INDEX IF EXISTS print_requirement_unique_printer_name;
