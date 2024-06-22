-- Enable the uuid-ossp extension to generate UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    title VARCHAR(50) NOT NULL,
    detail TEXT,
    priority VARCHAR(15),
    status VARCHAR(15),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create the project_members table
CREATE TABLE project_members (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, user_id)
);

-- Create indexes to speed up queries
CREATE INDEX idx_project_owner ON projects(owner_id);
CREATE INDEX idx_project_members_user ON project_members(user_id);
