CREATE TABLE message_files (
    message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    PRIMARY KEY (message_id, file_id)
);