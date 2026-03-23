CREATE TABLE task_collaborators (
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Composite Primary Key mencegah 1 user ditambahkan 2x di task yang sama
    PRIMARY KEY (task_id, user_id),
    
    -- Foreign Key dengan efek berantai
    CONSTRAINT fk_tc_task FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE,
    CONSTRAINT fk_tc_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);