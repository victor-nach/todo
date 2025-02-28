CREATE TABLE IF NOT EXISTS todos (
    id          UUID PRIMARY KEY, 
    title       TEXT NOT NULL,                            
    description TEXT,                                     
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    updated_at  TIMESTAMP                                 
);
