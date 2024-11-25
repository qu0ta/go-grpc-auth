INSERT INTO apps (id, name, secret)
VALUES (1, 'app1', 'secret1')
ON CONFLICT DO NOTHING;