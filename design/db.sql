CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    profile_avatar TEXT,
    profile_cover TEXT,
    bio TEXT,
    status VARCHAR(50) DEFAULT 'active', -- 'active', 'inactive', 'banned'
    role VARCHAR(50) DEFAULT 'user', -- 'user', 'admin', etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    author_id INT REFERENCES users(id) ON DELETE CASCADE,
    body TEXT,
    image_url TEXT,
    share_state VARCHAR(20), -- Public, Private, Circles, etc.
    likes_count INT DEFAULT 0,
    comments_count INT DEFAULT 0,
    hashtags TEXT[], -- Array of hashtags
    mentioned_users INT[], -- Array of user IDs mentioned
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX(author_id),
    INDEX(hashtags)
);


CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INT REFERENCES posts(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id INT REFERENCES comments(id) ON DELETE CASCADE, -- For replies
    comment_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX(post_id),
    INDEX(user_id)
);

CREATE TABLE likes (
    id SERIAL PRIMARY KEY,
    post_id INT REFERENCES posts(id) ON DELETE CASCADE,
    comment_id INT REFERENCES comments(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(post_id, user_id), -- Ensure a user can only like a post once
    UNIQUE(comment_id, user_id),
    INDEX(post_id),
    INDEX(comment_id)
);

CREATE TABLE shares (
    id SERIAL PRIMARY KEY,
    post_id INT REFERENCES posts(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    share_text TEXT, -- Optional text added when sharing
    shared_with INT[], -- Array of user IDs or circle IDs
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE followers (
    follower_id INT REFERENCES users(id) ON DELETE CASCADE,
    following_id INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(follower_id, following_id), -- Prevent duplicate follow entries
    INDEX(follower_id),
    INDEX(following_id)
);

CREATE TABLE circles (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    circle_name VARCHAR(100),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE circle_members (
    circle_id INT REFERENCES circles(id) ON DELETE CASCADE,
    member_id INT REFERENCES users(id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(circle_id, member_id)
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    recipient_id INT REFERENCES users(id) ON DELETE CASCADE,
    sender_id INT REFERENCES users(id),
    post_id INT REFERENCES posts(id),
    comment_id INT REFERENCES comments(id),
    type VARCHAR(50), -- e.g., 'like', 'comment', 'follow'
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);