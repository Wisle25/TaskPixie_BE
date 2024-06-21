CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    project_id UUID,
    user_id UUID NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    
    CONSTRAINT fk_task_project FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE SET NULL,
    CONSTRAINT fk_task_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
