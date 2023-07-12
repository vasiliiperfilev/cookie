INSERT INTO users (email, name, password_hash, user_type_id, image_id)
VALUES 
    ('test1@user', 'name 1', 'hash', 1, 'test1'),
    ('test2@user', 'name 2', 'hash', 2, 'test2'),
    ('test3@user', 'name 3', 'hash', 1, 'test3'),
    ('test4@user', 'name 4', 'hash', 2, 'test4'),
    ('testItemModel@user', 'item model', 'hash', 1, 'test4'),
    ('testOrderModel@user', 'order model', 'hash', 1, 'test4');

INSERT INTO conversations (last_message_id)
VALUES 
    (0),
    (0);

INSERT INTO conversations_users (conversation_id, user_id)
VALUES 
    (1, 2),
    (1, 6),
    (2, 6),
    (2, 4);

INSERT INTO items (supplier_id, unit_id, size, name, image_url)
VALUES
    (6, 1, 1, 'Test product', 'test'),
    (6, 2, 3, 'Test product 2', 'test 2');