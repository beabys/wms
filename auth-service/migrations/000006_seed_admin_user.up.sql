-- Seed default admin user for development
-- Email: admin@wms.local
-- Password: admin123
INSERT INTO users (id, email, password_hash, name, role, active, created_at, updated_at)
VALUES (
  '00000000-0000-0000-0000-000000000001',
  'admin@wms.local',
  '$2a$10$cNASgkJG0.8Ps4W/b1ESZePG6O.zikA9CgKVSbC80kWZ7Tyf4H8Dq',
  'System Admin',
  'admin',
  true,
  NOW(),
  NOW()
)
ON CONFLICT (email) DO NOTHING;
