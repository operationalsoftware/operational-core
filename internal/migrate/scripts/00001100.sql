CREATE TABLE pdf_generation_log (
    pdf_generation_log_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    template_name TEXT NOT NULL,
    input_data TEXT NOT NULL,
    file_id UUID REFERENCES file(file_id),
    pdf_title TEXT,
    created_by INT REFERENCES app_user(user_id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE pdf_print_log (
    pdf_print_log_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pdf_generation_log_id INTEGER NOT NULL REFERENCES pdf_generation_log(pdf_generation_log_id) ON DELETE CASCADE,
    template_name TEXT NOT NULL,
    requirement_name TEXT,
    printer_id INT,
    printer_name TEXT,
    printnode_job_id INT,
    status TEXT,
    error_message TEXT,
    created_by INT REFERENCES app_user(user_id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
