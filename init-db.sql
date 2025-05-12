-- Create databases
CREATE DATABASE auth;
CREATE DATABASE pds;
CREATE DATABASE bgs;

-- Connect to auth database
\c auth

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    did TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index on username
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_did ON users(did);

-- Connect to pds database
\c pds

-- Create repositories table
CREATE TABLE repositories (
    did TEXT PRIMARY KEY,
    head TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create commits table
CREATE TABLE commits (
    id TEXT PRIMARY KEY,
    repository_did TEXT NOT NULL REFERENCES repositories(did),
    prev TEXT,
    data BYTEA NOT NULL,
    signature BYTEA,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index on repository_did
CREATE INDEX idx_commits_repository_did ON commits(repository_did);

-- Create documents table
CREATE TABLE documents (
    id TEXT NOT NULL,
    repository_did TEXT NOT NULL REFERENCES repositories(did),
    type TEXT NOT NULL,
    value JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (repository_did, id)
);

-- Create index on repository_did and type
CREATE INDEX idx_documents_repository_did_type ON documents(repository_did, type);

-- Connect to bgs database
\c bgs

-- Create follows table
CREATE TABLE follows (
    follower TEXT NOT NULL,
    following TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower, following)
);

-- Create index on follower and following
CREATE INDEX idx_follows_follower ON follows(follower);
CREATE INDEX idx_follows_following ON follows(following);
