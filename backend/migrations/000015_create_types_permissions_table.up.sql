CREATE TABLE IF NOT EXISTS types_permissions(
    user_type_id int NOT NULL REFERENCES user_types(user_type_id) ON DELETE CASCADE,
    permission_id int NOT NULL REFERENCES permissions(permission_id) ON DELETE CASCADE,
    PRIMARY KEY(user_type_id, permission_id)
);

INSERT INTO types_permissions (user_type_id, permission_id)
VALUES 
    (2, 1),
    (1, 2),
    (1, 3),
    (1, 4),
    (2, 5),
    (1, 6),
    (2, 7);