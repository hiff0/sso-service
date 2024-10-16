INSERT INTO apps (id, name, secret)
VALUES (1, 'test-app', 'very-secret')
ON CONFLICT DO NOTHING;